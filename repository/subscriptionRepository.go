package repository

import (
	"database/sql"

	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
	"github.com/shopspring/decimal"
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
	var err error
	_, err = db.Exec("INSERT INTO subscriptions (name, description, duration_month, price, is_active, image) VALUES ($1, $2, $3, $4, $5, $6)", name, description, durationMonth, price, true, image)

	if err != nil {
		return nil, err
	}

	var s *model.Subscription
	s, err = GetSubscriptionByName(name, db)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func GetSubscriptionByName(name string, db *sql.DB) (*model.Subscription, error) {
	var err error
	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM subscriptions WHERE name = $1", name)
	if err != nil {
		return nil, err
	}
	var subscriptions []model.Subscription
	for rows.Next() {
		var s model.Subscription
		err = rows.Scan(&s.Id, &s.Name, &s.Description, &s.DurationMonth, &s.Price, &s.IsActive, &s.Image)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}
	if len(subscriptions) == 0 {
		return nil, appErr.NotFoundError("No subscription found.")
	}
	return &subscriptions[0], nil
}

func AnySubscriptionExists(db *sql.DB) (bool, error) {
	rows, rowErr := db.Query("SELECT id FROM subscriptions")
	if rowErr != nil {
		return false, rowErr
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func IsSubscriptionExist(sId int64, db *sql.DB) bool {
	rows, rowErr := db.Query("SELECT id FROM subscriptions WHERE id = $1", sId)
	if rowErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func GetSubscriptionById(sId int64, db *sql.DB) (*model.Subscription, error) {
	var err error
	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM subscriptions WHERE id = $1", sId)
	if err != nil {
		return nil, err
	}
	var subscriptions []model.Subscription
	for rows.Next() {
		var s model.Subscription
		err = rows.Scan(&s.Id, &s.Name, &s.Description, &s.DurationMonth, &s.Price, &s.IsActive, &s.Image)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}
	if len(subscriptions) == 0 {
		return nil, appErr.NotFoundError("No subscription found.")
	}
	return &subscriptions[0], nil
}
