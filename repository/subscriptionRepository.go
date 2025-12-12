package repository

import (
	"database/sql"
	"github.com/nhatflash/fbchain/model"
	"github.com/shopspring/decimal"
	appErr "github.com/nhatflash/fbchain/error"
)

func CheckSubscriptionNameExists(name string, db *sql.DB) bool {
	rows, rowErr := db.Query("SELECT name FROM subscriptions WHERE name = $1", name)
	if rowErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}


func CreateSubscription(name string, description string, durationMonth int, price decimal.Decimal, image string, db *sql.DB) (*model.Subscription, error) {
	_, insertErr := db.Exec("INSERT INTO subscriptions (name, description, duration_month, price, is_active, image) VALUES ($1, $2, $3, $4, $5, $6)", name, description, durationMonth, price, true, image)
	
	if insertErr != nil {
		return nil, insertErr
	}
	s, err := GetSubscriptionByName(name, db)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func GetSubscriptionByName(name string, db *sql.DB) (*model.Subscription, error) {
	rows, selectErr := db.Query("SELECT * FROM subscriptions WHERE name = $1", name)
	if selectErr != nil {
		return nil, selectErr
	}
	var subscriptions []model.Subscription
	for rows.Next() {
		var s model.Subscription
		scanErr := rows.Scan(&s.Id, &s.Name, &s.Description, &s.DurationMonth, &s.Price, &s.IsActive, &s.Image)
		if scanErr != nil {
			return nil, scanErr
		}
		subscriptions = append(subscriptions, s)
	}
	if len(subscriptions) == 0 {
		return nil, appErr.NotFoundError("No subscription found.")
	}
	return &subscriptions[0], nil
}