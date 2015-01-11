package builder

import (
	"bytes"
	"github.com/go-martini/martini"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestBuildForm(t *testing.T) {
	httpRecorder := httptest.NewRecorder()
	m := martini.Classic()
	m.Post("/", func(request *http.Request) {
		request.ParseForm()
		t.Log(">>>", request.PostForm)
	})
	req, err := http.NewRequest("POST", "/", strings.NewReader(`name=TestName&password=SomePassword&token=1234567890`))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "multipart/form-data")

	m.ServeHTTP(httpRecorder, req)
}

func genMultipartForm() (*bytes.Buffer, *multipart.Writer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("name", "SomeName")
	writer.WriteField("password", "SomePassword")
	writer.WriteField("token", "1234567890")

	return body, writer
}

type testMapForm struct {
	UserName     string `required:"true" field:"name" name:"Имя пользователя" type:"text"`
	UserPassword string `required:"true" field:"password" name:"Пароль пользователя" type:"password"`
	Resident     bool   `field:"resident" type:"radio" radio:"1;checked" name:"Резидент РФ"`
	NoResident   bool   `field:"resident" type:"radio" radio:"2" name:"Не резидент РФ"`
	Token        string `field:"token" type:"hidden" default:"true"`
	Gender       string `field:"gender" name:"Пол" type:"select" select:"Не известный=3;selected,Мужской=1,Женский=2"`
}

func TestMap(t *testing.T) {

	form := make(map[string][]string)
	form["name"] = append(form["name"], "SomeName")
	form["password"] = append(form["password"], "SomePassword")
	form["token"] = append(form["token"], "1234567890")
	form["resident"] = append(form["resident"], "1")
	form["gender"] = append(form["gender"], "2")

	fd := &testMapForm{}
	fd_result := &testMapForm{
		UserName:     "SomeName",
		UserPassword: "SomePassword",
		Resident:     true,
		NoResident:   false,
		Token:        "1234567890",
		Gender:       "Женский",
	}
	mapForm(fd, form, nil)
	if eq := reflect.DeepEqual(fd, fd_result); !eq {
		t.Error("Form not equal.\nActual:", fd, "\nExpected:", fd_result)
	} else {
		t.Log("Form equal.\nActual:", fd, "\nExpected:", fd_result)
	}
}
