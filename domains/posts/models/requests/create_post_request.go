package requests

type CreatePostRequest struct {
    UserID    string   `form:"user_id" validate:"required,uuid"`
    Caption   string   `form:"caption" validate:"required"`
    Tags      []string `form:"tags"`
    ImageURLs []string
}