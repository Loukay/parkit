package main

import (
	"context"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

// Controller has handlers for the API requests
type Controller struct {
	UploaderManager *manager.Uploader
	Bucket          string
}

// UploadImage uploads an image of the parking lot to the S3 bucket
func (controller Controller) UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	buffer, err := file.Open()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	defer buffer.Close()

	fileName := strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"

	go controller.uploadToBucket(fileName, buffer)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Image sent successfully"})

}

func (controller Controller) uploadToBucket(fileName string, file io.Reader) {
	result, err := controller.UploaderManager.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(controller.Bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		log.Printf("Error while uploading: %v", err)
		return
	}

	log.Printf("Image uploaded successfully: %+v", result)
}
