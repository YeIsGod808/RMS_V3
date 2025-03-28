package minioStore

import (
	"fmt"
	"strings"
	"sync"

	"github.com/RMS_V3/config"
	"github.com/RMS_V3/log"
	"github.com/RMS_V3/middleware/snowflake"
	"github.com/minio/minio-go/v6"
)

var (
	client         Minio
	minioStoreOnce sync.Once
)

type Minio struct {
	MinioClient  *minio.Client
	Endpoint     string
	Port         string
	VideoBuckets string
	PicBuckets   string
	FileBuckets  string
}

func GetMinio() Minio {
	minioStoreOnce.Do(initMinio)
	return client
}
func initMinio() {
	conf := config.GetGlobalConfig().MinioConfig

	// 服务端endpoint
	endpoint := conf.Host
	port := conf.Port
	endpoint = endpoint + ":" + port

	accessKeyID := conf.AccessKeyID
	secretAccessKey := conf.SecretAccessKey

	// bucket name
	videoBucket := conf.VideoBuckets
	picBucket := conf.PicBuckets
	fileBuckets := conf.FileBuckets

	// 初使化 minio client对象。
	minioClient, err := minio.New(
		endpoint,
		accessKeyID,
		secretAccessKey,
		false,
	)
	if err != nil {
		panic(err)
	}

	// 创建存储桶
	creatBucket(minioClient, videoBucket)
	creatBucket(minioClient, picBucket)
	creatBucket(minioClient, fileBuckets)
	client = Minio{minioClient, endpoint, port, videoBucket, picBucket, fileBuckets}
}
func creatBucket(minioClient *minio.Client, bucketName string) {
	found, err := minioClient.BucketExists(bucketName)
	if err != nil {
		log.Errorf("check %s bucketExists err:%s", bucketName, err)

	}

	if !found {
		minioClient.MakeBucket(bucketName, "us-east-1")
	}

	policy := `{"Version": "2012-10-17",
				"Statement": 
					[{
						"Action":["s3:GetObject"],
						"Effect": "Allow",
						"Principal": {"AWS": ["*"]},
						"Resource": ["arn:aws:s3:::` + bucketName + `/*"],
						"Sid": ""
					}]
				}`
	err = minioClient.SetBucketPolicy(bucketName, policy)
	if err != nil {
		log.Errorf("set %s bucket policy err:%s", bucketName, err)
	}
	log.Infof("create bucket %s success", bucketName)
}
func (m *Minio) UploadFile(filetype, file, userID string) (string, string, error) {
	var fileName strings.Builder
	var contentType, suffix, bucket string

	// 根据filetype设置contentType、suffix和bucket
	switch filetype {
	case "video":
		contentType = "video/mp4"
		suffix = ".mp4"
		bucket = m.VideoBuckets
	case "pdf":
		contentType = "application/pdf"
		suffix = ".pdf"
		bucket = m.FileBuckets
	case "pptx":
		contentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
		suffix = ".pptx"
		bucket = m.FileBuckets
	case "word":
		contentType = "application/msword" // 对于.doc文件；如果是.docx，则使用 "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		suffix = ".doc"
		bucket = m.FileBuckets
	case "image":
		contentType = "image/jpeg" // 也可以是其他图片格式如 "image/png"
		suffix = ".jpg"
		bucket = m.PicBuckets
	default:
		return "", "", fmt.Errorf("unsupported file type: %s", filetype)
	}

	// 通过用户id来保证不同用户即使在同一时刻上传到同一台服务器上的文件也会不一样
	fileName.WriteString(userID)
	fileName.WriteString("_")

	// 根据机器码、时间戳等生产分布式id，确保每个上传到MinIO/S3的对象都有一个独一无二的名字。
	snowflakeID := snowflake.GenID()

	fileName.WriteString(snowflakeID)
	fileName.WriteString(suffix)

	// 从本地路径file中取文件上传到minio服务器的特定bucket中，并将其在minio服务器上的名字保存为fileName
	n, err := m.MinioClient.FPutObject(bucket, fileName.String(), file, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Errorf("upload file error:%s", err.Error())
		return "", "", err
	}
	log.Infof("upload file %d byte success,fileName:%s", n, fileName.String())

	url := "http://" + m.Endpoint + "/" + bucket + "/" + fileName.String()
	// 返回一个url链接，指向上传的文件
	return url, fileName.String(), nil
}

// func (m *Minio) UploadFile(filetype, file, userID string) (string, string, error) {
// 	var fileName strings.Builder
// 	var contentType, Suffix, bucket string

// 	if filetype == "video" {
// 		contentType = "video/mp4"
// 		Suffix = ".mp4"
// 		bucket = m.VideoBuckets
// 	} else if filetype == "file" {
// 		contentType = "application/pdf"
// 		Suffix = ".pdf"
// 		bucket = m.FileBuckets
// 	} else {
// 		contentType = "image/jpeg"
// 		Suffix = ".jpg"
// 		bucket = m.PicBuckets
// 	}
// 	// 通过用户id来保证不同用户即使在同一时刻上传到同一台服务器上的文件也会不一样
// 	fileName.WriteString(userID)
// 	fileName.WriteString("_")

// 	// 根据机器码、时间戳等生产分布式id，确保每个上传到MinIO/S3的对象都有一个独一无二的名字。这避免了因文件名重复而导致的数据覆盖问题。
// 	snowflakeID := snowflake.GenID()

// 	fileName.WriteString(snowflakeID)
// 	fileName.WriteString(Suffix)
// 	// 从本地路径file中取文件上传到minio服务器的特定bucket中，并将其在minio服务器上的名字保存为fileName
// 	n, err := m.MinioClient.FPutObject(bucket, fileName.String(), file, minio.PutObjectOptions{ContentType: contentType})
// 	if err != nil {
// 		log.Errorf("upload file error:%s", err.Error())
// 		return "", "", err
// 	}
// 	log.Infof("upload file %d byte success,fileName:%s", n, fileName)
// 	url := "http://" + m.Endpoint + "/" + bucket + "/" + fileName.String()
// 	// 返回一个url链接，指向上传的文件
// 	return url, fileName.String(), nil
// }
