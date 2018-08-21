package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func unmarshalPayload(data []byte) (payload, error) {
	var r payload
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *payload) marshal() ([]byte, error) {
	return json.Marshal(r)
}

type payload struct {
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	Fallback string   `json:"fallback"`
	Pretext  *string  `json:"pretext"`
	Color    string   `json:"color"`
	Fields   []field  `json:"fields"`
	Actions  []action `json:"actions"`
}

type action struct {
	Type  string  `json:"type"`
	Text  string  `json:"text"`
	URL   *string `json:"url"`
	Style string  `json:"style,omitempty"`
}

type field struct {
	Title string  `json:"title"`
	Value *string `json:"value"`
	Short bool    `json:"short"`
}

var (
	url     = os.Getenv("SLACK_WEBHOOK_URL")
	app     = flag.String("app", "REQUIRED", "Pass the app name")
	env     = flag.String("env", "REQUIRED", "Build ENV")
	message = flag.String("message", "REQUIRED", "Usually the commit message")
	jobURL  = flag.String("job_url", "REQUIRED", "URL to the build")
	appURL  = flag.String("app_url", "REQUIRED", "URL to open the app")
)

func main() {
	if len(url) == 0 {
		log.Fatal("Set SLACK_WEBHOOK_URL")
	}

	flag.Parse()

	p := &payload{
		Attachments: []attachment{
			attachment{
				Fallback: "Release",
				Pretext:  message,
				Color:    "#764FA5",
				Fields: []field{
					field{
						Title: "App",
						Value: app,
						Short: true,
					},
					field{
						Title: "Env",
						Value: env,
						Short: true,
					},
				},
				Actions: []action{
					action{
						Type:  "button",
						Text:  "View Build 👾",
						URL:   jobURL,
						Style: "primary",
					},
					action{
						Type: "button",
						Text: "Open App 🚀",
						URL:  appURL,
					},
				},
			},
		},
	}

	b, err := p.marshal()
	if err != nil {
		log.Fatal("marshall", err)
	}
	data := bytes.NewBuffer(b)
	req, _ := http.NewRequest("POST", url, data)

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("http do", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
