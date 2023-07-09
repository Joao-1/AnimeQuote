package providers

import (
	"AnimeQuote/common"
	"fmt"
	"net/http"
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

var (
	tweetBody = "Hello World!"
	tweetId = "1676661887809855481"
	tweetResponse = fmt.Sprintf(`{"data": {"edit_history_tweet_ids": ["1676661887809855488"],"id": %q,"text": %q}}`, tweetId, tweetBody)
)

func TestTwitter(t *testing.T) {
	uploadMediaId := "123"

	server := common.Server(t, []common.Route{
		{
			Path: "/2/tweets",
			Handler: func(res http.ResponseWriter, req *http.Request, t *testing.T) {
				res.WriteHeader(http.StatusCreated)
				_, err := res.Write([]byte(tweetResponse))
				if err != nil { t.Fatal(err) }
			},
		},
		{
			Path: "/1.1/media/upload.json",
			Handler: func(res http.ResponseWriter, req *http.Request, t *testing.T) {
				res.WriteHeader(http.StatusOK)
				_, err := res.Write([]byte(fmt.Sprintf(`{"media_id_string": %q}`, uploadMediaId)))
				if err != nil { t.Fatal(err) }
			},
		},
		{
			Path: "/image",
			Handler: func(res http.ResponseWriter, req *http.Request, t *testing.T) {
				res.WriteHeader(http.StatusOK)
				_, err := res.Write([]byte("123"))
				if err != nil { t.Fatal(err) }
			},
		},
	})

	fakeImageToUpload := server.URL + "/image"

	twitter := new(Twitter) 
	twitter = twitter.Init("123", "456", "789", "012", server.URL)
	
	t.Run("Tweet", func(t *testing.T) {
		tweet, err := twitter.Tweet(TweetParams{Body: tweetBody})

		assert.Nil(t, err)
		assert.NotEmpty(t, tweet.Data.Id)
		assert.Equal(t, tweet.Data.Id, tweetId)
	})

	t.Run("TweetWithImage", func(t *testing.T) {
		tweet, err := twitter.Tweet(TweetParams{Body: tweetBody, Image: fakeImageToUpload})

		assert.Nil(t, err)
		assert.NotEmpty(t, tweet.Data.Id)
		assert.Equal(t, tweet.Data.Id, tweetId)
	})

	t.Run("UploadImage", func(t *testing.T) {
		uploadedMediaDetails, err := twitter.UploadImage(fakeImageToUpload)

		assert.Nil(t, err)
		assert.NotEmpty(t, uploadedMediaDetails.Id)
		assert.Equal(t, uploadedMediaDetails.Id, uploadMediaId)
	})
}

func TestGetCharacterImage(t *testing.T) {
	character := "Echidna";
	characterImage := "https://static.wikia.nocookie.net/rezero/images/0/02/Echidna_Anime_PV_2.png/revision/latest?cb=20200611134057"
	getResponse := fmt.Sprintf(`{"data": {"character": {"name": {"full": %q},"image": {"large": %q,"medium": ""}}}}`, character, characterImage)

	server := common.Server(t, []common.Route{
		{
			Path: "/",
			Handler: func(res http.ResponseWriter, req *http.Request, t *testing.T) {
				res.WriteHeader(http.StatusOK)
				_, err := res.Write([]byte(getResponse))
				if err != nil { t.Fatal(err) }
			},
		},
	})

	client := http.Client{}

	res, err := GetAnilistCharacterImageURL(character, server.URL, client)
	
	assert.Nil(t, err)
	assert.Equal(t, character, res.Data.Character.Name.Full)
	assert.Equal(t, characterImage, res.Data.Character.Image.Large)
}