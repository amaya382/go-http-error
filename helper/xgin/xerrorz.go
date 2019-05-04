package xgin

import (
	"github.com/amaya382/xerrorz"
	"github.com/gin-gonic/gin"
)

func SetHTTPErrJSON(c *gin.Context, errType xerrorz.ErrType, innerErrs ...*xerrorz.InnerErr) {
	err := xerrorz.NewHTTPErr(errType, innerErrs...)
	c.JSON(err.ErrDoc.Code, err)
}
