package message

type Message struct {
	// The message content
	Type    string `json:"type"`
	URL     string `json:"url"`
	Name    string `json:"name"`
	Content string `json:"content"`
}
