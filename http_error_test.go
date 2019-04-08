package httperror

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	fmt.Println("hogehoge")
	errRes := NewHTTPErr(InvalidArgument,
		NewInnerErr("fooService", "invalidArgument", "id", "requestBody", "Passed id is invalid"),
		NewInnerErr("fooService", "invalidArgument", "name", "requestBody", "Passed name is invalid"))
	bJSON, _ := json.Marshal(errRes)
	fmt.Println(string(bJSON))
}
