package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	tg "github.com/galeone/tfgo"
	"github.com/gofiber/fiber/v2"
)

type classification struct {
	Label      string
	Proability float32
}

var labels []string = []string{"Occupied", "Vacant"}

func main() {

	model := tg.LoadModel("/usr/share/detector/models/nadi/", []string{"serve"}, nil)

	awsCfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(awsCfg)

	uploadManager := manager.NewUploader(client)

	controller := Controller{
		UploaderManager: uploadManager,
		Bucket:          os.Getenv("S3_BUCKET_NAME"),
		Model:           model,
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
