package utility

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
)

var S3CdnUrl string

func init() {
	S3CdnUrl = os.Getenv("S3_CDN_URL")
}

type PagingDTO struct {
	Limit  int     `json:"limit"`
	LastID *string `json:"last_id"`
}

type PagingResponseDTO struct {
	LastId string `json:"last_id"`
}

func GetCurrentTimeInMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetTimeInMs(
	numberOfHours time.Duration,
) int64 {
	t := time.Now().Add(numberOfHours * time.Hour)
	return t.UnixNano() / int64(time.Millisecond)
}

func GetCdnUrl(
	path string,
	source string,
) string {
	if source == "firebase" {
		return path
	}
	return S3CdnUrl+ "/" + path
}

func GetPagingDTO(
	queryParams map[string][]string,
) PagingDTO {
	response := PagingDTO{
		Limit: 2,
	}
	if queryParams["limit"] != nil {
		limit := queryParams["limit"][0]
		if limit != "" {
			val, err := strconv.Atoi(limit)
			if err != nil {
				response.Limit = val
			}
		}
	}

	if queryParams["last_id"] != nil {
		lastId := queryParams["last_id"][0]
		if lastId != "" {
			response.LastID = &lastId
		}
	}
	return response
}

func SendHttpResponse(
	responseWriter http.ResponseWriter,
	statusCode int,
	response interface{},
) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	json.NewEncoder(responseWriter).Encode(response)
	return
}


type emptyHttpResponse struct {}
func SendEmptyHttpResponse(
	responseWriter http.ResponseWriter,
	statusCode int,
) {
	responseData := emptyHttpResponse{}
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	json.NewEncoder(responseWriter).Encode(responseData)
	return
}