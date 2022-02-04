package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

func main() {

	awsCfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(awsCfg)

	uploadManager := manager.NewUploader(client)

	controller := Controller{
		UploaderManager: uploadManager,
		Bucket:          os.Getenv("S3_BUCKET_NAME"),
	}

	var app *fiber.App = fiber.New(fiber.Config{
		Prefork: false,
	})

	app.Post("/images", controller.UploadImage)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Failed to listen to web server")
		panic(err)
	}

}
