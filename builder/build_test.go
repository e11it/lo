package builder

import (
	"bytes"
	"github.com/go-martini/martini"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestBuildForm(t *testing.T) {
	httpRecorder := httptest.NewRecorder()
	m := martini.Classic()
	m.Post("/", func(request *http.Request) {
		fd := &testMapForm{}
		fd_result := &testMapForm{
			UserName:     "SomeName",
			UserPassword: "SomePassword",
			Resident:     true,
			NoResident:   false,
			Token:        "1234567890",
			Gender:       "Женский",
		}
		if err := FormRead(fd, request); err != nil {
			t.Error("Error FormRead:", err.Error())
			return
		}

		if eq := reflect.DeepEqual(fd, fd_result); !eq {
			t.Error("Form not equal.\nActual:", fd, "\nExpected:", fd_result)
		} else {
			t.Log("Form equal.\nActual:", fd, "\nExpected:", fd_result)
		}
	})

	mpPayload, mpWriter := genMultipartForm()

	if req, err := http.NewRequest("POST", "/", mpPayload); err != nil {
		panic(err)
	} else {
		req.Header.Set("Content-Type", mpWriter.FormDataContentType())

		if err := mpWriter.Close(); err != nil {
			panic(err)
		}

		m.ServeHTTP(httpRecorder, req)
		switch httpRecorder.Code {
		case http.StatusNotFound:
			panic("Routing is messed up in test fixture (got 404): check methods and paths")
		case http.StatusInternalServerError:
			panic("Something bad happened")
		}
	}
}

func genMultipartForm() (*bytes.Buffer, *multipart.Writer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("name", "SomeName")
	writer.WriteField("password", "SomePassword")
	writer.WriteField("gender", "2")
	writer.WriteField("resident", "1")
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
	if err := mapForm(fd, form, nil); err != nil {
		t.Error("Error map form:", err.Error())
	}
	if eq := reflect.DeepEqual(fd, fd_result); !eq {
		t.Error("Form not equal.\nActual:", fd, "\nExpected:", fd_result)
	} else {
		t.Log("Form equal.\nActual:", fd, "\nExpected:", fd_result)
	}
}

func TestRequiredField(t *testing.T) {
	form := make(map[string][]string)
	form["name"] = append(form["name"], "SomeName")
	form["token"] = append(form["token"], "1234567890")
	form["resident"] = append(form["resident"], "1")
	form["gender"] = append(form["gender"], "2")

	fd := &testMapForm{}

	if err := mapForm(fd, form, nil); err == nil {
		t.Error("Should return error: \"No value for required field : UserPassword\"")
	}
}
