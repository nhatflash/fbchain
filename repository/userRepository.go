package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/nhatflash/fbchain/enum"
	"github.com/nhatflash/fbchain/model"
)


func GetSignInUser(email string, password string, db *sql.DB) *model.User {
	
	rows, err := db.Query("SELECT * FROM users WHERE email = $1 AND password = $2 LIMIT 1", email, password)
	if err != nil {
		log.Fatalln("Error when checking user in database", err)
		return nil
	}
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Role, &user.Phone, &user.Identity, &user.FirstName, &user.LastName, &user.Gender, &user.Birthdate, &user.PostalCode, &user.Address, &user.ProfileImage, &user.Status, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Fatalln("Error when mapping sql data into model", err)
			return nil
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return nil
	}
	return &users[0]
}



func RegisterTenant(email string, firstName string, lastName string, password string, gender *enum.Gender, birthdate *time.Time, db *sql.DB) (*model.User, error) {
	
	_, dbErr := db.Exec("INSERT INTO users (email, first_name, last_name, password, gender, birthdate, role, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", email, firstName, lastName, password, gender, birthdate, enum.TENANT, enum.PENDING)
	if dbErr != nil {
		return nil, errors.New("error when creating new tenant: " + dbErr.Error())
	}
	newUser, err := GetUserByEmail(email, db)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return newUser, nil
}


func CheckUserEmailExists(email string, db *sql.DB) bool {
	rows, err := db.Query("SELECT email FROM users WHERE email = $1", email)
	if err != nil {
		log.Fatalln("Error when checking email already exist in user table", err)
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func GetUserByEmail(email string, db *sql.DB) (*model.User, error) {
	rows, err := db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Role, &user.Phone, &user.Identity, &user.FirstName, &user.LastName, &user.Gender, &user.Birthdate, &user.PostalCode, &user.Address, &user.ProfileImage, &user.Status, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return nil, errors.New("no user found")
	}
	return &users[0], nil
}