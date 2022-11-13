package utils

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func GetRequest() {}

func POSTRequest(payload url.Values) []byte {
	request, err := http.NewRequest(
		http.MethodPost,
		"https://api.line.me/oauth2/v2.1/revoke",
		strings.NewReader(payload.Encode()),
	)
	if err != nil {
		log.Println(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		// err = fmt.Errorf("IO Read Error:: %s", err)
		log.Println(err)
	}

	defer response.Body.Close()

	log.Println(StreamToString(response.Body))

	// return response.Body
	return body
}

func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}
