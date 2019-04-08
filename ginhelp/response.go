package ginhelp

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ginResponse(c *gin.Context, httpCode int, errCode ErrCode, data interface{}, abort bool) {
	content := Response{
		Code:int(errCode),
		Msg:errCode.String(),
		Data:data,
	}

	if abort {
		c.AbortWithStatusJSON(httpCode, content)
		return
	}
	c.JSON(httpCode, content)
}

func GinResponse(c *gin.Context, httpCode int, errCode ErrCode, data interface{}) {
	ginResponse(c, httpCode, errCode, data, false)
}

func GinAbort(c *gin.Context, httpCode int, errCode ErrCode, data interface{}) {
	ginResponse(c, httpCode, errCode, data, true)
}

