package entity

import (
	"fmt"
	"regexp"
)

type Person struct {
	FirstName 	string  `json:"firstname" binding:"required"`
	LastName 	string 	`json:"lastname" binding:"required"`
	Age 		int8 	`json:"age" binding:"gte=1,lte=130"`
	Email 		string 	`json:"email" binding:"required,email"`
}

type Video struct {
	Title		string `json:"title" binding:"min=2,max=10" validate:"is-cool"`
	Description	string `json:"description" binding:"max=20"`
	URL			string `json:"url" binding:"required,url"`
	VideoID		string `json:"video_id"`
	Author 		Person `json:"author" binding:"required"`
}

func NewVideo(url string) (*Video, error) {
	videoID, err := extractVideoID(url)
	if err != nil {
		return nil, err
	}

	video := &Video{
		URL:     url,
		VideoID: videoID,
	}
	return video, nil
}

func extractVideoID(url string) (string, error) {
	regexPattern := `v=([^&]+)`

	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindStringSubmatch(url)

	if len(matches) >= 2 {
		videoID := matches[1]
		return videoID, nil
	}

	return "", fmt.Errorf("video ID not found in URL")
}