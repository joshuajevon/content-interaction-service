package http

import (
	"bootcamp-content-interaction-service/domains/posts"
	"bootcamp-content-interaction-service/domains/posts/models/requests"
	"bootcamp-content-interaction-service/shared/models/responses"
	"strconv"

	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PostHttp struct {
	postUc posts.PostUseCase
}

func NewPostHttp(postUc posts.PostUseCase) *PostHttp{
	return &PostHttp{
		postUc: postUc,
	}
}

func (handler *PostHttp) UpdatePost(c *gin.Context) {
    ctx := c.Request.Context()
    postID := c.Param("id")

    var form requests.UpdatePostRequest
    if err := c.ShouldBind(&form); err != nil {
        c.JSON(http.StatusBadRequest, responses.BasicResponse{Error: err.Error()})
        return
    }

    formData, err := c.MultipartForm()
    if err == nil {
        var imageURLs []string
        for _, file := range formData.File["images"] {
            timestamp := time.Now().Format("20060102150405")
            savePath := "storage/post/" + timestamp + "_" + file.Filename
            if err := c.SaveUploadedFile(file, savePath); err != nil {
                c.JSON(http.StatusInternalServerError, responses.BasicResponse{Error: "File upload error: " + err.Error()})
                return
            }
            imageURLs = append(imageURLs, savePath)
        }
        form.ImageURLs = imageURLs
    }

    result, err := handler.postUc.UpdatePost(ctx, postID, &form)
    if err != nil {
        c.JSON(http.StatusInternalServerError, responses.BasicResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}

func (handler *PostHttp) DeletePost(c *gin.Context) {
    ctx:= c.Request.Context()

    postId := c.Param("id")
    
    result, err := handler.postUc.DeletePost(ctx, postId)

     if err != nil {
        c.JSON(http.StatusInternalServerError, responses.BasicResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}

func (handler *PostHttp) ViewPostById(c *gin.Context) {
    ctx := c.Request.Context()
    postID := c.Param("id")
    result, err := handler.postUc.ViewPostById(ctx, postID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, responses.BasicResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}

func (handler *PostHttp) ViewAllPostByUserId(c *gin.Context) {
    ctx := c.Request.Context()

    result, err := handler.postUc.ViewAllPostByUserId(ctx)

    if err != nil {
        c.JSON(http.StatusInternalServerError, responses.BasicResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}

func (handler *PostHttp) ViewAllPost(c *gin.Context) {
    ctx := c.Request.Context()

    result, err := handler.postUc.ViewAllPost(ctx)

    if err != nil {
        c.JSON(http.StatusInternalServerError, responses.BasicResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}

func (handler *PostHttp) CreatePost(c *gin.Context) {
    ctx := c.Request.Context()
    var form requests.CreatePostRequest

    if err := c.ShouldBind(&form); err != nil {
        c.JSON(http.StatusBadRequest, responses.BasicResponse{Error: err.Error()})
        return
    }

    formData, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, responses.BasicResponse{Error: "Invalid multipart form: " + err.Error()})
        return
    }

    files := formData.File["images"]
    if len(files) == 0 {
        c.JSON(http.StatusBadRequest, responses.BasicResponse{Error: "At least one image is required"})
        return
    }

    var imagePaths []string
    for _, file := range files {
        timestamp := time.Now().Format("20060102150405")
        savePath := "storage/post/" + timestamp + "_" + file.Filename
        if err := c.SaveUploadedFile(file, savePath); err != nil {
            for _, path := range imagePaths {
                _ = os.Remove(path)
            }
            c.JSON(http.StatusInternalServerError, responses.BasicResponse{Error: "Failed to save image: " + err.Error()})
            return
        }
        imagePaths = append(imagePaths, savePath)
    }

    form.ImageURLs = imagePaths

    validate := validator.New()
    if err := validate.StructCtx(ctx, form); err != nil {
        for _, path := range imagePaths {
            _ = os.Remove(path)
        }
        c.JSON(http.StatusBadRequest, responses.BasicResponse{Error: err.Error()})
        return
    }

    result, err := handler.postUc.CreatePost(ctx, &form)
    if err != nil {
        for _, path := range imagePaths {
            _ = os.Remove(path)
        }
        c.JSON(http.StatusInternalServerError, responses.BasicResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusCreated, result)
}

func (handler *PostHttp) ViewPersonalFeed(c *gin.Context) {
	ctx := c.Request.Context()
	userId := c.Param("id")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, responses.BasicResponse{Error: "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, responses.BasicResponse{Error: "Invalid limit parameter"})
		return
	}

	offset := (page - 1) * limit

	result, err := handler.postUc.ViewPostByUserId(ctx, userId, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.BasicResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.BasicResponse{Data: result})
}
