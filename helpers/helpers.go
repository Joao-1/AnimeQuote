package helpers

import (
	"encoding/base64"
	"io"
	"net/http"
	"regexp"
)

func DownloadImage(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil { return "", err }

	defer res.Body.Close()

	imageData, err := io.ReadAll(res.Body)
	if err != nil { return "", err }
	
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)

	return imageBase64, nil
}

func ExtractTwitterImageURL(text string) string {
	if url := regexp.MustCompile(`https://t.co/([a-zA-Z0-9]+)`).FindStringSubmatch(text); len(url) >= 1 { return url[0] }
	return ""
}