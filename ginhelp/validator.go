package ginhelp

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

var (
	defaultFuncs = make(map[string]validator.Func)
)

func init() {
	//defaultFuncs["chinaPhone"] = IsUsername
	defaultFuncs["chinaPhone"] = IsChinaPhone
}

func RegisterValidation(funcs map[string]validator.Func) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("11")
	}

	for k, f := range defaultFuncs {
		if _, ok := funcs[k]; !ok {
			funcs[k] = f
		}
	}

	for k, f := range funcs {
		if err := v.RegisterValidation(k, f); err != nil {
			panic(err)
		}
	}

	v.RegisterAlias("username", "email|chinaPhone")
}

func IsUsername(fl validator.FieldLevel) bool {
	if s, ok := fl.Field().Interface().(string); ok {
		return len(s) > 4
	}
	return false
}

var chinaPhoneReg = regexp.MustCompile(`^1([38][0-9]|14[57]|5[^4])\d{8}$`)
func IsChinaPhone(fl validator.FieldLevel) bool {
	if s, ok := fl.Field().Interface().(string); ok {
		return chinaPhoneReg.MatchString(s)
	}
	return false
}
