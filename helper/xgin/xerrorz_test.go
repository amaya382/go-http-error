package xgin

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/amaya382/xerrorz"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

func TestSetHTTPErrJSON0(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SetHTTPErrJSON(c, xerrorz.InvalidArgument)

	// fmt.Println(res.Code)
	// fmt.Println(res.HeaderMap)
	// fmt.Println(res.Body)
}

func TestSetHTTPErrJSON1(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e1 := xerrors.Errorf("e1: %w", io.ErrClosedPipe)
	e2 := xerrors.Errorf("e2: %w", e1)
	e3 := xerrors.Errorf("e3: %w", e2)

	SetHTTPErrJSON(c, xerrorz.InvalidArgument,
		xerrorz.NewInnerErr(
			"fooService", "invalidArgument", "id", "requestBody", "Passed id is invalid", e3))

	// fmt.Println(res.Code)
	// fmt.Println(res.HeaderMap)
	// fmt.Println(res.Body)
}

func TestSetHTTPErrJSON2(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	e1 := xerrors.Errorf("e1: %w", io.ErrClosedPipe)
	e2 := xerrors.Errorf("e2: %w", e1)
	e3 := xerrors.Errorf("e3: %w", e2)

	f1 := xerrors.Errorf("f1: %w", io.ErrNoProgress)
	f2 := xerrors.Errorf("f2: %w", f1)

	SetHTTPErrJSON(c, xerrorz.InvalidArgument,
		xerrorz.NewInnerErr(
			"fooService", "invalidArgument", "id", "requestBody", "Passed id is invalid", e3),
		xerrorz.NewInnerErr(
			"fooService", "invalidArgument", "name", "requestBody", "Passed name is invalid", f2))

	// fmt.Println(res.Code)
	// fmt.Println(res.HeaderMap)
	// fmt.Println(res.Body)
}
