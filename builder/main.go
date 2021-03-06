package builder

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"
)

var (
	// Maximum amount of memory to use when parsing a multipart form.
	// Set this to whatever value you prefer; default is 10 MB.
	MaxMemory = int64(1024 * 1024 * 10)
)

func FormCreate(formStruct interface{}) (string, error) {
	var (
		fields []string
		err    error
	)

	if !isPointer(formStruct) {
		return "", errors.New("FormStruct should be a pointer to struct")
	}

	typ := reflect.TypeOf(formStruct)
	val := reflect.ValueOf(formStruct)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		rawField := typ.Field(i)

		// Skip ignored("-"" or empty) and unexported field in the struct
		if rawField.Tag.Get("field") == "-" || rawField.Tag.Get("field") == "" || !val.Field(i).CanInterface() {
			continue
		}

		field, err := GetField(rawField, val.Field(i))
		if err != nil {
			return "", err
		}
		fieldHtml, ferr := field.GetHTML()
		if ferr != nil {
			return "", err
		}
		fields = append(fields, fieldHtml)
	}
	html := strings.Join(fields, "\n")
	return html, err
}

func FormRead(formStruct interface{}, request *http.Request) error {

	if !isPointer(formStruct) {
		return errors.New("FormStruct should be a pointer to struct")
	}

	contentType := request.Header.Get("Content-Type")

	if request.Method == "POST" || contentType != "" {
		if strings.Contains(contentType, "form-urlencoded") {
			return errors.New("Unsupported Content-Type")
		} else if strings.Contains(contentType, "multipart/form-data") {
			return MultipartForm(formStruct, request)
		} else {
			if contentType == "" {
				return errors.New("Empty Content-Type")
			} else {
				return errors.New("Unsupported Content-Type")
			}
		}
	}
	return nil
}

func DumpForm(formStruct interface{}) {
	typ := reflect.TypeOf(formStruct)
	val := reflect.ValueOf(formStruct)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		rawField := typ.Field(i)
		valField := val.Field(i)
		log.Println(">", rawField.Name, "[", rawField.Tag.Get("field"), "]  :", valField.Interface())
	}
}
