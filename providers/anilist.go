package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AnilistResponse struct {
	Data struct {
		Character struct {
		  	Name struct {
				Full string `json:"full"`
		  	}
		  	Image struct { 
				Large string `json:"large"`
				Medium string `json:"medium"`
		  	} `json:"image"`
			}
	} `json:"data"`
}

func GetAnilistCharacterImageURL(characterName, serverUrL string, client http.Client) (AnilistResponse, error) {
	queryForGraphQL := map[string]string{
		"query": fmt.Sprintf(
			`{
				Character(search: %q) {
					name {
						full
					}
					image {
						large
						medium
					  }
				}
			}`,
			characterName,
		),
	}

	jsonQuery, err := json.Marshal(queryForGraphQL)
	if err != nil { return AnilistResponse{}, err }

	req, err := http.NewRequest("POST", serverUrL, bytes.NewBuffer(jsonQuery))
	if err != nil { return AnilistResponse{}, err }

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil { return AnilistResponse{}, err }

	if res.StatusCode == 404 { return AnilistResponse{}, fmt.Errorf("Character does not found") }

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil { return AnilistResponse{}, err }	
	
	var anilistResponse AnilistResponse
	errParse := json.Unmarshal(bodyBytes, &anilistResponse)
	if errParse != nil { return AnilistResponse{}, errParse }	
	
	return anilistResponse, nil

}