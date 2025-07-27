package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"voiceline_process_audio_lambda/transcribe"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func UploadAudioHandler(c *gin.Context) {

	// 1. Read base64 string from form field "audio_file"
	base64Data := c.PostForm("audio_file")
	if base64Data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Base64 audio data is required"})
		return
	}

	// 2. Remove data URL prefix if present (e.g. "data:audio/mp3;base64,")
	if idx := strings.Index(base64Data, ","); idx != -1 {
		base64Data = base64Data[idx+1:]
	}

	// 3. Decode base64 string
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 audio data"})
		return
	}

	// 4. Ensure /tmp/uploads directory exists
	uploadDir := "/tmp/uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads folder"})
			return
		}
	}

	// 5. Save decoded bytes to a file, e.g. audio_decoded.mp3
	savePath := fmt.Sprintf("%s/%s", uploadDir, "audio_decoded.mp3")
	err = os.WriteFile(savePath, decodedData, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save audio"})
		return
	}

	// 6. Get OPENAI_API_KEY from env (make sure it's set in Lambda environment variables)
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OpenAI API key is not set"})
		return
	}

	// 7. Call TranscribeAudio function to get the transcription text
	transcript, err := transcribe.TranscribeAudio(apiKey, savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to transcribe audio: %v", err)})
		return
	}

	// 8. Respond success with saved file path and transcription text
	c.JSON(http.StatusOK, gin.H{
		"message":    "Audio uploaded, decoded and transcribed successfully",
		"filePath":   savePath,
		"transcript": transcript,
	})
}

func init() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.POST("/voicelinemain", UploadAudioHandler)

	ginLambda = ginadapter.New(router)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
