package ginhelp

import (
	"github.com/gin-gonic/gin"
)

type respones struct {
	c         *gin.Context
	httpCode int
	code     int
	message  string
	result   interface{}
}

func New(c *gin.Context) *respones {
	return &respones{c: c}
}

func (r *respones) HttpCode(code int) *respones {
	r.httpCode = code
	return r
}

func (r *respones) Code(code int) *respones {
	r.code = code
	return r
}

func (r *respones) Msg(msg string) *respones {
	r.message = msg
	return r
}

func (r *respones) Result(data interface{}) *respones {
	r.result = data
	return r
}

// {"code":0, "result": "xxx"}
func (r *respones) Send() {
	if r.httpCode == 0 {
		r.httpCode = 200
	}

	content := gin.H{"code": 0, "result": r.result}
	r.c.JSON(r.httpCode, content)
}

// {"code": -32600, "message": "Invalid Request"}
func (r *respones) Abort() {
	if r.httpCode == 0 {
		r.httpCode = 500
	}

	if r.code == 0 {
		panic("logical error")
	}

	if len(r.message) == 0 {
		r.message = findMessage(r.code)
	}

	content := gin.H{"code": r.code, "message": r.message}
	r.c.AbortWithStatusJSON(r.httpCode, content)
}

var findMessage func(code int) string
func SetMessageCb(cb func(int) string) {
	findMessage = cb
}

func Abort(c *gin.Context, code int) {
	New(c).Code(code).Abort()
}

func Result(c *gin.Context, result interface{}) {
	New(c).Result(result).Send()
}

