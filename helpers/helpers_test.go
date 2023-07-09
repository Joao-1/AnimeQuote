package helpers

import (
	"AnimeQuote/common"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadImage(t *testing.T) {
	server := common.Server(t, []common.Route{
		{
			Path: "/image",
			Handler: func(res http.ResponseWriter, req *http.Request, t *testing.T) {
				res.WriteHeader(http.StatusOK)
				_, err := res.Write([]byte("123"))
				if err != nil { t.Fatal(err) }
			},
		},
})

	fakeImageResponse := "123"
	fakeImageResponse64 := "MTIz"

	fakeImgBase64, err := DownloadImage(server.URL + "/image")

	fakeImg, _ := base64.StdEncoding.DecodeString(fakeImgBase64)
	
	assert.Nil(t, err)
	assert.NotEmpty(t, fakeImgBase64)
	assert.Equal(t, fakeImageResponse64, fakeImgBase64)
	assert.Equal(t, fakeImageResponse, string(fakeImg))
}

func TestExtractTwitterImageURL(t *testing.T) {
	t.Run("Extracts image URL from tweet", func(t *testing.T) {
		url := "https://t.co/dinamic"

		extractURL := ExtractTwitterImageURL(fmt.Sprintf("text text %q text text", url))
	
		assert.Equal(t, url, extractURL)
	})

	t.Run("Returns empty string if no image URL is found", func(t *testing.T) {
		extractURL := ExtractTwitterImageURL("text text text text")
	
		assert.Equal(t, "", extractURL)
	})
}