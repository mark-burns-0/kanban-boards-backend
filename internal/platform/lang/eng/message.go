package eng

var responseMessages map[string]string = map[string]string{
	"created":  "Created successfully",
	"updated":  "Updated successfully",
	"deleted":  "Deleted successfully",
	"moved":    "Moved successfully",
	"archived": "Archived successfully",
}

func (p *Package) GetResponseMessage(key string) string {
	return responseMessages[key]
}
