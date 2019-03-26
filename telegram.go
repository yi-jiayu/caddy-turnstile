package turnstile

// User represents a Telegram user or bot.
type User struct {
	ID           int64  `json:"id"`
	LanguageCode string `json:"language_code"`
}

// Chat represents a chat.
type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

// Message represents a message.
type Message struct {
	From User `json:"from"`
	Chat Chat `json:"chat"`
}

// InlineQuery represents an incoming inline query.
type InlineQuery struct {
	ID   string `json:"id"`
	From User   `json:"from"`
}

// ChosenInlineResult represents a result of an inline query that was chosen by the user and sent to their chat partner.
type ChosenInlineResult struct {
	ResultID string `json:"result_id"`
	From     User   `json:"from"`
}

// CallbackQuery represents an incoming callback query from a callback button in an inline keyboard.
type CallbackQuery struct {
	ID      string  `json:"id"`
	From    User    `json:"from"`
	Message Message `json:"message"`
}

// Update represents an incoming update.
type Update struct {
	UpdateID           int64               `json:"update_id"`
	Message            *Message            `json:"message"`
	InlineQuery        *InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *CallbackQuery      `json:"callback_query"`
}
