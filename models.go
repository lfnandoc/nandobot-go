package main

import (
	"time"

	"gorm.io/gorm"

	"gorm.io/driver/sqlite"

	"github.com/lib/pq"
)

var DB *gorm.DB

type PlayerMatchInfo struct {
	gorm.Model
	PlayerID          uint
	PlayerName        string
	MatchID           string
	GameMode          string
	Win               bool
	Champion          string
	Skin              string
	Kills             int
	Deaths            int
	Assists           int
	EndTime           time.Time
	Duration          int
	TotalDamageDealt  int
	TotalHeal         int
	TeamPosition      string
	MythicItem        int
	Items             pq.Int64Array `gorm:"type:integer[6]"`
	ItemListImageLink string
	MessageSent       bool
}

type Player struct {
	gorm.Model
	Name    string
	UserId  string
	Matches []PlayerMatchInfo
}

func SetupDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("nandobot.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&Player{}, &PlayerMatchInfo{})
}
