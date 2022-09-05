package events

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"io"
	"net/http"
)

type HttpEndpoints interface {
	MakeCreateEvent() gin.HandlerFunc
	MakeListEvent() gin.HandlerFunc
	MakeListEventByUserId() gin.HandlerFunc
	MakeUploadDocument() gin.HandlerFunc
}

type httpEndpoints struct {
	ch setdata_common.CommandHandler
}

func NewHttpEndpoints(ch setdata_common.CommandHandler) HttpEndpoints {
	return &httpEndpoints{ch: ch}
}

func (h *httpEndpoints) MakeCreateEvent() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &CreateEventCommand{}
		userId, ok := context.Get("user_id")
		if !ok {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoUserIdInToken))
			return
		}
		cmd.UserId = userId.(string)
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusCreated, resp)
	}
}

func (h *httpEndpoints) MakeListEvent() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &ListEventCommand{}
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		respondJSON(context.Writer, http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeListEventByUserId() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &ListEventByUserIdCommand{}
		userId, ok := context.Get("user_id")
		if !ok {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoUserIdInToken))
			return
		}
		cmd.UserId = userId.(string)
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		respondJSON(context.Writer, http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeUploadDocument() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &UploadDocumentCommand{}
		userId, ok := context.Get("user_id")
		if !ok {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoUserIdInToken))
			return
		}
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

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "applications/json")
	w.WriteHeader(status)
	w.Write(response)
}
