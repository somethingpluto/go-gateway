package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"go_gateway/public"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"reflect"
	"regexp"
	"strings"
)

var val *validator.Validate

//设置Translation
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//参照：https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
		//设置支持语言
		en := en.New()
		zh := zh.New()

		//设置国际化翻译器
		uni := ut.New(zh, zh, en)
		val = validator.New()

		//根据参数取翻译器实例
		locale := c.DefaultQuery("locale", "zh")
		trans, _ := uni.GetTranslator(locale)

		//翻译器注册到validator
		switch locale {
		case "en":
			_ = en_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("en_comment")
			})
			break
		default:
			_ = zh_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})

			//自定义验证方法
			//https://github.com/go-playground/validator/blob/v9/_examples/custom-validation/main.go
			_ = val.RegisterValidation("valid_username", validUsername)

			// 验证服务名称
			_ = val.RegisterValidation("valid_service_name", validServiceName)

			// 规则 非空
			_ = val.RegisterValidation("valid_rule", validRule)

			_ = val.RegisterValidation("valid_url_rewrite", validUrlRewrite)

			//TODO:header_transfor验证规则
			_ = val.RegisterValidation("valid_header_transfor", func(fl validator.FieldLevel) bool {
				//for _, ms := range strings.Split(fl.Field().String(), "\n") {
				//	if len(strings.Split(ms, " ")) != 3 {
				//		return false
				//	}
				//}
				return true
			})

			//TODO:IPList验证规则
			_ = val.RegisterValidation("valid_iplist", func(fl validator.FieldLevel) bool {
				return true
			})

			_ = val.RegisterValidation("valid_weightlist", func(fl validator.FieldLevel) bool {
				for _, ms := range strings.Split(fl.Field().String(), "\n") {
					match, _ := regexp.Match(`^\d+$`, []byte(ms))
					if !match {
						return false
					}
				}
				return true
			})

			//自定义验证器
			//https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
			_ = val.RegisterTranslation("valid_username", trans, func(ut ut.Translator) error {
				return ut.Add("valid_username", "用户名不能为空", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_username", fe.Field())
				return t
			})

			_ = val.RegisterTranslation("valid_service_name", trans, func(ut ut.Translator) error {
				return ut.Add("valid_service_name", "服务名不能为空 ", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_service_name", fe.Field())
				return t
			})

			// 规则粗无提示
			_ = val.RegisterTranslation("valid_rule", trans, func(ut ut.Translator) error {
				return ut.Add("valid_rule", "{0} 规则不能为空 ", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_rule", fe.Field())
				return t
			})

			// 规则粗无提示
			_ = val.RegisterTranslation("valid_url_rewrite", trans, func(ut ut.Translator) error {
				return ut.Add("valid_url_rewrite", "不符合输入格式 ", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_url_rewrite", fe.Field())
				return t
			})

			_ = val.RegisterTranslation("valid_header_transfor", trans, func(ut ut.Translator) error {
				return ut.Add("valid_header_transfor", "header转换 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_header_transfor", fe.Field())
				return t
			})

			_ = val.RegisterTranslation("valid_iplist", trans, func(ut ut.Translator) error {
				return ut.Add("valid_iplist", "ip列表 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_iplist", fe.Field())
				return t
			})

			_ = val.RegisterTranslation("valid_weightlist", trans, func(ut ut.Translator) error {
				return ut.Add("valid_weightlist", "权重列表 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_weightlist", fe.Field())
				return t
			})

			break
		}
		c.Set(public.TranslatorKey, trans)
		c.Set(public.ValidatorKey, val)
		c.Next()
	}
}

// 用户名规则验证
func validUsername(fl validator.FieldLevel) bool {
	// 用户名不能为空
	match, _ := regexp.Match(`^[\s\S]*.*[^\s][\s\S]*$`, []byte(fl.Field().String()))
	return match
}

// 服务名称验证
func validServiceName(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return false
	}
	return true
}

// 验证输入规则不能为空
func validRule(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return false
	}
	return true
}

// 验证路由重写
func validUrlRewrite(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}

	return true
}
