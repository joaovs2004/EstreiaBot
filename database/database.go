package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Client struct {
  gorm.Model
  TelegramId  string
}

type TvShow struct {
  gorm.Model
  ShowId  string
  Name string
  LastSeason int
}

func InitDb() (*gorm.DB) {
  db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  db.AutoMigrate(&Client{}, &TvShow{})

  return db
}
