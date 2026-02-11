package pollinations

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	http "github.com/bogdanfinn/fhttp"

	"github.com/aandrew-me/tgpt/v2/src/client"
	"github.com/aandrew-me/tgpt/v2/src/structs"
)

func NewRequest(input string, params structs.Params) (*http.Response, error) {
	client, err := client.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	model := "openai-large"
	if params.ApiModel != "" {
		model = params.ApiModel
	}

	temperature := 0.6
	if params.Temperature != "" {
		if parsedTemp, err := strconv.ParseFloat(params.Temperature, 64); err == nil {
			temperature = parsedTemp
		}
	}

	top_p := 1.0
	if params.Top_p != "" {
		if parsedTopP, err := strconv.ParseFloat(params.Top_p, 64); err == nil {
			top_p = parsedTopP
		}
	}

	type pollinationsMessage struct {
		Content string `json:"content"`
		Role    string `json:"role"`
	}

	type pollinationsRequest struct {
		Messages    []pollinationsMessage `json:"messages"`
		Model       string               `json:"model"`
		Stream      bool                 `json:"stream"`
		Temperature float64              `json:"temperature"`
		TopP        float64              `json:"top_p"`
		Referrer    string               `json:"referrer"`
	}

	messages := []pollinationsMessage{
		{
			Content: params.SystemPrompt,
			Role:    "system",
		},
		{
			Content: input,
			Role:    "user",
		},
	}

	reqBody := pollinationsRequest{
		Messages:    messages,
		Model:       model,
		Stream:      true,
		Temperature: temperature,
		TopP:        top_p,
		Referrer:    "tgpt",
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("\nFailed to marshal request body.")
		fmt.Println("Error:", err)
		os.Exit(0)
	}

	var data = strings.NewReader(string(bodyBytes))

	req, err := http.NewRequest("POST", "https://text.pollinations.ai/openai", data)
	if err != nil {
		fmt.Println("\nSome error has occurred.")
		fmt.Println("Error:", err)
		os.Exit(0)
	}
	// Setting all the required headers
	req.Header.Set("Content-Type", "application/json")

	// Return response
	return (client.Do(req))
}

func GetMainText(line string) (mainText string) {
	var obj = "{}"
	if len(line) > 1 {
		obj = strings.Split(line, "data: ")[1]
	}

	var d structs.CommonResponse
	if err := json.Unmarshal([]byte(obj), &d); err != nil {
		return ""
	}

	if len(d.Choices) > 0 {
		mainText = d.Choices[0].Delta.Content
		return mainText
	}
	return ""
}
