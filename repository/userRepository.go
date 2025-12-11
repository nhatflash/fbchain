package repository

import (
	"database/sql"
	"log"
	"time"
	appError "github.com/nhatflash/fbchain/error"
	_ "github.com/lib/pq"
	"github.com/nhatflash/fbchain/enum"
	"github.com/nhatflash/fbchain/model"
)


func InitialRegisterTenant(email string, firstName string, lastName string, password string, gender *enum.Gender, birthdate *time.Time, db *sql.DB) (*model.User, error) {
	
	_, dbErr := db.Exec("INSERT INTO users (email, first_name, last_name, password, gender, birthdate, role, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", email, firstName, lastName, password, gender, birthdate, enum.TENANT, enum.PENDING)
	if dbErr != nil {
		return nil, appError.ErrInternal
	}
	newUser, userErr := GetUserByEmail(email, db)
	if userErr != nil {
		return nil, userErr
	}
	return newUser, nil
}


func CompletedRegisterTenant(email string, phone string, identity string, address string, postalCode string, profileImage string, db *sql.DB) (*model.User, error) {
	_, dbErr := db.Exec("UPDATE users SET phone = $1, identity = $2, address = $3, postal_code = $4, profile_image = $5, status = $6 WHERE email = $7", phone, identity, address, postalCode, profileImage, enum.ACTIVE, email)
	if dbErr != nil {
		return nil, appError.ErrInternal
	}

	updatedUser, userErr := GetUserByEmail(email, db)
	if userErr != nil {
		return nil, userErr
	}
	return updatedUser, nil
}


func CheckUserEmailExists(email string, db *sql.DB) bool {
	rows, err := db.Query("SELECT email FROM users WHERE email = $1 LIMIT 1", email)
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
	rows, rowErr := db.Query("SELECT * FROM users WHERE email = $1 LIMIT 1", email)
	if rowErr != nil {
		return nil, rowErr
	}
	var users []model.User
	for rows.Next() {
		var user model.User
		scanErr := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Role, &user.Phone, &user.Identity, &user.FirstName, &user.LastName, &user.Gender, &user.Birthdate, &user.PostalCode, &user.Address, &user.ProfileImage, &user.Status, &user.CreatedAt, &user.UpdatedAt)
		if scanErr != nil {
			return nil, scanErr
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return nil, appError.ErrNotFound
	}
	return &users[0], nil
}


func GetUserByPhone(phone string, db *sql.DB) (*model.User, error) {
	rows, rowErr := db.Query("SELECT * FROM users WHERE phone = $1 LIMIT 1", phone)
	if rowErr != nil {
		return nil, rowErr
	}
	var users []model.User
	for rows.Next() {
		var user model.User
		scanErr := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Role, &user.Phone, &user.Identity, &user.FirstName, &user.LastName, &user.Gender, &user.Birthdate, &user.PostalCode, &user.Address, &user.ProfileImage, &user.Status, &user.CreatedAt, &user.UpdatedAt)
		if scanErr != nil {
			return nil, scanErr
		}
		users = append(users, user)
	}
	if (len(users) == 0) {
		return nil, appError.ErrNotFound
	} 
	return &users[0], nil
}