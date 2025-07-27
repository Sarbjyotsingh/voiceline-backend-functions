package transcribe

import (
	"context"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func TranscribeAudio(apiKey, audioFilePath string) (string, error) {
	// Set a timeout (e.g., 2 minutes)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client := openai.NewClient(apiKey)

	resp, err := client.CreateTranscription(ctx, openai.AudioRequest{
		Model:    "whisper-1",
		FilePath: audioFilePath,
	})
	if err != nil {
		return "", err
	}

	return resp.Text, nil
}
