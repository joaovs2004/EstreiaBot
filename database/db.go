package database

import (
	"fmt"

	"estreiaBot/api"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func InitDb() {
	DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&Client{}, &TvShow{}, &ClientSubscription{})
}

// checks if a user exists in the database, and creates them if they don't
func CreateUser(telegramID int64) {
	var user Client

	// Check if the user exists
	result := DB.First(&user, "telegram_id = ?", telegramID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println(telegramID)
			user = Client{TelegramID: telegramID}
			DB.Create(&user)
		}
	}
}

// checks if the show exists in the database, and creates them if they don't
func CreateShow(showId string, showName string) {
	var show TvShow

	// Check if the show exists
	result := DB.First(&show, "show_id = ?", showId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			showLastSeason := api.GetLastSeason(showId)
			show = TvShow{ShowID: showId, Name: showName, LastSeason: showLastSeason}
			DB.Create(&show)
		}
	}
}

// checks if the client subscription exists in the database, and creates them if they don't
func CreateClientSubscription(clientID int64, showID string) {
	var clientSubscription ClientSubscription

	// Check if the subscription exists
	result := DB.First(&clientSubscription, "client_id = ? AND show_id = ?", clientID, showID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			clientSubscription = ClientSubscription{ClientID: clientID, ShowID: showID}
			DB.Create(&clientSubscription)
		}
	}
}

func GetClientSubscriptions(clientID int64) []TvShow {
	var subscriptions []TvShow
	var clientSubscriptions []ClientSubscription

	// Get all subscriptions for the client
	DB.Where("client_id = ?", clientID).Find(&clientSubscriptions)

	// Get all shows for the subscriptions
	for _, subscription := range clientSubscriptions {
		var show TvShow
		DB.First(&show, "show_id = ?", subscription.ShowID)
		subscriptions = append(subscriptions, show)
	}

	return subscriptions
}

func RemoveClientSubscription(clientID int64, showID string) {
	var clientSubscription ClientSubscription

	// Check if the subscription exists
	result := DB.First(&clientSubscription, "client_id = ? AND show_id = ?", clientID, showID)
	if result.Error == nil {
		DB.Delete(&clientSubscription)
	}
}
