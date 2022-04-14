package auth

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kirigaikabuto/city-api/api_keys"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"net/http"
)

type ApiKeyMdw interface {
	MakeApiKeyMiddleware() gin.HandlerFunc
}

type apiKeyMdw struct {
	apiKeyStore api_keys.ApiKeyStore
}

func NewApiKeyMdw(apiKeyStore api_keys.ApiKeyStore) ApiKeyMdw {
	return &apiKeyMdw{apiKeyStore: apiKeyStore}
}

func (a *apiKeyMdw) MakeApiKeyMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		apiKeyVal := context.Request.Header.Get("Api-Key")
		if apiKeyVal == "" {
			respondJSON(context.Writer, http.StatusBadRequest, setdata_common.ErrToHttpResponse(ErrNoApiKeyHeaderValue))
			return
		}
		_, err := a.apiKeyStore.GetByKey(apiKeyVal)
		if err != nil && err == api_keys.ErrApiKeyNotFound {
			respondJSON(context.Writer, http.StatusBadRequest, setdata_common.ErrToHttpResponse(ErrIncorrectApiKey))
			return
		} else if err != nil {
			respondJSON(context.Writer, http.StatusBadRequest, setdata_common.ErrToHttpResponse(err))
			return
		}
		context.Next()
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
