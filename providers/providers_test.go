package providers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}

func TestTwitter(t *testing.T) {
	tweetBody := "Hello World!"
	tweetId := "1676661887809855481"
	tweetResponse :=  fmt.Sprintf(`{"data": {"edit_history_tweet_ids": ["1676661887809855488"],"id": %q,"text": %q}}`, tweetId, tweetBody)
	fakeImageResponse := "123"

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		switch req.URL.Path {
			case "/2/tweets":
				res.WriteHeader(http.StatusCreated)
				_, err := res.Write([]byte(tweetResponse))
				if err != nil { t.Fatal(err) }
			case "/image":
				res.WriteHeader(http.StatusOK)
				_, err := res.Write([]byte(fakeImageResponse))
				if err != nil { t.Fatal(err) }
			case "/1.1/media/upload.json":
				fmt.Println("req: ", req)
				res.WriteHeader(http.StatusOK)
				_, err := res.Write([]byte(`{"media_id_string": "123"}`))
				if err != nil { t.Fatal(err) }
		}
	}))

	fakeImageToDownload := server.URL + "/image"

	twitter := new(Twitter) 
	twitter = twitter.Init("123", "456", "789", "012", server.URL)
	
	t.Run("Tweet", func(t *testing.T) {
		tweet, err := twitter.Tweet(TweetParams{Body: tweetBody})

		assert.Nil(t, err)
		assert.NotEmpty(t, tweet.Data.Id)
		assert.Equal(t, tweet.Data.Id, tweetId)
	})

	t.Run("UploadImage", func(t *testing.T) {
		mediaId := "123"

		uploadedMediaDetails, err := twitter.UploadImage(fakeImageToDownload)

		assert.Nil(t, err)
		assert.NotEmpty(t, uploadedMediaDetails.id)
		assert.Equal(t, uploadedMediaDetails.id, mediaId)
	})
	
	t.Run("Download image", func(t *testing.T) {
		oneTwoThreebase64 := "MTIz"

		fakeImgBase64, err := downloadImage(fakeImageToDownload)

		fakeImg, _ := base64.StdEncoding.DecodeString(fakeImgBase64)

		assert.Nil(t, err)
		assert.NotEmpty(t, fakeImgBase64)
		assert.Equal(t, fakeImgBase64, oneTwoThreebase64)
		assert.Equal(t, fakeImageResponse, string(fakeImg))
	})
}