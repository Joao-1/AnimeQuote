package main

import (
	"AnimeQuote/providers"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/joao-1/animechan-go"
	"github.com/joho/godotenv"
)



func main() {
	godotenv.Load()
	
	twitter := new(providers.Twitter)
	twitter = twitter.Init(os.Getenv("CONSUMER_ID"), os.Getenv("CONSUMER_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"), "https://api.twitter.com")

	client := http.Client{}

	animechan := animechan.Animechan{Client: &client}

	for range time.Tick(time.Hour * 1) {
	quote, err := animechan.Random().Only()
	if err != nil { panic(err) }

	regex := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	formatedQuote := fmt.Sprintf(`
	%q - %s, %s

	#anime #quotes #%s #%s 
	`, quote.Quote, quote.Character, quote.Anime, regex.ReplaceAllString(quote.Character, ""), regex.ReplaceAllString(quote.Anime, ""))

	anilistResponse, err := providers.GetAnilistCharacterImageURL(quote.Character, "https://graphql.anilist.co", client)
	if err != nil { panic(err) }

	media, err := twitter.UploadImage(anilistResponse.Data.Character.Image.Large)
	if err != nil { panic(err) }

	tweet, _ := twitter.Tweet(providers.TweetParams{Body: formatedQuote, Image: media.Id})

	fmt.Printf("%+v\n", tweet)
	}
}