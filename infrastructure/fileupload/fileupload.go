package fileupload

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type fileUpload struct{}

func NewFileUpload() *fileUpload {
	return &fileUpload{}
}

type UploadFileInterface interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}

//So what is exposed is Uploader
var _ UploadFileInterface = &fileUpload{}

func (f fileUpload) UploadFile(file *multipart.FileHeader) (string, error) {
	fileOpen, err := file.Open()
	if err != nil {
		return "", errors.New("cannot open file")
	}

	defer fileOpen.Close()

	size := file.Size
	//The image should not be more than 500KB
	fmt.Println("the size: ", size)
	if size > int64(512000) {
		return "", errors.New("sorry, please upload an Image of 500KB or less")
	}

	//only the first 512 bytes are used to sniff the content type of a file,
	//so, so no need to read the entire bytes of a file.
	buffer := make([]byte, size)

	_, err = fileOpen.Read(buffer)
	if err != nil {
		return "", err
	}

	fileType := http.DetectContentType(buffer)

	if !strings.HasPrefix(fileType, "image") {
		return "", errors.New("please upload a valid image")
	}

	fileBytes := bytes.NewReader(buffer)
	filePath := FormatFile(file.Filename)
	path := "/profile-photos/" + filePath

	params := &s3.PutObjectInput{
		Bucket:        aws.String("chodapi"), //this is the name i saved the bucket that contains the image
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
		ACL:           aws.String("public-read"),
	}

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("DO_SPACES_KEY"), os.Getenv("DO_SPACES_SECRET"), os.Getenv("DO_SPACES_TOKEN")),
		Endpoint:    aws.String(os.Getenv("DO_SPACES_ENDPOINT")),
		Region:      aws.String(os.Getenv("DO_SPACES_REGION")),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return "", err
	}

	s3Client := s3.New(newSession)

	_, err = s3Client.PutObject(params)
	if err != nil {
		fmt.Println("actual error: ", err)
		return "", errors.New("something went wrong uploading image")
	}

	return filePath, nil
}
