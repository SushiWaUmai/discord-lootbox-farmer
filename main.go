package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type RateLimitResponse struct {
	Message    string  `json:"message"`
	RetryAfter float64 `json:"retry_after"`
	Global     bool    `json:"global"`
}

var (
	AUTHORIZATION      string
	X_SUPER_PROPERTIES string
)

func main() {
	loadEnv()

	client := &http.Client{}
	url := "https://discord.com/api/v9/users/@me/lootboxes/open"

	for {
		resp, err := makeRequest(client, url)
		defer resp.Body.Close()

		if err != nil {
			log.Println("Error:", err)
			continue
		}
		log.Printf("POST request sent. Status: %s\n", resp.Status)

		if resp.StatusCode == http.StatusTooManyRequests {
			var limitResp RateLimitResponse
			body, _ := io.ReadAll(resp.Body)
			if err := json.Unmarshal(body, &limitResp); err != nil {
				log.Println("Error parsing rate limit response:", err)
				continue
			}

			log.Printf("Rate limited. Retrying after %.2f seconds\n", limitResp.RetryAfter)
			time.Sleep(time.Duration(limitResp.RetryAfter * float64(time.Second)))
		} else {
			time.Sleep(2 * time.Second)
		}
	}
}

func makeRequest(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", AUTHORIZATION)
	req.Header.Set("X-Super-Properties", X_SUPER_PROPERTIES)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return resp, nil
}

func loadEnv() {
	AUTHORIZATION = os.Getenv("AUTHORIZATION")
	if AUTHORIZATION == "" {
		log.Fatalln("Missing env property: \"AUTHORIZATION\"")
	}

	X_SUPER_PROPERTIES = os.Getenv("X_SUPER_PROPERTIES")
	if AUTHORIZATION == "" {
		log.Fatalln("Missing env property: \"X_SUPER_PROPERTIES\"")
	}
}
