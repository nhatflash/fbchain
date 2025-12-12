package repository

import (
	"database/sql"
	"time"
	appError "github.com/nhatflash/fbchain/error"
	_ "github.com/lib/pq"
	"github.com/nhatflash/fbchain/enum"
	"github.com/nhatflash/fbchain/model"
)


func CheckUserEmailExists(email string, db *sql.DB) bool {
	rows, err := db.Query("SELECT email FROM users WHERE email = $1 LIMIT 1", email)
	if err != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}


func CheckUserPhoneExists(phone string, db *sql.DB) bool {
	rows, dbErr := db.Query("SELECT phone FROM users WHERE phone = $1 LIMIT 1", phone)
	if dbErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func CheckUserIdentityExists(identity string, db *sql.DB) bool {
	rows, dbErr := db.Query("SELECT identity FROM users WHERE identity = $1 LIMIT 1", identity)
	if dbErr != nil {
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


func CreateTenantUser(firstName string, lastName string, email string, password string, birthdate *time.Time, gender *enum.Gender, phone string, identity string, address string, postalCode string, profileImage string, db *sql.DB) (*model.User, error) {
	_, insertErr := db.Exec("INSERT INTO users (email, password, role, phone, identity, first_name, last_name, gender, birthdate, postal_code, address, profile_image, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", email, password, enum.TENANT, phone, identity, firstName, lastName, gender, birthdate, postalCode, address, profileImage, enum.ACTIVE);
	if insertErr != nil {
		return nil, insertErr
	}
	tenantUser, err := GetUserByEmail(email, db)
	if err != nil {
		return nil, err
	}
	return tenantUser, nil
}


func CreateAdminUser(email string, password string, phone string, identity string, firstName string, lastName string, gender *enum.Gender, birthdate *time.Time, postalCode string, address string, profileImage string, db *sql.DB) error {
	_, insertErr := db.Exec("INSERT INTO users (email, password, role, phone, identity, first_name, last_name, gender, birthdate, postal_code, address, profile_image, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", email, password, enum.ADMIN, phone, identity, firstName, lastName, gender, birthdate, postalCode, address, profileImage, enum.ACTIVE)

	if insertErr != nil {
		return insertErr
	}
	return nil
}


func CheckIfAdminUserAlreadyExists(db *sql.DB) (bool, error) {
	rows, rowErr := db.Query("SELECT id FROM users WHERE role = $1", enum.ADMIN)
	if rowErr != nil {
		return false, rowErr
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}