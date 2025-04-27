package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

const openAIEndpoint = "https://api.openai.com/v1/chat/completions"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []Choice `json:"choices"`
}

type MaterialType string

const (
	ProgressUpdate MaterialType = "progress update"
	CommitMessage  MaterialType = "commit message"
)

func loadEnv() error {
	return godotenv.Load(".env")
}

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff")
	out, err := cmd.Output()
	return string(out), err
}

func sendToOpenAI(diff string, materialType MaterialType) (string, error) {
	if materialType != ProgressUpdate && materialType != CommitMessage {
		return "", fmt.Errorf("invalid material type")
	}

	err := loadEnv()
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not found in .env file")
	}

	var prompt string
	if materialType == ProgressUpdate {
		prompt = `
Based on the following git diff, generate a clear and simple summary with two sections:
1. "> What I Did Today" – list items with simple dash formatting (-), no bold or extra symbols.
2. "> Challenges and How I Solved Them" – each challenge should start with "Challenge:" and solution with "Solution:", with the solution indented under the challenge using the simple dash (-) formatting.

Keep the language non-technical and easy to understand for someone without programming knowledge.

# Important
Format the document to be like this

## Format
Things I got done today:
	- <item>

Challenges and how I solved them:
	- <challenge> -> solved by <solution>

# Note
Don't output extra things, just the list. Also don't append '#' or any markdown formatting, just the tab indentation to identify list items from the list title

Here is the git diff:
` + diff
	} else if materialType == CommitMessage {
		prompt = `
Based on the diff provided, generate a git commit message.

# Important
Do not add additional commentary texts, just the git commmit message

Here is the git diff:
` + diff
	}

	requestBody := RequestBody{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	bodyBytes, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", openAIEndpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return res.Choices[0].Message.Content, nil
}

func copyToClipboardMac(text string) {
	cmd := exec.Command("pbcopy")
	in, _ := cmd.StdinPipe()
	cmd.Start()
	in.Write([]byte(text))
	in.Close()
	cmd.Wait()
}

func main() {
	cm := flag.Bool("cm", false, "Generate commit message based on current working directory diff")
	pu := flag.Bool("pu", false, "Generate progress update in the form of a list containing 'What I did Today' and 'Challenges and How I Solved Them' as titles")
	flag.Parse()

	diff, err := getGitDiff()
	if err != nil {
		fmt.Println("Error getting git diff:", err)
		return
	}

	if len(diff) == 0 {
		fmt.Println("No changes detected.")
		return
	}

	var response string

	if *cm {
		fmt.Println("Generating commit message...")
		response, err = sendToOpenAI(diff, CommitMessage)
		if err != nil {
			fmt.Println("Error calling OpenAI API:", err)
			return
		}
	}

	if *pu {
		fmt.Println("Generating progress update...")
		response, err = sendToOpenAI(diff, ProgressUpdate)
		if err != nil {
			fmt.Println("Error calling OpenAI API:", err)
			return
		}
	}

	if response != "" {
		fmt.Println("\n--- Generated Output ---\n")
		fmt.Println(response)

		if _, err := exec.LookPath("pbcopy"); err == nil {
			copyToClipboardMac(response)
			fmt.Println("\n(Output copied to clipboard ✔️)")
		} else {
			fmt.Println("\n(Note: 'pbcopy' not found, skipped clipboard copy)")
		}
	}
}
