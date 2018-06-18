package httputil

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
	"github.com/sportivaid/go-template/src/common/apperror"
)

type response struct {
	StatusCode  int         `json:"status_code"`
	Messages    []string    `json:"messages"`
	ProcessTime float64     `json:"process_time"`
	Data        interface{} `json:"data"`
}

func WriteResponse(c *gin.Context, messages []string, processTime float64, data interface{}) {
	c.JSON(
		http.StatusOK,
		response{
			StatusCode:  http.StatusOK,
			Messages:    messages,
			ProcessTime: processTime,
			Data:        data,
		},
	)
}

func WriteErrorResponse(c *gin.Context, processTime float64, err error) {
	errCode := apperror.GetErrorCodes(err)
	c.JSON(
		errCode.HTTPcode,
		response{
			StatusCode:  errCode.StatusCode,
			Messages:    []string{err.Error()},
			ProcessTime: processTime,
			Data:        nil,
		},
	)
}

func WriteDecodeErrorResponse(c *gin.Context, processTime float64, data interface{}) {
	err := apperror.DecodeError
	errCode := apperror.GetErrorCodes(err)
	c.JSON(
		errCode.HTTPcode,
		response{
			StatusCode:  errCode.StatusCode,
			Messages:    []string{err.Error()},
			ProcessTime: processTime,
			Data:        data,
		},
	)
}

func DecodeFormRequest(r *http.Request, req interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		return json.NewDecoder(r.Body).Decode(&req)
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		decoder := form.NewDecoder()
		r.ParseForm()
		return decoder.Decode(&req, r.Form)
	}
	return apperror.DecodeError
}
