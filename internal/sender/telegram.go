package sender

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TelegramSender struct {
	botToken string
	client   *http.Client
}

func NewTelegram(botToken string) *TelegramSender {
	return &TelegramSender{
		botToken: botToken,
		client:   &http.Client{},
	}
}

func (t *TelegramSender) Send(ctx context.Context, recipient, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	resp, err := t.client.PostForm(apiURL, url.Values{
		"chat_id": {recipient},
		"text":    {message},
	})
	if err != nil {
		return fmt.Errorf("telegram request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram bad status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
