package data

import (
	"time"

	t "github.com/sikozonpc/notebase/types"
)


func NewHighlight(text, location, note string, userId, bookId int) *t.Highlight {
	return &t.Highlight{
		Text:     text,
		Location: location,
		Note:     note,
		UserId:   userId,
		BookId:   bookId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}