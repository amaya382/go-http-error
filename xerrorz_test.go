package xerrorz

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"

	"golang.org/x/xerrors"
)

const sampleJSON = `
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
}`

func TestJSONEquality0(t *testing.T) {
	e1 := xerrors.Errorf("e1: %w", io.ErrClosedPipe)
	e2 := xerrors.Errorf("e2: %w", e1)
	e3 := xerrors.Errorf("e3: %w", e2)
	errRes := NewHTTPErr(InvalidArgument,
		NewInnerErr("fooService", "invalidArgument", "id", "requestBody", "Passed id is invalid", e3),
		NewInnerErr("fooService", "invalidArgument", "name", "requestBody", "Passed name is invalid", io.ErrNoProgress))
	bJSON, err := json.Marshal(errRes)
	if err != nil {
		t.Fatal("Failed to marshal an err object")
	}

	fmt.Printf("%+v\n", errRes)

	var o1, o2 interface{}
	err = json.Unmarshal(bJSON, &o1)
	if err != nil {
		t.Fatal("Failed to unmarshal an err json")
	}
	err = json.Unmarshal([]byte(sampleJSON), &o2)
	if err != nil {
		t.Fatal("Failed to unmarshal a sample json")
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Fatal("Inconsistent json was generated")
	}
}

func TestJSONEquality1(t *testing.T) {
	e1 := xerrors.Errorf("e1: %w", io.ErrClosedPipe)
	e2 := xerrors.Errorf("e2: %w", e1)
	e3 := xerrors.Errorf("e3: %w", e2)
	errRes := NewHTTPErr(InvalidArgument,
		NewInnerErr("fooService", "invalidArgument", "INVALID", "requestBody", "Passed id is invalid", e3),
		NewInnerErr("fooService", "invalidArgument", "name", "requestBody", "Passed name is invalid", io.ErrNoProgress))
	bJSON, err := json.Marshal(errRes)
	if err != nil {
		t.Fatal("Failed to marshal an err object")
	}

	var o1, o2 interface{}
	err = json.Unmarshal(bJSON, &o1)
	if err != nil {
		t.Fatal("Failed to unmarshal an err json")
	}
	err = json.Unmarshal([]byte(sampleJSON), &o2)
	if err != nil {
		t.Fatal("Failed to unmarshal a sample json")
	}

	if reflect.DeepEqual(o1, o2) {
		t.Fatal("Inconsistent json was generated")
	}
}

func TestContents0(t *testing.T) {
	e1 := xerrors.Errorf("e1: %w", io.ErrClosedPipe)
	e2 := xerrors.Errorf("e2: %w", e1)
	e3 := xerrors.Errorf("e3: %w", e2)
	errRes := NewHTTPErr(InvalidArgument,
		NewInnerErr("fooService", "invalidArgument", "id", "requestBody", "Passed id is invalid", e3),
		NewInnerErr("fooService", "invalidArgument", "name", "requestBody", "Passed name is invalid", io.ErrNoProgress))

	if errRes.ErrDoc.Code != 400 {
		t.Fatal("Invalid status code")
	}
}
