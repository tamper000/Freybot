package transcribe

type STTResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Transcription []struct {
			Caption string `json:"caption"`
		} `json:"transcription"`
	} `json:"data"`
}
