package s3

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/wweir/mafia/drivers"
	"github.com/wweir/mafia/pkg/fspath"
)

var _ drivers.FTPAdaptor = new(FTP)

type FTP struct {
	*common
}

func NewFTP(timeout time.Duration) (*FTP, error) {
	c, err := newCommon(timeout)
	if err != nil {
		return nil, err
	}

	return &FTP{
		common: c,
	}, nil
}

type fakeWriteAt struct {
	*io.PipeWriter
}

func (rw *fakeWriteAt) WriteAt(p []byte, off int64) (n int, err error) {
	return rw.PipeWriter.Write(p)
}

func (ftp *FTP) GetFile(path string, offset int64) (int64, io.ReadCloser, error) {
	bucket, object := fspath.SplitPrefixDir(path)
	pipeR, pipeW := io.Pipe()
	downloader := manager.NewDownloader(ftp.Client, func(d *manager.Downloader) {
		d.Concurrency = 1
	})

	go func() {
		_, err := downloader.Download(context.Background(),
			&fakeWriteAt{PipeWriter: pipeW},
			&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(object),
			})
		log.Err(err).Msg("download from S3")
		pipeW.Close()
	}()

	return 0, pipeR, nil
}
