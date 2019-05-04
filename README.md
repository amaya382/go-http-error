# xerrorz
xerrorz provides GCP-like pretty formatted http error responses and also nested error logs powered by xerrors.



## Supported F/W
xerrorz has helper functions specialized for several web frameworks
* `net/http`
* [gin-gonic](https://github.com/gin-gonic/gin)


## Error Contents
You can pick up from preset errors for top-level http error contents (`HTTPErrDoc`). See [ErrType](https://github.com/amaya382/xerrorz/blob/feat/xerrorz/xerrorz.go#L130)

```go
type HTTPErr struct {
	ErrDoc HTTPErrDoc `json:"error"`
}

type HTTPErrDoc struct {
	Errors  []*InnerErr `json:"errors"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

type InnerErr struct {
	Domain       string `json:"domain"`       // global, {yourServiceName}, usage,...
	Reason       string `json:"reason"`       // invalidParameter, required,...
	Location     string `json:"location"`     // Authorization, {paramName},...
	LocationType string `json:"locationType"` // header, parameter, requestBody,...
	Message      string `json:"message"`      // {description}
	Cause        error  `json:"-"`            // For internal-use error reporting, NOT included in error jsons
}
```

```json
{
  "error": {
    "errors": [
      {
        "domain": "fooService",
        "reason": "invalidArgument",
        "location": "id",
        "locationType": "requestBody",
        "message": "Passed id is invalid"
      }
    ],
    "code": 400,
    "message": "Invalid argument"
  }
}
```


## Usage
Example scenario: `io.ErrClosedPipe` causes `invalidArgument.id` and `io.ErrNoProgress` causes `invalidArgument.name`

```go
import (
	"github.com/amaya382/xerrorz"
	"encoding/json"
	"fmt"
	"io" // For illustration
	"golang.org/x/xerrors"
)

// Example nested errors causing `invalidArgument.id`
e1 := xerrors.Errorf("e1: %w", io.ErrClosedPipe) // Original error
e2 := xerrors.Errorf("e2: %w", e1)               // Propagated once
e3 := xerrors.Errorf("e3: %w", e2)               // Propagated twice

errRes := xerrorz.NewHTTPErr(xerrorz.InvalidArgument,
	xerrorz.NewInnerErr("fooService", "invalidArgument", "id", "requestBody", "Passed id is invalid", e3),
	xerrorz.NewInnerErr("fooService", "invalidArgument", "name", "requestBody", "Passed name is invalid",
		io.ErrNoProgress))
bJSON, _ := json.Marshal(errRes)
fmt.Println(string(bJSON))  // (1) Print json
fmt.Printf("%+v\n", errRes) // (2) Print error for logging
```

### (1) Generated Response JSON
```json
{
  "error": {
    "errors": [
      {
        "domain": "fooService",
        "reason": "invalidArgument",
        "location": "id",
        "locationType": "requestBody",
        "message": "Passed id is invalid"
      },
      {
        "domain": "fooService",
        "reason": "invalidArgument",
        "location": "name",
        "locationType": "requestBody",
        "message": "Passed name is invalid"
      }
    ],
    "code": 400,
    "message": "Invalid argument"
  }
}
```

### (2) Nested and Detailed Error Log powered by xerrors
```
[HTTP Status 400] Invalid argument:
    github.com/amaya382/xerrorz.TestJSONEquality0
        /home/amaya/work/xerrorz/xerrorz_test.go:41
  - Invalid argument:
    github.com/amaya382/xerrorz.NewHTTPErr
        /home/amaya/work/xerrorz/xerrorz.go:103
  - Passed id is invalid:
    github.com/amaya382/xerrorz.TestJSONEquality0
        /home/amaya/work/xerrorz/xerrorz_test.go:42
  - e3:
    github.com/amaya382/xerrorz.TestJSONEquality0
        /home/amaya/work/xerrorz/xerrorz_test.go:40
  - e2:
    github.com/amaya382/xerrorz.TestJSONEquality0
        /home/amaya/work/xerrorz/xerrorz_test.go:39
  - e1:
    github.com/amaya382/xerrorz.TestJSONEquality0
        /home/amaya/work/xerrorz/xerrorz_test.go:38
  - io: read/write on closed pipe:
  - Passed name is invalid:
    github.com/amaya382/xerrorz.TestJSONEquality0
        /home/amaya/work/xerrorz/xerrorz_test.go:43
  - multiple Read calls return no data or error
```


## Usage for gin
Helper functions set a status code, a content-type header, and a body.

```go
import (
	"github.com/amaya382/xerrorz"
	"github.com/amaya382/xerrorz/helper/xgin"
	"github.com/gin-gonic/gin"
	"io" // For illustration
)

	var c *gin.Context // Use a context from handlers

  // Status Code: 400
  // Header: Content-Type:application/json
  // Body: {error json}
	xgin.SetHTTPErrJSON(c, xerrorz.InvalidArgument,
		xerrorz.NewInnerErr("fooService", "invalidArgument", "id", "requestBody", "Passed id is invalid", nil),
		xerrorz.NewInnerErr("fooService", "invalidArgument", "name", "requestBody", "Passed name is invalid",
			io.ErrNoProgress))
```
