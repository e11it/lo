package builder

import (
	"errors"
	//"log"
	"mime/multipart"
	"net/http"
	"reflect"
)

func MultipartForm(formStruct interface{}, request *http.Request) error {
	// This if check is necessary due to https://github.com/martini-contrib/csrf/issues/6
	if request.MultipartForm == nil {
		// Workaround for multipart forms returning nil instead of an error
		// when content is not multipart; see https://code.google.com/p/go/issues/detail?id=6334
		if multipartReader, err := request.MultipartReader(); err != nil {
			return errors.New("DeserializationError: " + err.Error())
		} else {
			form, parseErr := multipartReader.ReadForm(MaxMemory)
			if parseErr != nil {
				return errors.New("DeserializationError: " + parseErr.Error())
			}
			request.MultipartForm = form
		}
	}

	if err := mapForm(formStruct, request.MultipartForm.Value, request.MultipartForm.File); err != nil {
		return err
	}
	return nil
}

func mapForm(formStruct interface{}, form map[string][]string, formfile map[string][]*multipart.FileHeader) error {
	typ := reflect.TypeOf(formStruct)
	val := reflect.ValueOf(formStruct)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {

		rawField := typ.Field(i)
		valField := val.Field(i)

		fieldName := rawField.Tag.Get("field")
		// Skip ignored("-"" or empty) and unexported field in the struct
		if fieldName == "-" || fieldName == "" || !valField.CanSet() {
			continue
		}

		field, ferr := GetField(rawField, valField)
		if ferr != nil {
			return ferr
		}

		if inputValue, exists := form[fieldName]; exists {
			numElems := len(inputValue)
			if val.Kind() == reflect.Slice && numElems > 0 {
				return errors.New("Unsupported form type: slice")
			} else {
				// log.Println("Set " + rawField.Name + ", field:" + fieldName + " = " + inputValue[0])
				//
				if err := field.SetValueFromString(inputValue[0]); err != nil {
					return err
				}
			}
			continue
		} else if _, exists_file := formfile[fieldName]; exists_file {
			return errors.New("Unsupported form type: file")
		} else {
			if field.Required {
				return errors.New("No value for required field : " + rawField.Name)
			}
		}
	}

	return nil
}

func isPointer(obj interface{}) bool {
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		return true
	}
	return false
}
