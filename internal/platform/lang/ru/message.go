package ru

var responseMessages map[string]string = map[string]string{
	"created":  "Успешно создано",
	"updated":  "Успешно обновлено",
	"deleted":  "Успешно удалено",
	"moved":    "Успешно перемещено",
	"archived": "Успешно архивировано",
}

func (p *Package) GetResponseMessage(key string) string {
	return responseMessages[key]
}
