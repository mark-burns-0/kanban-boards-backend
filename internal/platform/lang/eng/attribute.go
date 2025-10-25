package eng

type Package struct{}

var attribute = map[string]string{
	"user_id":               "User",
	"category_id":           "Category",
	"platform_id":           "Platform",
	"password":              "Password",
	"password_confirmation": "Password confirmation",
	"email":                 "E-mail",
	"name":                  "Name",
	"firstname":             "Firstname",
	"lastname":              "Lastname",
	"patronymic":            "Patronymic",
	"text":                  "Text",
	"is_private":            "Private",
	"creator_id":            "Creator",
	"color":                 "Color",
	"tag":                   "Tag",
	"board_id":              "Board",
	"card_id":               "Card",
	"description":           "Description",
	"column_id":             "Column",
	"position":              "Position",
}

func (p *Package) GetAttribute(field string) string {
	return attribute[field]
}
