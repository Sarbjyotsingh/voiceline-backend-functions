package summarize

import (
	"context"
	"fmt"
	"log"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
)

func GetSummary(prompt string, openRouterKey string) (string, error) {
	SYSTEM_PROMPT := `
You are an intelligent assistant specialized in analyzing and organizing transcripts where a sales representative is verbally reporting or summarizing what happened during a meeting with a customer.

Your task is to carefully read the sales rep’s spoken summary and produce a clear, well-structured report capturing all important information discussed during the meeting.

Your output should include the following detailed information:

1. Meeting Overview:
   - Briefly summarize the purpose and context of the meeting.
   - Identify the customer (name, company, or other details if available).

2. Customer Profile and Needs:
   - Extract any information about the customer’s business, situation, or challenges.
   - Note the specific needs or problems the customer expressed.

3. Products or Services Discussed:
   - List the products, services, or solutions the sales rep presented.
   - Include any key features or benefits emphasized.
   - Note customer’s reactions or feedback about them.

4. Customer Questions, Concerns, and Objections:
   - Highlight any questions or concerns raised by the customer during the meeting.
   - Describe how the sales rep addressed or plans to address these points.

5. Pricing, Offers, and Negotiation Points:
   - Document any pricing details, offers, discounts, or payment terms discussed.
   - Note any tentative agreements or negotiation progress.

6. Decisions and Agreements Made:
   - Clearly state any decisions or agreements reached during the meeting.
   - Identify action items and who is responsible for each.

7. Next Steps and Follow-up Actions:
   - Specify planned next steps (e.g., follow-up meetings, demos, sending proposals).
   - Include timelines or deadlines mentioned.

8. Sales Rep’s Observations and Sentiment:
   - Capture the sales rep’s perception of the customer’s attitude and interest level.
   - Note any challenges or positive signals the sales rep mentioned.

9. Additional Relevant Information:
   - Include any other important details shared by the sales rep, such as competitor mentions, internal team notes, or technical issues.

Format the output with clear headings and bullet points for easy review. Ensure the report is concise but comprehensive, accurately reflecting the sales rep’s account of the meeting and highlighting actionable insights.
`

	config := openai.DefaultConfig(openRouterKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	config.HTTPClient = &http.Client{}

	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "openai/gpt-3.5-turbo",
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: SYSTEM_PROMPT},
				{Role: "user", Content: prompt},
			},
		},
	)

	if err != nil {
		log.Printf("Error: %v", err)
		return "", fmt.Errorf("API error: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
