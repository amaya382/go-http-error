# go-http-error
go-http-error provides pretty formatted GCP-like http error responses.


## Usage
```go
errRes := httperror.NewHTTPErr(httperror.InvalidArgument,
	httperror.NewInnerErr("fooService", "invalidArgument", "id", "requestBody", "Passed id is invalid"),
	httperror.NewInnerErr("fooService", "invalidArgument", "name", "requestBody", "Passed name is invalid"))
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

