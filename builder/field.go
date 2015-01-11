package builder

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type formField struct {
	Required            bool
	Field               string
	ReflectField        reflect.StructField
	ReflectValue        reflect.Value
	Default             bool
	Name                string
	Type                string
	Ext                 string
	ElBuilder           TypeFElBuilder
	BodyBuilder         TypeFBodyBuilder
	FormValPreprocessor TypeFormValPreprocessor
}

// Custom builders
// Function to build element. <element>
type TypeFElBuilder func() []string

// Function to build body of element. <element>BODY</element>
type TypeFBodyBuilder func() string

// Preprocessor for assigned form value
type TypeFormValPreprocessor func(string) (string, error)

func GetField(field reflect.StructField, val reflect.Value) (*formField, error) {
	// Field with default values
	rField := &formField{
		Required:     false,
		Default:      false,
		Name:         field.Tag.Get("name"),
		Field:        field.Tag.Get("field"),
		ReflectField: field,
		ReflectValue: val,
	}

	if field.Tag.Get("required") == "true" {
		rField.Required = true
	}

	if field.Tag.Get("default") == "true" {
		rField.Default = true
	}

	//
	rField.ValidateTypeSetCustoms(field.Tag.Get("type"))
	// For select and radio get tag with type name
	if yes := rField.HasExt(); yes {
		if field.Tag.Get(rField.Type) == "" {
			return nil, errors.New("No tag for " + rField.Type)
		}
		rField.Ext = field.Tag.Get(rField.Type)
	}

	return rField, nil
}

// Проверяет type. Определяет обработчики
func (self *formField) ValidateTypeSetCustoms(val string) {
	switch val {
	case "text", "password", "hidden", "button", "textarea", "checkbox":
		self.Type = val
		self.ElBuilder = self.defaultBuilder
		self.BodyBuilder = nil
		self.FormValPreprocessor = nil
		return
	case "radio":
		self.Type = val
		self.ElBuilder = self.radioBuilder
		self.BodyBuilder = nil
		self.FormValPreprocessor = self.getValueRadioFrom
		return
	case "select":
		self.Type = val
		self.ElBuilder = nil
		self.BodyBuilder = self.selectBodyBuilder
		self.FormValPreprocessor = self.getValueSelectFrom
		return
	default:
		log.Println("Type not set. Used default type 'text'")
		self.ValidateTypeSetCustoms("text")
		return
	}
}

// Возвращает true, если предполагается наличие дополнительного tag в стракутре
func (self *formField) HasExt() bool {
	if self.Type == "radio" || self.Type == "select" {
		return true
	}
	return false
}

// Генерирует html код для поля
func (self *formField) GetHTML() (string, error) {
	var (
		html    string
		body    string
		options []string
	)
	label := self.getHtmlLabel()
	//
	openEl, closeEl := self.getOpenCloseElement()
	if openEl == "" {
		return "", errors.New("Can't get open html element for" + self.ReflectField.Name)
	}
	options = append(options, fmt.Sprintf("name=\"%s\"", self.Field))

	if self.ElBuilder != nil {
		options = append(options, self.ElBuilder()...)
	}
	//
	strOptions := strings.Join(options, " ")
	//
	if self.BodyBuilder != nil {
		body = self.BodyBuilder()
	}
	html = fmt.Sprintf("%s<%s %s>%s%s<br>", label, openEl, strOptions, body, closeEl)
	return html, nil
}

// Присваивает значение поля.
// @val string строка полученная из формы
func (self *formField) SetValueFromString(val string) error {
	if val == "" && self.Required {
		return errors.New("No value for required field : " + self.ReflectField.Name)
	}
	// Preprocess input value
	if self.FormValPreprocessor != nil {
		if valnew, err := self.FormValPreprocessor(val); err != nil {
			return err
		} else {
			val = valnew
		}
	}
	// Assign input value
	return setWithProperType(self.ReflectValue.Kind(), val, self.ReflectValue)
}

// Возвращает label для текущего элемента
func (self *formField) getHtmlLabel() string {
	if len(self.Name) > 0 {
		return fmt.Sprintf("<label for=\"%s\">%s</label>\n", self.Field, self.Name)
	}
	return ""
}

// Возвращает имя открывающего и, если необходимо, закрывающего элемента
func (self *formField) getOpenCloseElement() (string, string) {
	switch self.Type {
	case "text", "password", "hidden", "button", "checkbox", "radio":
		return fmt.Sprintf("input type=\"%s\"", self.Type), ""
	case "select", "textarea":
		return self.Type, fmt.Sprintf("</%s>", self.Type)
	}
	return "", ""
}

// Pretify Value based on type.
func (self *formField) valueIsDefined() (string, bool) {
	var value string
	switch self.ReflectValue.Kind() {
	case reflect.Invalid:
		// Value not defined
		return "", false
	case reflect.String:
		value = self.ReflectValue.String()
		break
	case reflect.Bool:
		value = strconv.FormatBool(self.ReflectValue.Bool())
		break
	case reflect.Int, reflect.Int64:
		if self.ReflectValue.Type() == reflect.ValueOf((time.Duration)(0)).Type() {
			// Not a good idea )
			// value = rval.Interface().(time.Duration).String()
			value = strconv.FormatInt(self.ReflectValue.Int(), 10)
		} else {
			value = strconv.FormatInt(self.ReflectValue.Int(), 10)
		}
		break
	case reflect.Uint, reflect.Uint64:
		value = strconv.FormatUint(self.ReflectValue.Uint(), 10)
		break
	case reflect.Float64:
		value = strconv.FormatFloat(self.ReflectValue.Float(), 'f', 2, 64)
		break
	}

	if len(value) > 0 {
		return value, true
	}
	return "", false
}

//
func (self *formField) defaultBuilder() []string {
	var (
		options []string
	)
	if value, isset := self.valueIsDefined(); isset && self.Default {
		options = append(options, fmt.Sprintf("value=\"%s\"", value))
	}
	return options
}

// RADIO
func (self *formField) radioBuilder() []string {
	var (
		options []string
	)
	/*
	 * type:”radio” radio:”VALUE[;checked]”
	 * Где:
	 * VALUE– Значение radiobutton
	 * ;checked – Флаг выбора по умолчанию
	 */
	if strings.HasSuffix(self.Ext, ";checked") {
		value := strings.TrimSuffix(self.Ext, ";checked")
		options = append(options, fmt.Sprintf("value=\"%s\"", value))
		options = append(options, fmt.Sprintf("checked"))
	} else {
		options = append(options, fmt.Sprintf("value=\"%s\"", self.Ext))
	}

	return options
}

func (self *formField) getValueRadioFrom(formVal string) (string, error) {
	value := strings.TrimSuffix(self.Ext, ";checked")
	if len(value) < 0 {
		return "", errors.New("No value for radio element")
	}
	if formVal == value {
		return "true", nil
	}
	return "false", nil
}

func (self *formField) getValueSelectFrom(formVal string) (string, error) {
	options := strings.Split(self.Ext, ",")
	for _, str := range options {
		var (
			name, endStr, value string
		)
		eqIndex := strings.Index(str, "=")
		if eqIndex < 0 {
			return "", errors.New(fmt.Sprintln("Incorrect tag value. Pattern(NAME=VALUE[;selected]) not matched. Str:", str, "Full str: ", self.Ext))
		}
		name = str[0:eqIndex]
		endStr = str[eqIndex+1 : len(str)]
		value = strings.TrimSuffix(endStr, ";selected")

		if value == formVal {
			return name, nil
		}
	}
	return "", errors.New("Form value: " + formVal + "Doesn't match any defined values: " + self.Ext)
}

// Return body of select element(<options> html)
func (self *formField) selectBodyBuilder() string {
	var html string = "\n"
	/*
	 * type:"select" select:"[NAME=VALUE[;selected]],"
	 */
	options := strings.Split(self.Ext, ",")
	for _, str := range options {
		var (
			name, endStr, value, selected string
		)
		eqIndex := strings.Index(str, "=")
		if eqIndex < 0 {
			// TODO: Error processor
			errors.New(fmt.Sprintln("Incorrect tag value. Pattern(NAME=VALUE[;selected]) not matched. Str:", str, "Full str: ", self.Ext))
		}
		name = str[0:eqIndex]
		endStr = str[eqIndex+1 : len(str)]
		if strings.HasSuffix(endStr, ";selected") {
			value = strings.TrimSuffix(endStr, ";selected")
			selected = " selected"
		} else {
			value = endStr
		}

		html = fmt.Sprintf("%s\t<option value=\"%s\"%s>%s</option>\n", html, value, selected, name)

	}
	return html
}

// This sets the value in a struct of an indeterminate type to the
// matching value from the request (via Form middleware) in the
// same type, so that not all deserialized values have to be strings.
// Supported types are string, int, float, and bool.
func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val == "" {
			val = "0"
		}
		// !Custom setup for TimeDuration.
		/* // What will be returned from client?
		 * if structField.Type() == reflect.ValueOf((time.Duration)(0)).Type() {
		 *
		 * }
		 */
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return errors.New("TypeError: Value could not be parsed as integer")
		} else {
			structField.SetInt(intVal)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val == "" {
			val = "0"
		}
		uintVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return errors.New("TypeError: Value could not be parsed as unsigned integer")
		} else {
			structField.SetUint(uintVal)
		}
	case reflect.Bool:
		if val == "" {
			val = "false"
		}
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return errors.New("TypeError: Value could not be parsed as boolean")
		} else {
			structField.SetBool(boolVal)
		}
	case reflect.Float32:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return errors.New("TypeError: Value could not be parsed as 32-bit float")
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.Float64:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return errors.New("TypeError: Value could not be parsed as 64-bit float")
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.String:
		structField.SetString(val)
	}

	return nil
}
