package highlight

import (
	"encoding/json"
	"log"
	"mime/multipart"
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

func parseKindleExtractFromString(file string) (*t.RawExtractBook, error) {
	raw := new(t.RawExtractBook)
	f := []byte(file)

	err := json.Unmarshal(f, raw)
	if err != nil {
		log.Println("error unmarshalling file: ", err)
		return nil, err
	}

	return raw, nil
}

func parseKindleExtractFromFile(file multipart.File) (*t.RawExtractBook, error) {
	decoder := json.NewDecoder(file)

	raw := new(t.RawExtractBook)
	if err := decoder.Decode(raw); err != nil {
		log.Println("error decoding file: ", err)
		return nil, err
	}

	return raw, nil
}
