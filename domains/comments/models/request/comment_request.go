package request

type CommentRequest struct {
	Msg     string     `json:"msg" validate:"required"`
}
