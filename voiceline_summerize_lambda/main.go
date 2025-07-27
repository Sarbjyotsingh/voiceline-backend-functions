package main

import (
	"context"
	"log"
	"os"
	"voiceline_summerize_lambda/summarize"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	apiKey := os.Getenv("OPEN_ROUTER_API_KEY")
	router := gin.Default()
	router.GET("/voicelinetest", func(c *gin.Context) {
		text := c.Query("text")

		if text == "" {
			c.JSON(400, gin.H{
				"error": "Text parameter is required",
			})
			return
		}
		summary, err := summarize.GetSummary(text, apiKey)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Failed to generate summary",
			})
			return
		}

		c.JSON(200, gin.H{
			"received": summary,
		})
	})

	ginLambda = ginadapter.New(router)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	log.Printf("Lambda function started")
	lambda.Start(handler)
}
