package storage

// Storage is an interface for interacting with the File System
// or a cloud storage service like GCP
type Storage interface {
	Read(filename string) (string, error)
}
