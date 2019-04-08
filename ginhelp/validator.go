package ginhelp

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"reflect"
)

var (
	defaultFuncs = make(map[string]validator.Func)
)

func init() {
	defaultFuncs["username"] = IsUsername
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
}


func IsUsername(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	if s, ok := field.Interface().(string); ok {
		return len(s) > 4
	}
	return false
}
