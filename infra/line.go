package infra

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type AccessTokenMeta struct {
	Token string `json:"access_token"`
	Type  string `json:"token_type"`
	Exp   int64  `json:"expires_in"`
	Id    string `json:"key_id"`
}

func fetchJwt() string {
	// Open private key file
	file, err := os.Open(os.Getenv("LINE_PRIVATE_KEY_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read private contents from the opened file
	buffer, err := ioutil.ReadAll(file)

	if err != nil {
		log.Printf("Failed to read the privatekey: %s\n", err)
	}
	privkey, err := jwk.ParseKey(buffer)
	if err != nil {
		log.Printf("Failed to parse the privatekey to JWK format: %s\n", err)
	}

	// Build JWT payload
	// Refer: https://developers.line.biz/ja/docs/messaging-api/generate-json-web-token/#generate-jwt
	token, err := jwt.NewBuilder().
		Subject(os.Getenv("LINE_CHANNEL_ID")).
		Issuer(os.Getenv("LINE_CHANNEL_ID")).
		Audience([]string{"https://api.line.me/"}).
		Expiration(time.Now().Add(30*time.Minute)).
		Claim("token_exp", 60*60*24*30).
		Build()
	if err != nil {
		fmt.Printf("failed to build token: %s\n", err)
	}

	// Generate JWT (variable signed is a jwt)
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, privkey))
	if err != nil {
		fmt.Printf("Failed to sign token: %s\n", err)
	}
	return string(signed)
}

func FetchChannelAccessToken() string {
	// Referrence: https://developers.line.biz/ja/reference/messaging-api/#issue-channel-access-token-v2-1
	payload := url.Values{}
	payload.Set("grant_type", "client_credentials")
	payload.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	payload.Add("client_assertion", string(fetchJwt()))

	reqBody := strings.NewReader(payload.Encode())
	req, err := http.NewRequest(http.MethodPost, "https://api.line.me/oauth2/v2.1/token", reqBody)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var resBody io.Reader = res.Body
	var access_token_meta AccessTokenMeta
	err = json.NewDecoder(resBody).Decode(&access_token_meta)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(access_token_meta.Token)
	return access_token_meta.Token
}

func RevokeAccessToken(accessToken string) {
	payload := url.Values{}
	payload.Set("client_id", os.Getenv("LINE_CHANNEL_ID"))
	payload.Add("client_secret", os.Getenv("LINE_CHANNEL_SECRET"))
	payload.Add("access_token", accessToken)

	reqBody := strings.NewReader(payload.Encode())
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.line.me/oauth2/v2.1/revoke",
		reqBody,
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	log.Println(res.StatusCode)
}
