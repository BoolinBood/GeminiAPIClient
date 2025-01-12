package ai

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	API_URL = "https://api.openai.com/v1/audio/speech"
)

func OpenAiTextToSpeech(textToSpeech string) {
	request, err := http.NewRequest("POST", API_URL, nil)

	if err != nil {
		log.Println(err)
	}

	request.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	request.Header.Add("Content-Type", "application/json")

	data := map[string]string{
		"model": "tts-1",
		"input": textToSpeech,
		"voice": "sage",
	}
	requestBody, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
	request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	log.Println(response.StatusCode)
}
