package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 自定义绑定密码验证规则，使用正则表达

var PasswdValidator validator.Func = func(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if ok {
		return regexp.MustCompile(`^([0-9].*[a-zA-Z]|[a-zA-Z].*[0-9]).*$`).MatchString(data)
	}
	return ok
	
}
