package bucket

import (
	"io"
	"os"
	"reflect"
)

const (
	AwsProvider BucketType = iota
)

func New(t BucketType, cfg any) (b *Bucket, err error) {
	rt := reflect.TypeOf(cfg)

	switch bt {
	case "AwsProvider":
		//TODO implementar aws provider
	default:
		return nil, fmt.ErrorF("type not implemented")
	}

	return
}

type BucketType int

type BucketInterface interface {
	Upload(io.Reader, string) error
	Download(string, string) (*os.File, error)
	Delete(string) error
}

type Bucket struct {
	p BucketInterface
}

func (b *Bucket) Upload(file io.Reader, key string) error {
	return b.p.Upload(file, key)
}

func (b *Bucket) Download(src, dst string) (file *os.File, err error) {
	return b.p.Download(src, dst)
}

func (b *Bucket) Delete(key string) error {
	return b.p.Delete(key)
}
