package storage

type MemoryStorage struct{}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (m *MemoryStorage) Read(filename string) (string, error) {
	return fileContent, nil
}

var fileContent = `
	{
  "asin": "SOMERANDOMASIN",
  "title": "Some random book on kindle",
  "authors": "Some random author",
  "highlights": [
    {
      "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,",
      "isNoteOnly": false,
      "location": {
        "url": "kindle://book?action=open&asin=SOMERANDOMASIN&location=307",
        "value": 307
      },
      "note": "This is a note"
    },

    {
      "text": "consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,",
      "isNoteOnly": false,
      "location": {
        "url": "kindle://book?action=open&asin=SOMERANDOMASIN&location=742",
        "value": 742
      },
      "note": null
    }
  ]
}
`
