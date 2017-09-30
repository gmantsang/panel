package shards

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type webhookMessage struct {
	Content   string  `json:"content"`
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Embeds    []embed `json:"embeds"`
}

type embed struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Color       uint      `json:"color"`
	Timestamp   time.Time `json:"timestamp"`
}

func (handler *Handler) sendBrokenWebhook(shard Shard) {
	embeds := make([]embed, 1)
	embeds[0] = embed{
		Title:       "Shard Broken!",
		Description: fmt.Sprintf("**Bot:** %s\n**Shard**: %d\n(Reported By: <@%s>)", shard.Bot, shard.Number, shard.UserID),
		Type:        "rich",
		Color:       0x560404,
		Timestamp:   shard.Timestamp,
	}
	msg := webhookMessage{
		Username:  "Shards",
		AvatarURL: "https://cdn.discordapp.com/emojis/293495010719170560.png",
		Embeds:    embeds,
	}
	handler.sendWebhook(msg)
}

func (handler *Handler) sendFixedWebhook(shard Shard, user string) {
	embeds := make([]embed, 1)
	embeds[0] = embed{
		Title:       "Shard Fixed",
		Description: fmt.Sprintf("**Bot:** %s\n**Shard**: %d\n(Fixed By: <@%s>)", shard.Bot, shard.Number, user),
		Type:        "rich",
		Color:       0x018259,
		Timestamp:   time.Now(),
	}
	msg := webhookMessage{
		Username:  "Shards",
		AvatarURL: "https://cdn.discordapp.com/emojis/293495010719170560.png",
		Embeds:    embeds,
	}
	handler.sendWebhook(msg)
}

func (handler *Handler) sendWebhook(msg webhookMessage) {
	if handler.Config.WebhookURL == "" {
		return
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	buf := bytes.NewBuffer(raw)

	_, err = http.Post(handler.Config.WebhookURL, "application/json", buf)
	if err != nil {
		log.Println(err)
	}
}
