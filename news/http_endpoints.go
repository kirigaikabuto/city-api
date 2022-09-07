package news

import (
	"bytes"
	"github.com/gin-gonic/gin"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"io"
	"net/http"
)

type HttpEndpoints interface {
	MakeCreateNews() gin.HandlerFunc
	MakeListNews() gin.HandlerFunc
	MakeUpdateNews() gin.HandlerFunc
	MakeGetNewsById() gin.HandlerFunc
	MakeGetNewsByAuthorId() gin.HandlerFunc
	MakeUploadPhoto() gin.HandlerFunc
}

type httpEndpoints struct {
	ch setdata_common.CommandHandler
}

func NewHttpEndpoints(ch setdata_common.CommandHandler) HttpEndpoints {
	return &httpEndpoints{ch: ch}
}

func (h *httpEndpoints) MakeCreateNews() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &CreateNewsCommand{}
		userId, ok := context.Get("user_id")
		if !ok {
			context.AbortWithStatusJSON(http.StatusBadRequest, setdata_common.ErrToHttpResponse(ErrNoUserIdInToken))
			return
		}
		cmd.AuthorId = userId.(string)
		err := context.ShouldBindJSON(&cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, setdata_common.ErrToHttpResponse(err))
			return
		}
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusCreated, resp)
	}
}

func (h *httpEndpoints) MakeListNews() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &ListNewsCommand{}
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeUpdateNews() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &UpdateNewsCommand{}
		cmd.Id = context.Request.URL.Query().Get("id")
		err := context.ShouldBindJSON(&cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, setdata_common.ErrToHttpResponse(err))
			return
		}
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeGetNewsById() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &GetNewsByIdCommand{}
		cmd.Id = context.Request.URL.Query().Get("id")
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeGetNewsByAuthorId() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &GetNewsByAuthorId{}
		userId, ok := context.Get("user_id")
		if !ok {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoUserIdInToken))
			return
		}
		cmd.AuthorId = userId.(string)
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeUploadPhoto() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &UploadPhotoCommand{}
		userId, ok := context.Get("user_id")
		if !ok {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoUserIdInToken))
			return
		}
		cmd.Id = context.Request.URL.Query().Get("id")
		cmd.UserId = userId.(string)
		buf := bytes.NewBuffer(nil)
		file, _, err := context.Request.FormFile("file")
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		_, err = io.Copy(buf, file)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		err = file.Close()
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		cmd.File = buf
		cmd.ContentType = http.DetectContentType(buf.Bytes())
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusOK, resp)
	}
}
