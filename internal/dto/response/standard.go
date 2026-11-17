package response

type StandardResponse struct {
	Status    bool       `json:"status"`
	Code      Code       `json:"code"`
	Message   string     `json:"message"`
	Data      any        `json:"data,omitempty"`
	ErrDetail []FieldErr `json:"error_details,omitempty"`
}

type Code string
type FieldErr struct {
	Field string `json:"field"`
	Msg   string `json:"message"`
}

type Error struct {
	Msg     string
	Code    Code
	Details []FieldErr
}

func (err *Error) Error() string {
	return err.Msg
}
