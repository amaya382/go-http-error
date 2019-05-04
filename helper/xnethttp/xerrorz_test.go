package xnethttp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/amaya382/xerrorz"
	"golang.org/x/xerrors"
)

const sampleJSON0 = `
{
    "error": {
        "errors": [],
        "code": 400,
        "message": "Invalid argument"
    }
}
`

const sampleJSON1 = `
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
`

const sampleJSON2 = `
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
`

func TestSetHTTPErrJSON0(t *testing.T) {
	res := httptest.NewRecorder()
	var w http.ResponseWriter
	w = res

	SetHTTPErrJSON(w, xerrorz.InvalidArgument)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("Invalid status code: %d\n", res.Code)
	}

	if res.HeaderMap.Get("Content-Type") != "application/json" {
		t.Fatalf("Invalid header: Content-Type:%s\n", res.HeaderMap.Get("Content-Type"))
	}

	var o1, o2 interface{}
	err := json.Unmarshal(res.Body.Bytes(), &o1)
	if err != nil {
		t.Fatalf("Failed to unmarshal an err json: %+v\n", err)
	}
	err = json.Unmarshal([]byte(sampleJSON0), &o2)
	if err != nil {
		t.Fatalf("Failed to unmarshal a sample json: %+v\n", err)
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Fatalf("Inconsistent json was generated: %+v\n", err)
	}
}

func TestSetHTTPErrJSON1(t *testing.T) {
	res := httptest.NewRecorder()
	var w http.ResponseWriter
	w = res

	e1 := xerrors.Errorf("e1: %w", io.ErrClosedPipe)
	e2 := xerrors.Errorf("e2: %w", e1)
	e3 := xerrors.Errorf("e3: %w", e2)

	SetHTTPErrJSON(w, xerrorz.InvalidArgument,
		xerrorz.NewInnerErr(
			"fooService", "invalidArgument", "id", "requestBody", "Passed id is invalid", e3))

	if res.Code != http.StatusBadRequest {
		t.Fatal("Invalid status code")
	}

	if res.HeaderMap.Get("Content-Type") != "application/json" {
		t.Fatal("Invalid header")
	}

	var o1, o2 interface{}
	err := json.Unmarshal(res.Body.Bytes(), &o1)
	if err != nil {
		t.Fatalf("Failed to unmarshal an err json: %+v\n", err)
	}
	err = json.Unmarshal([]byte(sampleJSON1), &o2)
	if err != nil {
		t.Fatalf("Failed to unmarshal a sample json: %+v\n", err)
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Fatalf("Inconsistent json was generated: %+v\n", err)
	}
}

func TestSetHTTPErrJSON2(t *testing.T) {
	res := httptest.NewRecorder()
	var w http.ResponseWriter
	w = res

	e1 := xerrors.Errorf("e1: %w", io.ErrClosedPipe)
	e2 := xerrors.Errorf("e2: %w", e1)
	e3 := xerrors.Errorf("e3: %w", e2)

	f1 := xerrors.Errorf("f1: %w", io.ErrNoProgress)
	f2 := xerrors.Errorf("f2: %w", f1)

	SetHTTPErrJSON(w, xerrorz.InvalidArgument,
		xerrorz.NewInnerErr(
			"fooService", "invalidArgument", "id", "requestBody", "Passed id is invalid", e3),
		xerrorz.NewInnerErr(
			"fooService", "invalidArgument", "name", "requestBody", "Passed name is invalid", f2))

	if res.Code != http.StatusBadRequest {
		t.Fatal("Invalid status code")
	}

	if res.HeaderMap.Get("Content-Type") != "application/json" {
		t.Fatal("Invalid header")
	}

	var o1, o2 interface{}
	err := json.Unmarshal(res.Body.Bytes(), &o1)
	if err != nil {
		t.Fatalf("Failed to unmarshal an err json: %+v\n", err)
	}
	err = json.Unmarshal([]byte(sampleJSON2), &o2)
	if err != nil {
		t.Fatalf("Failed to unmarshal a sample json: %+v\n", err)
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Fatalf("Inconsistent json was generated: %+v\n", err)
	}
}
