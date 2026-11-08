package response

type Sender struct {
	SenderID  int64
	UserName  string
	ImagePath string
}

type Attachment struct {
	AttachmentID int64
	FileType     string
	FileSize     string
	FilePath     string
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
