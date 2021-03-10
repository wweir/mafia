package s3

import (
	"context"
	"io/fs"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/wweir/mafia/drivers"
	"github.com/wweir/mafia/pkg/fsmock"
	"github.com/wweir/mafia/pkg/fspath"
)

var _ drivers.FSAdaptor = new(common)

type common struct {
	*s3.Client
	timeout time.Duration

	drivers.MockFSFull
}

func newCommon(timout time.Duration) (*common, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timout)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("cn-northwest-1"))
	if err != nil {
		return nil, err
	}

	return &common{
		timeout: timout,
		Client:  s3.NewFromConfig(cfg),
	}, nil
}

func (c *common) DeleteDir(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	bucket, object := fspath.SplitPrefixDir(path)

	if object == "" {
		log.Info().Interface("bucket", bucket).Msg("delete")
		// _, err := c.Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		// 	Bucket: aws.String(bucket),
		// })
		return nil
	}

	p := s3.NewListObjectsV2Paginator(c.Client, &s3.ListObjectsV2Input{
		Bucket:     aws.String(bucket),
		Prefix:     aws.String(object + "/"),
		FetchOwner: true,
	})

	files := []types.ObjectIdentifier{{Key: aws.String(object + "/")}}
	for p.HasMorePages() {
		out, err := p.NextPage(ctx)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, content := range out.Contents {
			files = append(files, types.ObjectIdentifier{
				Key: content.Key,
			})
		}
	}

	log.Info().Interface("files", files).Msg("delete")
	// _, err := c.Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
	// 	Bucket: aws.String(bucket),
	// 	Delete: &types.Delete{
	// 		Objects: files,
	// 	},
	// })
	return nil
}
func (c *common) DeleteFile(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	bucket, object := fspath.SplitPrefixDir(path)

	_, err := c.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	})
	return err
}

func (c *common) Rename(fromPath string, toPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	c.Client.CopyObject(ctx, &s3.CopyObjectInput{})

	fromBucket, fromObject := fspath.SplitPrefixDir(fromPath)
	toBucket, toObject := fspath.SplitPrefixDir(toPath)

	switch "" {
	case fromBucket, toBucket:
		return errors.New("bucket name must be setted")
	case fromObject, toObject:
		return errors.New("object name must be setted")
	default:
		if fromBucket != toBucket {
			return errors.Errorf("not impliment")
		}
		return nil
	}
}
func (c *common) Mkdir(name string, perm os.FileMode) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	bucket, object := fspath.SplitPrefixDir(name)
	if object == "" {
		_, err := c.Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		return err
	}
	return nil

	// c.Client.PutObject(ctx,&s3.PutObjectInput{

	// })
}

func (c *common) Stat(name string) (os.FileInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	bucket, object := fspath.SplitPrefixDir(name)
	if object == "" {
		return fsmock.MockFileInfo(name, true, 0, nil), nil
	}

	out, err := c.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	})
	if err != nil {
		return fsmock.MockFileInfo(name, true, 0, nil), nil
	}
	return fsmock.MockFileInfo(name, false, out.ContentLength, out.LastModified), nil
}

func (c *common) ReadDir(name string) ([]fs.FileInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	bucket, object := fspath.SplitPrefixDir(name)
	if bucket == "" {
		out, err := c.Client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		fis := make([]fs.FileInfo, 0, len(out.Buckets))
		for b := range out.Buckets {
			fis = append(fis, fsmock.MockFileInfo(*out.Buckets[b].Name, true, 0, nil))
		}
		return fis, nil
	}

	if object != "" {
		object += "/"
	}
	p := s3.NewListObjectsV2Paginator(c.Client, &s3.ListObjectsV2Input{
		Bucket:     aws.String(bucket),
		Prefix:     aws.String(object),
		Delimiter:  aws.String("/"),
		FetchOwner: true,
	})

	var fis = []fs.FileInfo{}
	for p.HasMorePages() {
		out, err := p.NextPage(ctx)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for p := range out.CommonPrefixes {
			fis = append(fis, fsmock.MockFileInfo(
				*out.CommonPrefixes[p].Prefix, true, 0, nil))
		}
		for c := range out.Contents {
			content := out.Contents[c]
			if *content.Key != object {
				fis = append(fis, fsmock.MockFileInfo(
					*content.Key, false, content.Size, content.LastModified))
			}
		}
	}
	return fis, nil
}
