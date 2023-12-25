package highlight

import (
	"log"
	"time"

	t "github.com/sikozonpc/notebase/types"
)

func New(text, location, note string, userId, bookId int) *t.Highlight {
	return &t.Highlight{
		Text:      text,
		Location:  location,
		Note:      note,
		UserID:    userId,
		BookID:    bookId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

// Fetches the file from the cloud storage in the KindleExtract format and reads it by creating a any new entities (books and highlights) if needed.
func parseKindleExtractFile(file string, userId int) ([]*t.Highlight, error) {
	log.Println("file: ", file)

	return nil, nil
}