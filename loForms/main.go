package loForms

type MyForm struct {
	UserName     string `required:"true" field:"name" name:"Имя пользователя" type:"text"`
	UserPassword string `required:"true" field:"password" name:"Пароль пользователя" type:"password"`
	Resident     bool   `field:"resident" type:"radio" radio:"1;checked" name:"Резидент РФ"`
	NoResident   bool   `field:"resident" type:"radio" radio:"2" name:"Не резидент РФ"`
	Gender       string `field:"gender" name:"Пол" type:"select" select:"Не известный=3;selected,Мужской=1,Женский=2"`
	Age          int64  `field:"age" name:"Возраст" type:"text" default:"true"`
	Token        string `field:"token" type:"hidden" default:"true"`
}
