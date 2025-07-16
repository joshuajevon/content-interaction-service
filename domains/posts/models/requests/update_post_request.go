package requests

type UpdatePostRequest struct {
    Caption   string   `form:"caption" validate:"required"`
    Tags      []string `form:"tags"`
    ImageURLs []string
}