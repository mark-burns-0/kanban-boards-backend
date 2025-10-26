package utils

type autoParseKey struct {
	key   string
	value string
}

var (
	user = autoParseKey{
		key:   "UserID",
		value: "userID",
	}
	card = autoParseKey{
		key:   "CardID",
		value: "card_id",
	}
	board = autoParseKey{
		key:   "ID",
		value: "id",
	}
	column = autoParseKey{
		key:   "ColumnID",
		value: "column_id",
	}
	comment = autoParseKey{
		key:   "ID",
		value: "comment_id",
	}
)

func getAutoParseKeys() []autoParseKey {
	return []autoParseKey{
		user,
		card,
		board,
		column,
		comment,
	}
}
