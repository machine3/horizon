package s3

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"

	"g.hz.netease.com/horizon/pkg/util/errors"
)

type Interface interface {
	// PutObject 上传对象，path为对象路径，content为对象内容的reader，metadata为对象的元数据
	PutObject(ctx context.Context, path string, content io.ReadSeeker, metadata map[string]string) error
	// GetObject 获取对象内容
	GetObject(ctx context.Context, path string) ([]byte, error)
	// CopyObject 同一个桶内的对象拷贝，从srcPath拷贝到destPath
	CopyObject(ctx context.Context, srcPath, destPath string) error
	// ListObjects 获取以prefix开头的对象，最多返回maxKeys个;
	// NOTE: s3的ListObjects返回结果是以key的字母序正序排序的，而不是按照上传时间排序的
	// Ref: https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjects.html
	ListObjects(ctx context.Context, prefix string, maxKeys int64) ([]*awss3.Object, error)
	// DeleteObjects 删除所有以prefix开头的对象
	DeleteObjects(ctx context.Context, prefix string) error
	//GetSignedObjectURL  获取上传文件的url
	GetSignedObjectURL(path string, expire time.Duration) (string, error)
}

type Params struct {
	// s3存储的AccessKey
	AccessKey string
	// s3存储的SecretKey
	SecretKey string
	// s3存储的Region
	Region string
	// s3存储的Endpoint
	Endpoint string
	// s3存储的桶
	Bucket string
	// 是否要跳过https证书校验
	SkipVerify bool
	// 是否使用子目录形式
	S3ForcePathStyle bool
	// 对象的ContentType
	ContentType string
	// s3 sdk log level
	LogLevel *aws.LogLevelType
}

func (params *Params) check() error {
	const op = "s3 params check"
	if len(params.AccessKey) == 0 {
		return errors.E(op, "AccessKey must be specified")
	}
	if len(params.SecretKey) == 0 {
		return errors.E(op, "SecretKey must be specified")
	}
	if len(params.Region) == 0 {
		return errors.E(op, "Region must be specified")
	}
	if len(params.Bucket) == 0 {
		return errors.E(op, "Bucket must be specified")
	}
	return nil
}

type Driver struct {
	Params
	S3 *awss3.S3
}

func NewDriver(params Params) (Interface, error) {
	const op = "new s3 driver"
	if err := params.check(); err != nil {
		return nil, err
	}
	d := &Driver{Params: params}

	awsConfig := aws.NewConfig()
	cred := credentials.NewStaticCredentials(params.AccessKey, params.SecretKey, "")
	awsConfig.WithCredentials(cred)
	awsConfig.WithRegion(params.Region)
	awsConfig.WithS3ForcePathStyle(params.S3ForcePathStyle)
	if len(params.Endpoint) > 0 {
		awsConfig.WithEndpoint(params.Endpoint)
	}
	if params.SkipVerify {
		awsConfig.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		})
	}
	if params.LogLevel != nil {
		awsConfig.WithLogLevel(*params.LogLevel)
	}
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, errors.E(op, err)
	}

	d.S3 = awss3.New(sess)

	return d, nil
}

// PutObject 上传对象，path为对象路径，content为对象内容的reader，metadata为对象的元数据
func (d *Driver) PutObject(ctx context.Context, path string, content io.ReadSeeker, metadata map[string]string) error {
	_, err := d.S3.PutObjectWithContext(ctx, &awss3.PutObjectInput{
		Body:        content,
		Bucket:      aws.String(d.Bucket),
		ContentType: aws.String(d.ContentType),
		Key:         aws.String(path),
		Metadata: func() map[string]*string {
			if metadata == nil {
				return nil
			}
			ret := make(map[string]*string)
			for k, v := range metadata {
				ret[k] = aws.String(v)
			}
			return ret
		}(),
	})
	return err
}

func (d *Driver) GetObject(ctx context.Context, path string) ([]byte, error) {
	output, err := d.S3.GetObjectWithContext(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = output.Body.Close() }()

	content, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (d *Driver) GetSignedObjectURL(path string, expire time.Duration) (string, error) {
	req, _ := d.S3.GetObjectRequest(&awss3.GetObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(path),
	})
	urlStr, err := req.Presign(expire)

	if err != nil {
		return "", err
	}
	return urlStr, nil
}

// CopyObject 同一个桶内的对象拷贝，从srcPath拷贝到destPath
func (d *Driver) CopyObject(ctx context.Context, srcPath, destPath string) error {
	_, err := d.S3.CopyObjectWithContext(ctx, &awss3.CopyObjectInput{
		Bucket:     aws.String(d.Bucket),
		CopySource: aws.String(fmt.Sprintf("/%s/%s", d.Bucket, srcPath)),
		Key:        aws.String(destPath),
	})
	return err
}

// ListObjects 获取以prefix开头的对象，最多返回maxKeys个;
// NOTE: s3的ListObjects返回结果是以key的字母序正序排序的，而不是按照上传时间排序的
// Ref: https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjects.html
func (d *Driver) ListObjects(ctx context.Context, prefix string, maxKeys int64) ([]*awss3.Object, error) {
	output, err := d.S3.ListObjectsWithContext(ctx, &awss3.ListObjectsInput{
		Bucket:  aws.String(d.Bucket),
		MaxKeys: aws.Int64(maxKeys),
		Prefix:  aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}
	return output.Contents, nil
}

// DeleteObjects 删除所有以prefix开头的对象
func (d *Driver) DeleteObjects(ctx context.Context, prefix string) error {
	// maxKey的取值范围为[0, 1000]
	maxKeys := int64(1000)
	var objects []*awss3.Object
	var err error
	for {
		objects, err = d.ListObjects(ctx, prefix, maxKeys)
		if err != nil {
			return err
		}
		if objects == nil {
			return nil
		}
		if _, err := d.S3.DeleteObjectsWithContext(ctx, &awss3.DeleteObjectsInput{
			Bucket: aws.String(d.Bucket),
			Delete: &awss3.Delete{
				Objects: func() []*awss3.ObjectIdentifier {
					identifiers := make([]*awss3.ObjectIdentifier, 0)
					for _, obj := range objects {
						identifiers = append(identifiers, &awss3.ObjectIdentifier{
							Key: obj.Key,
						})
					}
					return identifiers
				}(),
			},
		}); err != nil {
			return err
		}
	}
}