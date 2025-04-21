package database

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	TelegramID int64
}

type TvShow struct {
	gorm.Model
	ShowID     string
	Name       string
	LastSeason int
}

type ClientSubscription struct {
	gorm.Model
	ClientID int64
	ShowID   string
}
