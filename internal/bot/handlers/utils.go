package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tamper000/freybot/internal/config"
	"github.com/tamper000/freybot/internal/models"
	"github.com/tamper000/freybot/internal/providers"
	"golang.org/x/net/html"
)

func GetModelsByGroup(group string) []config.AIModel {
	return config.Models[config.ModelGroup(group)]
}

func GetModelByCallback(group string, callback string) config.AIModel {
	modelGroup := config.Models[config.ModelGroup(group)]

	for _, info := range modelGroup {
		if info.CallbackData == callback {
			return info
		}
	}

	return config.AIModel{}
}

func GetProviderByUser(user *models.User) config.ProviderModel {
	modelGroup := config.Models[config.ModelGroup(user.Group)]

	for _, info := range modelGroup {
		if info.ApiName == user.Model {
			return info.Provider
		}
	}

	return ""
}

func GetModelByUser(user *models.User) config.AIModel {
	modelGroup := config.Models[config.ModelGroup(user.Group)]

	for _, info := range modelGroup {
		if info.ApiName == user.Model {
			return info
		}
	}

	return config.AIModel{}
}

func IsSupportPhotoByUser(user *models.User) bool {
	modelGroup := config.Models[config.ModelGroup(user.Group)]

	for _, info := range modelGroup {
		if info.ApiName == user.Model {
			return info.Image
		}
	}

	return false
}

func GetPhotoModelByApiName(name string) config.AIModel {
	for _, info := range config.PhotoModels {
		if info.ApiName == name {
			return info
		}
	}

	return config.AIModel{}
}

func splitText(text string) []string {
	var (
		result    = []string{}
		chunkSize = 3600
		current   = ""
	)

	splitted := strings.Split(text, "\n")

	for _, chunk := range splitted {
		if len(current)+len(chunk)+len("\n") > chunkSize {
			result = append(result, current)
			current = chunk
		} else {
			current += "\n" + chunk
		}
	}

	if len(result) == 0 {
		result = append(result, current)
	} else if result[len(result)-1] != current {
		result = append(result, current)
	}

	return result
}

func SplitHTML(content string) ([]string, error) {
	reader := strings.NewReader(content)
	tokenizer := html.NewTokenizer(reader)

	var (
		current   strings.Builder
		intag     strings.Builder
		result    = []string{}
		chunkSize = 3600
		opens     = 0
		isCode    = false
		clearText strings.Builder
	)

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			if opens != 0 {
				return result, errors.New("Нейросеть вернула битый текст.")
			}
			if current.Len() > 0 {
				result = append(result, current.String())
			}
			if intag.Len() > 0 {
				result = append(result, intag.String())
			}
			break
		}

		switch tokenType {
		case html.StartTagToken:
			token := tokenizer.Token()
			text := token.String()
			fmt.Println(text, isCode)
			if token.Data == "code" {
				isCode = true
				intag.WriteString(text)
				opens++
			} else {
				if isCode {
					text = strings.ReplaceAll(text, "<", "&lt;")
					text = strings.ReplaceAll(text, ">", "&gt;")
					intag.WriteString(text)
				} else {
					intag.WriteString(text)
					opens++
				}
			}
		case html.EndTagToken:
			token := tokenizer.Token()
			text := token.String()

			fmt.Println(text, isCode)
			if token.Data == "code" {
				isCode = false
				intag.WriteString(text)
				opens--
			} else {
				if isCode {
					text = strings.ReplaceAll(text, "<", "&lt;")
					text = strings.ReplaceAll(text, ">", "&gt;")
					intag.WriteString(text)
				} else {
					intag.WriteString(text)
					opens--
				}
			}
		default:
			text := tokenizer.Token().String()
			clearText.WriteString(text)
			text = strings.ReplaceAll(text, "<", "&lt;")
			text = strings.ReplaceAll(text, ">", "&gt;")
			intag.WriteString(text)
			if intag.Len() > chunkSize && opens != 0 {
				return result, errors.New("Нейросеть вернула битый текст.")
			}
		}

		if opens == 0 {
			if intag.Len()+current.Len() <= chunkSize {
				current.WriteString(intag.String())
				intag.Reset()
			} else if intag.Len() > chunkSize {
				result = append(result, current.String())
				current.Reset()
				chunks := splitText(intag.String())
				fmt.Println(len(chunks))
				result = append(result, chunks[:len(chunks)-1]...)
				intag.Reset()
				current.WriteString(chunks[len(chunks)-1])
			}
		}
	}

	return result, nil
}

func (h *Handler) TranscribeAudio(ctx context.Context, bot *telego.Bot, voice *telego.Voice) (string, error) {
	file, err := bot.GetFile(ctx, &telego.GetFileParams{
		FileID: voice.FileID,
	})
	if err != nil {
		return "", err
	}

	url := bot.FileDownloadURL(file.FilePath)
	fileData, err := tu.DownloadFile(url)
	if err != nil {
		return "", err
	}

	return h.pClient.NewMessageVoice(fileData)
}

func (h *Handler) GetClientByProvider(provider config.ProviderModel) providers.Client {
	switch provider {
	case config.IoNet:
		return h.ioClient
	case config.Pollinations:
		return h.pClient
	case config.OpenRouter:
		return h.opClient
	}

	return nil
}
