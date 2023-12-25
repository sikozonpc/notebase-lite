package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/storage"
	"github.com/sikozonpc/notebase/config"
)

type GCPStorage struct {
	client *storage.Client

	w   io.Writer
	ctx context.Context
	// failed indicates that one or more of the steps failed.
	failed bool
}

func NewGCPStorage(ctx context.Context) (*GCPStorage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Println(ctx, "failed to create client: %v", err)
		return nil, err
	}
	defer client.Close()

	buf := &bytes.Buffer{}
	config := &GCPStorage{
		w:      buf,
		ctx:    ctx,
		client: client,
	}

	return config, nil
}

func (d *GCPStorage) errorf(format string, args ...interface{}) {
	d.failed = true
	fmt.Fprintln(d.w, fmt.Sprintf(format, args...))
	// cloud logging: log.Println(d.ctx, format, args...)
}

func (s *GCPStorage) Read(filename string) (string, error) {
	bucketName := config.Envs.GCPBooksBucketName
	bucket := s.client.Bucket(bucketName)

	rc, err := bucket.Object(filename).NewReader(s.ctx)
	if err != nil {
		s.errorf("unable to open file from bucket %q, file %q: %v", bucketName, filename, err)
		return "", err
	}
	defer rc.Close()

	slurp, err := io.ReadAll(rc)
	if err != nil {
		s.errorf("unable to read data from bucket %q, file %q: %v", bucketName, filename, err)
		return "", err
	}

	fmt.Fprintf(s.w, "%s\n", bytes.SplitN(slurp, []byte("\n"), 2)[0])

	if len(slurp) > 1024 {
		fmt.Fprintf(s.w, "...%s\n", slurp[len(slurp)-1024:])
	} else {
		fmt.Fprintf(s.w, "%s\n", slurp)
	}

	return string(slurp), nil
}
