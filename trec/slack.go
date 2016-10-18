package tre

type message struct {
	Username    string      `json:"username,omitempty"`
	IconEmoji   string      `json:"icon_emoji,omitempty"`
	Channel     string      `json:"channel,omitempty"`
	Text        string      `json:"text,omitempty"`
	Attachments attachments `json:"attachments,omitempty"`
}

type attachment struct {
	Fallback   string `json:"fallback,omitempty"`
	Color      string `json:"color,omitempty"`
	Pretext    string `json:"pretext,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
	AuthorLink string `json:"author_link,omitempty"`
	AuthorIcon string `json:"author_icon,omitempty"`
	Title      string `json:"title,omitempty"`
	TitleLink  string `json:"title_link,omitempty"`
	Text       string `json:"text,omitempty"`
	Fields     fields `json:"fields,omitempty"`
	ImageURL   string `json:"image_url,omitempty"`
	ThumbURL   string `json:"thumb_url,omitempty"`
	Footer     string `json:"footer,omitempty"`
	FooterIcon string `json:"footer_icon,omitempty"`
	TS         int    `json:"ts,omitempty"`
}

type attachments []attachment

type field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

type fields []field
