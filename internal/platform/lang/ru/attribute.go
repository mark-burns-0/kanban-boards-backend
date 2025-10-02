package ru

type Package struct{}

var attribute = map[string]string{
	"user_id":     "Пользователь",
	"category_id": "Категория",
	"platform_id": "Платформа",
	"passowrd":    "Пароль",
	"mail":        "Почта",
	"name":        "Название",
	"firstname":   "Имя",
	"lastname":    "Фамилия",
	"patronymic":  "Отчество",
	"text":        "Текст",
}

func (p *Package) GetAttribute(field string) string {
	return attribute[field]
}
