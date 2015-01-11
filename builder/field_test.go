package builder

import (
	"reflect"
	"testing"
	"time"
)

type (
	FormTests struct {
		description string
		payload     *formField
		expected    string
	}
)

var FieldT = []FormTests{
	{
		description: "Input type Text",
		payload:     &formField{Field: "name", Name: "Имя пользователя", Type: "text"},
		expected: `<label for="name">Имя пользователя</label>
<input type="text" name="name"><br>`,
	},
	{
		description: "Input type Password",
		payload:     &formField{Field: "password", Name: "Пароль пользователя", Type: "password"},
		expected: `<label for="password">Пароль пользователя</label>
<input type="password" name="password"><br>`,
	},
	{
		description: "Input type Button",
		payload:     &formField{Field: "button", Type: "button"},
		expected:    `<input type="button" name="button"><br>`,
	},
	{
		description: "Input type Hidden",
		payload:     &formField{Field: "token", Type: "hidden", Default: true, ReflectValue: reflect.ValueOf("345625145123451234123412342345")},
		expected:    `<input type="hidden" name="token" value="345625145123451234123412342345"><br>`,
	},
	{
		description: "Input type Radio",
		payload:     &formField{Field: "resident", Name: "Резидент РФ", Type: "radio", Default: true, Ext: "1;checked"},
		expected: `<label for="resident">Резидент РФ</label>
<input type="radio" name="resident" value="1" checked><br>`,
	},
	{
		description: "Input type Select",
		payload:     &formField{Field: "gender", Name: "Пол", Type: "select", Default: true, Ext: "Не известный=3;selected,Мужской=1,Женский=2"},
		expected: `<label for="gender">Пол</label>
<select name="gender">
	<option value="3" selected>Не известный</option>
	<option value="1">Мужской</option>
	<option value="2">Женский</option>
</select><br>`,
	},
}
var RValueT = []FormTests{
	{
		description: "Int",
		payload:     &formField{ReflectValue: reflect.ValueOf((int)(32))},
		expected:    "32",
	},
	{
		description: "Int64",
		payload:     &formField{ReflectValue: reflect.ValueOf((int64)(322342342341234))},
		expected:    "322342342341234",
	},
	{
		description: "Bool",
		payload:     &formField{ReflectValue: reflect.ValueOf((bool)(true))},
		expected:    "true",
	},
	{
		description: "String",
		payload:     &formField{ReflectValue: reflect.ValueOf((string)("absdjkfa фываорфдылвоар фшгутфыва 1453245"))},
		expected:    "absdjkfa фываорфдылвоар фшгутфыва 1453245",
	},
	{
		description: "Uint",
		payload:     &formField{ReflectValue: reflect.ValueOf((uint)(2346))},
		expected:    "2346",
	},
	{
		description: "Time.Duration",
		payload:     &formField{ReflectValue: reflect.ValueOf((time.Duration(10) * time.Second))},
		expected:    "10000000000",
	},
}

func TestValueToString(t *testing.T) {
	for _, testCase := range RValueT {
		actStr, isSet := testCase.payload.valueIsDefined()
		if (actStr != testCase.expected) || !isSet {
			t.Errorf("Description <<%s>>.\nExpected: %s\nActual: %s\nIsSet: %t", testCase.description, testCase.expected, actStr, isSet)
		}
	}
}

func TestFiled(t *testing.T) {
	for _, testCase := range FieldT {
		testCase.payload.ValidateTypeSetCustoms(testCase.payload.Type)
		actStr := testCase.payload.GetHTML()
		if actStr != testCase.expected {
			t.Errorf("Description <<%s>>.\nExpected: %s\nActual: %s", testCase.description, testCase.expected, actStr)
		}
	}
}
