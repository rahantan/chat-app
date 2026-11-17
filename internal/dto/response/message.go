package response

type Sender struct {
	UserID    int64  `json:"user_id"`
	UserName  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

type Attachment struct {
	AttachmentID int64  `json:"attachment_id"`
	FileType     string `json:"file_type"`
	FileSize     string `json:"file_size"`
	FileUrl      string `json:"file_url"`
}

type Message struct {
	MsgID       int64        `json:"message_id"`
	MsgType     string       `json:"message_type"`
	Content     string       `json:"content"`
	Sender      *Sender      `json:"sender"`
	Attachments []Attachment `json:"attachments"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAtAt string       `json:"updated_at"`
}
