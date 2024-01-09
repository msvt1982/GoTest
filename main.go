package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"

	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/Get", getPublicFile)
	r.POST("/PublicS3Files", publicS3Files)
	r.Run(":31982")
}

func getPublicFile(c *gin.Context) {
	filesource := c.Query("filesource")
	filedest := c.Query("filedest")
	//Public(filesource, filedest)
	c.JSON(200, gin.H{
		"filesource": filesource,
		"filedest":   filedest,
	})
}

func publicS3Files(c *gin.Context) {
	s3files := []struct {
		Filesource string `json:"filesource"`
		Filedest   string `json:"filedest"`
	}{}

	c.BindJSON(&s3files)

	if len(s3files) > 0 {
		ctx := context.Background()
		endpoint := "sub.domain.ru:9000"
		accessKeyID := "Fh9ofdWc4O4NK0rj9Kpq"
		secretAccessKey := "1enNMGqoz2DNljc3jojdYOBN5wAHEJ5cCT7adGWE"
		useSSL := true

		// Initialize minio client object.
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatalln(err)
		}

		// Make a new bucket called mymusic.
		bucketName := "public"
		location := "us-east-1"

		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
			if errBucketExists == nil && exists {
				log.Printf("We already own %s\n", bucketName)
			} else {
				log.Fatalln(err)
			}
		} else {
			log.Printf("Successfully created %s\n", bucketName)
		}

		log.Printf("Total elements: %d\n", len(s3files))

		for i := 0; i < len(s3files); i++ {
			log.Printf("Filesource: %s\n", s3files[i].Filesource)

			// Upload the zip file
			filePath := s3files[i].Filesource //os.Args[1]   //"MonitorERP.zip"
			objectName := s3files[i].Filedest //os.Args[2] //"golden-oldies1.zip"
			contentType := "image/jpeg"

			// Upload the zip file with FPutObject
			info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("Successfully uploaded %s of size %d. Num file %d\n", objectName, info.Size, i)
		}

	}
}
