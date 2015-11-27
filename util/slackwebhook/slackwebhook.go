package slackwebhook

type Payload struct {
	Channel string `json:"channel"`
	Text string `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Fallback string `json:"fallback"`
	Color string `json:"color"`
	Pretext string `json:"pretext"`
	Title string `json:"title"`
	TitleLink string `json:"title_link"`
	Text string `json:"text"`
	Fields []Field `json:"fields"`
	MarkdownIn []string `json:"mrkdwn_in"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool `json:"short"`
}