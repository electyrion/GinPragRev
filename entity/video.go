package entity

import (
	"fmt"
	"regexp"
	"time"
)

type Person struct {
	ID 			uint64	`gorm:"primary_key;auto_icrement" json:"id"`
	FirstName 	string  `json:"firstname" binding:"required" gorm:"type:varchar(32)"`
	LastName 	string 	`json:"lastname" binding:"required" gorm:"type:varchar(32)"`
	Age 		int8 	`json:"age" binding:"gte=1,lte=130"`
	Email 		string 	`json:"email" binding:"required,email" gorm:"type:varchar(256)"`
}

type Video struct {
	ID 			uint64 `gorm:"primary_key;auto_icrement" json:"id"`
	Title		string `json:"title" binding:"min=2,max=100" validate:"is-cool" gorm:"type:varchar(100)"`
	Description	string `json:"description" binding:"max=20" gorm:"type:varchar(200)"`
	URL			string `json:"url" binding:"required,url" gorm:"type:varchar(250);UNIQUE"`
	URLid		string `json:"url_id"`
	Author 		Person `json:"author" binding:"required" gorm:"foreignkey:PersonID"`
	PersonID 	uint64 `json:-`
	CreatedAt	time.Time `json:"-" gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt	time.Time `json:"-" gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func NewVideo(url string) (*Video, error) {
	urlID, err := extractVideoID(url)
	if err != nil {
		return nil, err
	}

	video := &Video{
		URL:     url,
		URLid: urlID,
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