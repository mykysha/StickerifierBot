package domain

// Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

// Message represents a Telegram message.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// Update represents an incoming Telegram update.
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}
