package applications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpEndpoints interface {
	MakeCreateApplication() gin.HandlerFunc
	MakeCreateApplicationWithAuth() gin.HandlerFunc
	MakeListApplication() gin.HandlerFunc
	MakeUploadApplicationFile() gin.HandlerFunc
	MakeListApplicationByType() gin.HandlerFunc
	MakeGetApplicationById() gin.HandlerFunc
	MakeUpdateStatus() gin.HandlerFunc
	MakeAuthorizedUserListApplications() gin.HandlerFunc
	MakeUpdateApplication() gin.HandlerFunc
	MakeRemoveApplication() gin.HandlerFunc
	MakeListByAddressWithAuth() gin.HandlerFunc
	MakeListByAddress() gin.HandlerFunc
	MakeUploadMultipleFiles() gin.HandlerFunc

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
		setupResponse(context.Writer)
		cmd := &CreateApplicationCommand{}
		err := context.ShouldBindJSON(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
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

func (h *httpEndpoints) MakeCreateApplicationWithAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &CreateApplicationCommand{}
		userId, ok := context.Get("user_id")
		if !ok {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoUserIdInToken))
			return
		}
		cmd.UserId = userId.(string)
		err := context.ShouldBindJSON(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
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

func (h *httpEndpoints) MakeListApplication() gin.HandlerFunc {
	return func(c *gin.Context) {
		setupResponse(c.Writer)
		cmd := &ListApplicationsCommand{}
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(c.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.JSON(200, resp)
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

func (h *httpEndpoints) MakeUpdateStatus() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &UpdateApplicationStatusCommand{}
		id := context.Request.URL.Query().Get("id")
		if id == "" {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoApplicationId))
			return
		}
		respJs, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		err = json.Unmarshal(respJs, &cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
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

func (h *httpEndpoints) MakeAuthorizedUserListApplications() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &ListApplicationsByUserIdCommand{}
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

func (h *httpEndpoints) MakeUpdateApplication() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &UpdateApplicationCommand{}
		id := context.Request.URL.Query().Get("id")
		if id == "" {
			context.AbortWithStatusJSON(http.StatusBadRequest, setdata_common.ErrToHttpResponse(ErrNoApplicationId))
			return
		}
		err := context.ShouldBindJSON(&cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		cmd.Id = id
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeRemoveApplication() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &RemoveApplicationCommand{}
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

func (h *httpEndpoints) MakeListByAddressWithAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &ListByAddressCommand{}
		userId, ok := context.Get("user_id")
		if !ok {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoUserIdInToken))
			return
		}
		address := context.Request.URL.Query().Get("address")
		if address == "" {
			context.AbortWithStatusJSON(http.StatusBadRequest, setdata_common.ErrToHttpResponse(ErrNoAddress))
			return
		}
		cmd.UserId = userId.(string)
		cmd.Address = address
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusOK, resp)
	}
}

func (h *httpEndpoints) MakeListByAddress() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &ListByAddressCommand{}
		address := context.Request.URL.Query().Get("address")
		if address == "" {
			context.AbortWithStatusJSON(http.StatusBadRequest, setdata_common.ErrToHttpResponse(ErrNoAddress))
			return
		}
		cmd.Address = address
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.JSON(http.StatusCreated, resp)
	}
}

func (h *httpEndpoints) MakeUploadMultipleFiles() gin.HandlerFunc {
	return func(context *gin.Context) {
		cmd := &UploadMultipleFilesCommand{}
		id := context.Request.URL.Query().Get("id")
		if id == "" {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(ErrNoApplicationId))
			return
		}
		cmd.Id = id

		form, err := context.MultipartForm()
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		formFilesHeader := form.File["files"]
		files := []FileObj{}
		fmt.Println(formFilesHeader)
		for _, fileHeader := range formFilesHeader {
			buf := bytes.NewBuffer(nil)
			file, err := fileHeader.Open()
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
			files = append(files, FileObj{
				File:        buf,
				ContentType: http.DetectContentType(buf.Bytes()),
			})
		}
		cmd.Files = files
		resp, err := h.ch.ExecCommand(cmd)
		if err != nil {
			respondJSON(context.Writer, http.StatusInternalServerError, setdata_common.ErrToHttpResponse(err))
			return
		}
		respondJSON(context.Writer, http.StatusOK, resp)
	}
}

func setupResponse(w http.ResponseWriter) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Allow-Origin")
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
