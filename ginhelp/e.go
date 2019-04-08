package ginhelp

import "net/http"

type ErrCode int
const (
	_ ErrCode = iota
)

var msgMap = map[ErrCode]string{}

func (code ErrCode) String() string {
	if msg, ok := msgMap[code]; ok {
		return msg
	}

	return http.StatusText(int(code))
}

func Register(code ErrCode, msg string) {
	msgMap[code] = msg
}

func init() {

}