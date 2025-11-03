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
	"color":       "Цвет",
	"tag":         "Тег",
	"board_id":    "Доска",
	"card_id":     "Карточка",
	"description": "Описание",
	"column_id":   "Столбец",
	"position":    "Позиция",
}

func (p *Package) GetAttribute(field string) string {
	return attribute[field]
}
