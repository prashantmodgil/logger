package minio
import (
	"github.com/minio/minio-go/v6"
	"log"
	"io/ioutil"
	"fmt"
		"bytes"
		"os"


)
var MinioClient *minio.Client
func InitializeMinio(){
	// from env
	endpoint :=os.Getenv("endpoint")
	accessKeyID := os.Getenv("accessKeyID")
	secretAccessKey := os.Getenv("secretAccessKey")
	useSSL := true

	var err error
	MinioClient, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}
}

func UploadLogToMinio(logFileName, date string, hour int,) error {
    _, err := MinioClient.StatObject("mylogs", fmt.Sprintf("%v/log_%d.txt", date,hour), minio.StatObjectOptions{})

    if err == nil {
        existingLogs, err := MinioClient.GetObject("mylogs", fmt.Sprintf("%v/log_%d.txt", date,hour), minio.GetObjectOptions{})
        if err != nil {
            return err
        }
        defer existingLogs.Close()

        existingLogBytes, err := ioutil.ReadAll(existingLogs)
        if err != nil {
            return err
        }

        newLogBytes, err := ioutil.ReadFile(logFileName)
        if err != nil {
            return err
        }
        updatedLogBytes := append(existingLogBytes, newLogBytes...)

        _, err = MinioClient.PutObject("mylogs", fmt.Sprintf("%v/log_%d.txt", date,hour), bytes.NewReader(updatedLogBytes), int64(len(updatedLogBytes)), minio.PutObjectOptions{})
        if err != nil {
            return err
        }

    } else {
        _, err := MinioClient.FPutObject("mylogs", fmt.Sprintf("%v/log_%d.txt", date,hour), logFileName, minio.PutObjectOptions{})
        if err != nil {
            return err
        }
    }

    return nil
}
