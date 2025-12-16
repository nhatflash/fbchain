package repository

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/nhatflash/fbchain/enum"
	appError "github.com/nhatflash/fbchain/error"
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
	var err error
	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM users WHERE email = $1 LIMIT 1", email)
	if err != nil {
		return nil, err
	}
	var users []model.User
	for rows.Next() {
		var u model.User
		err = rows.Scan(&u.Id, &u.Email, &u.Password, &u.Role, &u.Phone, &u.Identity, &u.FirstName, &u.LastName, &u.Gender, &u.Birthdate, &u.PostalCode, &u.Address, &u.ProfileImage, &u.Status, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, appError.ErrNotFound
	}
	return &users[0], nil
}

func GetUserByPhone(phone string, db *sql.DB) (*model.User, error) {
	var err error
	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM users WHERE phone = $1 LIMIT 1", phone)
	if err != nil {
		return nil, err
	}
	var users []model.User
	for rows.Next() {
		var u model.User
		err = rows.Scan(&u.Id, &u.Email, &u.Password, &u.Role, &u.Phone, &u.Identity, &u.FirstName, &u.LastName, &u.Gender, &u.Birthdate, &u.PostalCode, &u.Address, &u.ProfileImage, &u.Status, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, appError.ErrNotFound
	}
	return &users[0], nil
}

func CreateTenantUser(firstName string, lastName string, email string, password string, birthdate *time.Time, gender *enum.Gender, phone string, identity string, address string, postalCode string, profileImage string, db *sql.DB) (*model.User, error) {
	var err error

	_, err = db.Exec("INSERT INTO users (email, password, role, phone, identity, first_name, last_name, gender, birthdate, postal_code, address, profile_image, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", email, password, enum.ROLE_TENANT, phone, identity, firstName, lastName, gender, birthdate, postalCode, address, profileImage, enum.USER_ACTIVE)
	if err != nil {
		return nil, err
	}
	tUser, err := GetUserByEmail(email, db)
	if err != nil {
		return nil, err
	}
	return tUser, nil
}

func CreateAdminUser(email string, password string, phone string, identity string, firstName string, lastName string, gender *enum.Gender, birthdate *time.Time, postalCode string, address string, profileImage string, db *sql.DB) error {
	var err error
	_, err = db.Exec("INSERT INTO users (email, password, role, phone, identity, first_name, last_name, gender, birthdate, postal_code, address, profile_image, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", email, password, enum.ROLE_ADMIN, phone, identity, firstName, lastName, gender, birthdate, postalCode, address, profileImage, enum.USER_ACTIVE)

	if err != nil {
		return err
	}
	return nil
}

func CheckIfAdminUserAlreadyExists(db *sql.DB) (bool, error) {
	rows, rowErr := db.Query("SELECT id FROM users WHERE role = $1", enum.ROLE_ADMIN)
	if rowErr != nil {
		return false, rowErr
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}


func ListAllUsers(db *sql.DB) ([]model.User, error) {
	var rows *sql.Rows
	var err error
	rows, err = db.Query("SELECT * FROM users ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	var users []model.User
	for rows.Next() {
		var u model.User
		err = rows.Scan(&u.Id, &u.Email, &u.Password, &u.Role, &u.Phone, &u.Identity, &u.FirstName, &u.LastName, &u.Gender, &u.Birthdate, &u.PostalCode, &u.Address, &u.ProfileImage, &u.Status, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, appError.NotFoundError("No user found.")
	}
	return users, nil
}


func GetUserById(db *sql.DB, id int64) (*model.User, error) {
	var rows *sql.Rows
	var err error
	rows, err = db.Query("SELECT * FROM users WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	var users []model.User
	for rows.Next() {
		var u model.User
		err = rows.Scan(&u.Id, &u.Email, &u.Password, &u.Role, &u.Phone, &u.Identity, &u.FirstName, &u.LastName, &u.Gender, &u.Birthdate, &u.PostalCode, &u.Address, &u.ProfileImage, &u.Status, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, appError.NotFoundError("No user found.")
	}
	return &users[0], nil
}
