package builder

import (
	"testing"
)

type MyForm struct {
	UserName     string `required:"true" field:"name" name:"Имя пользователя" type:"text"`
	UserPassword string `required:"true" field:"password" name:"Пароль пользователя" type:"password"`
	Resident     bool   `field:"resident" type:"radio" radio:"1;checked" name:"Резидент РФ"`
	NoResident   bool   `field:"resident" type:"radio" radio:"2" name:"Не резидент РФ"`
	Gender       string `field:"gender" name:"Пол" type:"select" select:"Не известный=3;selected,Мужской=1,Женский=2"`
	Age          int64  `field:"age" name:"Возраст" type:"text" default:"true"`
	Token        string `field:"token" type:"hidden" default:"true"`
}

var MyFormExp = `<label for="name">Имя пользователя</label>
<input type="text" name="name"><br>
<label for="password">Пароль пользователя</label>
<input type="password" name="password"><br>
<label for="resident">Резидент РФ</label>
<input type="radio" name="resident" value="1" checked><br>
<label for="resident">Не резидент РФ</label>
<input type="radio" name="resident" value="2"><br>
<label for="gender">Пол</label>
<select name="gender">
	<option value="3" selected>Не известный</option>
	<option value="1">Мужской</option>
	<option value="2">Женский</option>
</select><br>
<label for="age">Возраст</label>
<input type="text" name="age" value="18"><br>
<input type="hidden" name="token" value="345625145123451234123412342345"><br>`

func TestFormToHTML(t *testing.T) {
	var fd *MyForm = &MyForm{
		Age:   18,
		Token: "345625145123451234123412342345",
	}
	form, _ := FormCreate(fd)
	/*if form != MyFormExp {
		t.Error("From html doesn't match expected. \n------ Actual:\n", form, "\n------ Expected:\n", MyFormExp)
	}*/
	if len(form) != len(MyFormExp) {
		t.Error("Html length mismatch. Actual:", len(form), "Expected:", len(MyFormExp))
	}
	loop := len(form)
	if loop > len(MyFormExp) {
		loop = len(MyFormExp)
	}
	for i := 0; i < loop; i++ {
		if form[i:i+1] != MyFormExp[i:i+1] {
			t.Error("Form not match since:", i, "Matched part is:", form[:i+1])
			return
		}
	}
}
