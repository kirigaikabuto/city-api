package applications

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpEndpoints interface {
	MakeCreateApplication() gin.HandlerFunc
	MakeListApplication() gin.HandlerFunc
	MakeUploadApplicationFile() gin.HandlerFunc
	MakeListApplicationByType() gin.HandlerFunc
	MakeGetApplicationById() gin.HandlerFunc

	MakeSearchPlace() gin.HandlerFunc
}

type httpEndpoints struct {
	ch setdata_common.CommandHandler
}

func NewHttpEndpoints(ch setdata_common.CommandHandler) HttpEndpoints {
	return &httpEndpoints{ch: ch}
}

func (h *httpEndpoints) MakeCreateApplication() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &CreateApplicationCommand{}
		jsonData, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		err = json.Unmarshal(jsonData, &cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		respondJSON(context.Writer, http.StatusCreated, resp)
	}
}

func (h *httpEndpoints) MakeListApplication() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &ListApplicationsCommand{}
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		respondJSON(context.Writer, http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeSearchPlace() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &SearchPlaceCommand{}
		name := context.Request.URL.Query().Get("name")
		if name == "" {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrSearchPlaceNoAddressName))
			return
		}
		cmd.Name = name
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		respondJSON(context.Writer, http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeUploadApplicationFile() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &UploadApplicationFileCommand{}
		id := context.Request.URL.Query().Get("id")
		if id == "" {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoApplicationId))
			return
		}
		cmd.Id = id
		buf := bytes.NewBuffer(nil)
		file, _, err := context.Request.FormFile("file")
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		_, err = io.Copy(buf, file)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		err = file.Close()
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		cmd.File = buf
		cmd.ContentType = http.DetectContentType(buf.Bytes())
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		respondJSON(context.Writer, http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeListApplicationByType() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &ListApplicationsByTypeCommand{}
		appType := context.Request.URL.Query().Get("type")
		if appType == "" {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoApplicationType))
			return
		}
		cmd.AppType = appType
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		respondJSON(context.Writer, http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeGetApplicationById() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &GetApplicationByIdCommand{}
		id := context.Request.URL.Query().Get("id")
		if id == "" {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoApplicationId))
			return
		}
		cmd.Id = id
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
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
