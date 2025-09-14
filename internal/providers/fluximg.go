package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"net/url"
	"regexp"
	"strings"
	"time"

	proxy "golang.org/x/net/proxy"
)

var re = regexp.MustCompile(`\b\d{6}\b`)

type mail struct {
	RequestID string `json:"request-id"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Expires int64  `json:"expires"`
		Inbox   []struct {
			ID         string `json:"id"`
			Inbox      string `json:"inbox"`
			SenderName string `json:"senderName"`
			From       string `json:"from"`
			Subject    string `json:"subject"`
			TextBody   string `json:"textBody"`
			ReceivedAt int64  `json:"receivedAt"`
			ReadFlag   int    `json:"readFlag"`
		} `json:"inbox"`
	} `json:"data"`
}

type FluxClient struct {
	dialer  proxy.Dialer
	domains []string
}

func NewFluxClient(proxyStr string) (*FluxClient, error) {
	parsed, err := url.Parse(proxyStr)
	if err != nil {
		return nil, err
	}
	dialer, err := proxy.FromURL(parsed, nil)
	if err != nil {
		return nil, err
	}

	return &FluxClient{dialer: dialer}, nil
}

func (flux *FluxClient) NewImage(photoBytes []byte, prompt string) ([]byte, error) {
	jar, _ := cookiejar.New(nil)
	clientWithProxy := &http.Client{
		Transport: &http.Transport{Dial: flux.dialer.Dial},
		Jar:       jar,
		Timeout:   time.Second * 10,
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: time.Second * 10,
	}

	if err := flux.register(clientWithProxy, client); err != nil {
		return []byte{}, err
	}

	return flux.generateImage(
		client,
		photoBytes, prompt,
	)
}

func (flux *FluxClient) generateImage(client *http.Client,
	photoBytes []byte, prompt string) ([]byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="image"; filename="photo.jpg"`)
	header.Set("Content-Type", "image/jpeg")

	part, err := writer.CreatePart(header)
	if err != nil {
		return []byte{}, err
	}
	_, err = io.Copy(part, bytes.NewReader(photoBytes))
	if err != nil {
		return []byte{}, err
	}

	// fields := map[string]string{
	// 	"prompt":         prompt,
	// 	"model":          "qwen-image-edit",
	// 	"aspect_ratio":   "match_input_image",
	// 	"guidance":       "10",
	// 	"steps":          "28",
	// 	"quality":        "85",
	// 	"go_fast":        "false",
	// 	"output_format":  "jpeg",
	// 	"output_quality": "100",
	// }

	// for key, value := range fields {
	// 	if err := writer.WriteField(key, value); err != nil {
	// 		return []byte{}, err
	// 	}
	// }

	fields := map[string]string{
		"prompt":        prompt,
		"model":         "nano-banana",
		"output_format": "jpg",
	}

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return []byte{}, err
		}
	}

	if err := writer.Close(); err != nil {
		return []byte{}, err
	}

	resp, err := client.Post(
		"https://flux-1.net/api/tools/qwen-image-edit",
		writer.FormDataContentType(),
		&buf,
	)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != 200 {
		return []byte{}, errors.New("you have been blocked")
	}

	var result struct {
		Output string `json:"output"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return []byte{}, err
	}
	resp.Body.Close()

	resp, err = client.Get(result.Output)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (flux *FluxClient) getEmail(client *http.Client) (string, error) {
	baseURL := "https://tempmail.so/us/api/inbox"

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result mail
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.Name, nil
}

func (flux *FluxClient) getCode(client *http.Client) (string, error) {
	var result mail
	baseURL := "https://tempmail.so/us/api/inbox"

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")

	for _ = range 3 {
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return "", err
		}

		if len(result.Data.Inbox) == 0 {
			time.Sleep(time.Second * 3)
			continue
		}

		mail := result.Data.Inbox[0]
		if mail.From == `"no-reply@flux-1.net"` {
			text := mail.TextBody
			matches := re.FindStringSubmatch(text)
			if len(matches) == 0 {
				return "", fmt.Errorf("verification code not found")
			}

			return matches[0], nil
		}
	}

	return "", err
}

func (flux *FluxClient) register(clientWithProxy *http.Client, client *http.Client) error {
	email, err := flux.getEmail(client)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(map[string]string{
		"email": email,
	})
	if err != nil {
		return err
	}

	// Sending auth code
	resp, err := clientWithProxy.Post(
		"https://flux-1.net/api/auth/send-code",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("you have been blocked")
	}

	// Get csrf token
	resp, err = client.Get("https://flux-1.net/api/auth/csrf")
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("you have been blocked")
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	resp.Body.Close()

	csrfToken, ok := result["csrfToken"]
	if !ok || csrfToken == "" {
		return errors.New("failed to get csrfToken")
	}

	// Get auth code
	code, err := flux.getCode(client)

	if err != nil || code == "" {
		return errors.New("failed to get auth code")
	}

	// Login
	data := url.Values{}
	data.Set("email", email)
	data.Set("code", code)
	data.Set("redirect", "false")
	data.Set("callbackUrl", "/")
	data.Set("csrfToken", csrfToken)

	_, err = client.Post(
		"https://flux-1.net/api/auth/callback/credentials",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("you have been blocked")
	}

	return nil
}
