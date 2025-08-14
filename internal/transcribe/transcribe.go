package transcribe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

var (
	client    = &http.Client{}
	urlUpload = "https://transcribetext.com/upload_async"
	urlStatus = "https://transcribetext.com/task_status/"
)

var fields = map[string]string{
	"transcription_engine": "groq",
	"timestamp_format":     "semantic",
	"enable_alignment":     "true",
	"enable_diarization":   "false",
}

type CreatedJob struct {
	JobID   string `json:"job_id"`
	Success bool   `json:"success"`
}

type StatusJob struct {
	Result struct {
		Success       bool   `json:"success"`
		Transcription string `json:"transcription"`
	} `json:"result"`
	Status string `json:"status"`
}

func TranscribeAudio(fileData []byte) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fileWriter, err := writer.CreateFormFile("file", "data.ogg")
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}
	_, err = fileWriter.Write(fileData)
	if err != nil {
		return "", fmt.Errorf("error writing field: %w", err)
	}
	for key, value := range fields {
		err = writer.WriteField(key, value)
		if err != nil {
			return "", fmt.Errorf("error writing file data: %w", err)
		}
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing writer: %w", err)
	}

	req, err := http.NewRequest("POST", urlUpload, &body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Origin", "https://transcribetext.com")
	req.Header.Set("Referer", "https://transcribetext.com/")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=0")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server internal error: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var result CreatedJob
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("error get job_id")
	}
	jobID := result.JobID

	for range 10 {
		time.Sleep(1 * time.Second)

		statusReq, err := http.NewRequest("GET", fmt.Sprintf("%s%s", urlStatus, jobID), nil)
		if err != nil {
			return "", fmt.Errorf("error creating request: %w", err)
		}

		statusReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0")
		statusReq.Header.Set("Accept", "*/*")
		statusReq.Header.Set("Accept-Language", "en-US,en;q=0.5")
		statusReq.Header.Set("Referer", "https://transcribetext.com/")
		statusReq.Header.Set("Origin", "https://transcribetext.com")
		statusReq.Header.Set("Sec-Fetch-Site", "same-origin")
		statusReq.Header.Set("Sec-Fetch-Mode", "cors")

		statusResp, err := client.Do(statusReq)
		if err != nil {
			return "", fmt.Errorf("error sending request: %w", err)
		}

		statusBody, _ := io.ReadAll(statusResp.Body)
		statusResp.Body.Close()

		var statusData StatusJob
		json.Unmarshal(statusBody, &statusData)

		if statusData.Status == "completed" && statusData.Result.Success {
			return statusData.Result.Transcription, nil
		}
	}

	return "", fmt.Errorf("transcription not found")
}
