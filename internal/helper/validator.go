package helper

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/id"
	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ut "github.com/go-playground/universal-translator"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

var Validate = validator.New()

func Validation(req interface{}) error {
	var uni *ut.UniversalTranslator
	var trans ut.Translator
	var validate = validator.New()

	id := id.New()
	uni = ut.New(id, id)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ = uni.GetTranslator("id")

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})

	id_translations.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(req)
	if err != nil {
		errSlice := err.(validator.ValidationErrors)

		var se []string
		for _, e := range errSlice {
			message := strings.Replace(e.Translate(trans), "_", " ", 1)
			f := fmt.Sprintf(`{"field":"%s", "message":"%s"}`, strings.ToLower(e.Field()), message)
			se = append(se, f)
		}

		st := status.New(codes.Unknown, "validation")
		v := &errdetails.DebugInfo{
			StackEntries: se,
			Detail:       "Please check your payload request",
		}
		ds, errD := st.WithDetails(v)

		if errD != nil {
			return st.Err()
		}
		return ds.Err()
	}

	return nil
}
