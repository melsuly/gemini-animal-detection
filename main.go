package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file")
	}
}

func main() {
	ApiKey := os.Getenv("API_KEY")
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(ApiKey))
	if err != nil {
		logrus.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	model.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
	}

	prompt := `Is this image contain any animal or not? Response should be in JSON format: { "isDetected": "boolean" }`

	imageFile, err := os.Open("example/bus.jpeg")
	if err != nil {
		logrus.Fatal(err)
	}
	defer imageFile.Close()

	opts := genai.UploadFileOptions{DisplayName: "Image to Handle"}

	doc, err := client.UploadFile(ctx, "", imageFile, &opts)
	if err != nil {
		logrus.Fatal(err)
	}
	defer client.DeleteFile(ctx, doc.Name)

	resp, err := model.GenerateContent(ctx,
		genai.FileData{URI: doc.URI},
		genai.Text(prompt))
	if err != nil {
		logrus.Fatal(err)
	}

	if len(resp.Candidates) > 0 {
		c := resp.Candidates[0]
		if len(c.Content.Parts) > 0 {
			fmt.Println(c.Content.Parts[0])
		}
	}
}
