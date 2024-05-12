package highlight

import (
	"encoding/json"
	"log"
	"mime/multipart"

	t "github.com/sikozonpc/notebase/types"
)

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
