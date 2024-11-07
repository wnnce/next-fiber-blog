package handler

import (
    "io"
    "log"
	"net/http"
	"net/url"
	"strings"
)

func Handler(response http.ResponseWriter, request *http.Request) {
	originUrl := request.URL.String()
    firstQueryIndex := strings.Index(originUrl, "path")
    var requestUrl string
    if firstQueryIndex > 2 {
        if originUrl[firstQueryIndex-1] == '?' {
            requestUrl = originUrl[6 : firstQueryIndex-1]
        } else {
            requestUrl = originUrl[6 : firstQueryIndex-2]
        }
    } else {
        requestUrl = originUrl[6:]
    }
    requestUrl, _ = url.PathUnescape(requestUrl)
    log.Printf("originUrl: %s, requestUrl: %s", originUrl, requestUrl)
	imageResp, err := http.Get("https://file.qiniu.vnc.ink" + requestUrl)
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	defer imageResp.Body.Close()
	for key, values := range imageResp.Header {
		for _, value := range values {
			response.Header().Add(key, value)
		}
	}
	imageBody, _ := io.ReadAll(imageResp.Body)
	_, _ = response.Write(imageBody)
}