package transcribe

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

const url = "https://anytranscribe.com/wp-admin/admin-ajax.php"

func TranscribeAudio(audioData []byte) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fields := map[string]string{
		"action":                  "audio_transcription_generate",
		"audio_transcriber_nonce": "be8aaa40c9",
		"language":                "undefined",
	}

	for key, value := range fields {
		err := writer.WriteField(key, value)
		if err != nil {
			return "", err
		}
	}

	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", `form-data; name="audio_file"; filename="blob"`)
	header.Set("Content-Type", "audio/ogg")
	part, err := writer.CreatePart(header)
	if err != nil {
		return "", err
	}
	_, err = part.Write(audioData)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://anytranscribe.com/")
	req.Header.Set("Origin", "https://anytranscribe.com")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=4")
	req.Header.Set("Content-Type", writer.FormDataContentType()) // Важно: boundary

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Bad status code")
	}

	var stt STTResponse
	if err := json.NewDecoder(resp.Body).Decode(&stt); err != nil {
		return "", nil
	}

	if !stt.Success || len(stt.Data.Transcription) == 0 {
		return "", errors.New("Failed to process response")
	}

	var text string

	for _, caption := range stt.Data.Transcription {
		text += caption.Caption
		text += "\n"
	}

	return text, nil
}
