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
	fmt.Println("Starting bot...")

	err := godotenv.Load()
	if err != nil { panic(err) }
	
	twitter := new(providers.Twitter)
	twitter = twitter.Init(os.Getenv("CONSUMER_ID"), os.Getenv("CONSUMER_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"), "https://api.twitter.com")

	client := http.Client{}

	animechan := animechan.Animechan{Client: &client}

	fmt.Println("Making first post...")
	MakePost(twitter, client, animechan)
	for range time.Tick(time.Hour * 1) {
		fmt.Println("Making new post...")
		MakePost(twitter, client, animechan)
	}
}

func MakePost(twitter *providers.Twitter, client http.Client, animechan animechan.Animechan) {
	quote, err := animechan.Random().Only()
	fmt.Println(quote)
	if err != nil { fmt.Println(err); MakePost(twitter, client, animechan) }

	regex := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	formatedQuote := fmt.Sprintf(`
	%q - %s, %s

	#anime #quotes #%s #%s 
	`, quote.Quote, quote.Character, quote.Anime, regex.ReplaceAllString(quote.Character, ""), regex.ReplaceAllString(quote.Anime, ""))

	anilistResponse, err := providers.GetAnilistCharacterImageURL(quote.Character, "https://graphql.anilist.co", client)
	if err != nil { fmt.Println(err); MakePost(twitter, client, animechan) }

	media, err := twitter.UploadImage(anilistResponse.Data.Character.Image.Large)
	if err != nil { fmt.Println(err); MakePost(twitter, client, animechan) }

	tweet, _ := twitter.Tweet(providers.TweetParams{Body: formatedQuote, Image: media.Id})

	fmt.Printf("%+v\n", tweet)
}