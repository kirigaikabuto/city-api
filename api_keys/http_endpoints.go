package api_keys

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"net/http"
)

type HttpEndpoints interface {
	MakeCreateApiKey() gin.HandlerFunc
	MakeListApiKey() gin.HandlerFunc
}

type httpEndpoints struct {
	ch setdata_common.CommandHandler
}

func NewHttpEndpoints(ch setdata_common.CommandHandler) HttpEndpoints {
	return &httpEndpoints{ch: ch}
}

func (h *httpEndpoints) MakeCreateApiKey() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &CreateApiKeyCommand{}
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, err)
			return
		}
		respondJSON(context.Writer, http.StatusCreated, resp)
	}
}

func (h *httpEndpoints) MakeListApiKey() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &ListApiKeyCommand{}
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, err)
			return
		}
		respondJSON(context.Writer, http.StatusOK, resp)
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
