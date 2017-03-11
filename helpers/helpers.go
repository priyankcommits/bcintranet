package helpers

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

func ImageToBase64(url string) string {
	// Convert url image to base64 encoding
	res, _ := http.Get(url)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	imgBase64Str := base64.StdEncoding.EncodeToString(bodyBytes)
	return imgBase64Str
}
