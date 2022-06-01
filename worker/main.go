package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string
	Processed bool
}

type Response struct {
	Email          string `json:"email"`
	TimeFromInsert int    `json:"time_from_insert"`
}

func selectNew(db *gorm.DB) {
	var result []User

	db.Model(User{}).
		Where("processed = ?", false).
		Select("email", "created_at").
		Find(&result)

	db.Model(User{}).
		Where("processed = ?", false).
		Update("processed", "true")

	res := make([]Response, len(result))
	now := time.Now()

	for i, x := range result {
		res[i] = Response{Email: x.Email, TimeFromInsert: int(now.Sub(x.CreatedAt).Seconds())}
	}
	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
}

func main() {
	db, e := gorm.Open(postgres.Open("host=127.0.0.1 port=5432 user=accounts password=accounts dbname=accounts sslmode=disable"), &gorm.Config{})
	if e != nil {
		panic("failed to connect database")
	}

	fmt.Println("started")
	s := gocron.NewScheduler(time.UTC)
	s.Cron("* * * * *").Do(selectNew, db)
	s.StartBlocking()

}
