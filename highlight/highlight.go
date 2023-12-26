package highlight

import (
	"encoding/json"
	"log"
	"time"

	t "github.com/sikozonpc/notebase/types"
)

func New(text, location, note, bookId string, userId int) *t.Highlight {
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
func parseKindleExtractFile(file string, userId int) (*t.RawExtractBook, error) {
	raw := new(t.RawExtractBook)
	f := []byte(file)

	err := json.Unmarshal(f, raw)
	if err != nil {
		log.Println("error unmarshalling file: ", err)
		return nil, err
	}

	return raw, nil
}
