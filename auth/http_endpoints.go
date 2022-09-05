package auth

import (
	"github.com/gin-gonic/gin"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"net/http"
)

type HttpEndpoints interface {
	MakeLoginEndpoint() gin.HandlerFunc
	MakeRegisterEndpoint() gin.HandlerFunc
}

type httpEndpoints struct {
	ch setdata_common.CommandHandler
}

func NewHttpEndpoints(ch setdata_common.CommandHandler) HttpEndpoints {
	return &httpEndpoints{ch: ch}
}

func (h *httpEndpoints) MakeLoginEndpoint() gin.HandlerFunc {
	return func(context *gin.Context) {
		setupResponse(context.Writer)
		cmd := &LoginCommand{}
		err := context.ShouldBindJSON(&cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, setdata_common.ErrToHttpResponse(err))
			return
		}
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeRegisterEndpoint() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &RegisterCommand{}
		err := context.ShouldBindJSON(&cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, setdata_common.ErrToHttpResponse(err))
			return
		}
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusCreated, resp)
	}
}

func setupResponse(w http.ResponseWriter) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Allow-Origin")
}
