package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/nhatflash/fbchain/enum"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
)

type UserRepository struct {
	Db 		*sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

func (ur *UserRepository) CheckUserEmailExists(email string) bool {
	rows, err := ur.Db.Query("SELECT email FROM users WHERE email = $1 LIMIT 1", email)
	if err != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func (ur *UserRepository) CheckUserPhoneExists(phone string) bool {
	rows, dbErr := ur.Db.Query("SELECT phone FROM users WHERE phone = $1 LIMIT 1", phone)
	if dbErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func (ur *UserRepository) CheckUserIdentityExists(identity string) bool {
	rows, dbErr := ur.Db.Query("SELECT identity FROM users WHERE identity = $1 LIMIT 1", identity)
	if dbErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var err error
	var rows *sql.Rows
	rows, err = ur.Db.Query("SELECT * FROM users WHERE email = $1 LIMIT 1", email)
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
		return nil, appErr.NotFoundError("No user found.")
	}
	return &users[0], nil
}

func (ur *UserRepository) GetUserByPhone(phone string) (*model.User, error) {
	var err error
	var rows *sql.Rows
	rows, err = ur.Db.Query("SELECT * FROM users WHERE phone = $1 LIMIT 1", phone)
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
		return nil, appErr.NotFoundError("User not found.")
	}
	return &users[0], nil
}

func (ur *UserRepository) CreateTenantUser(ctx context.Context, firstName string, lastName string, email string, password string, birthdate *time.Time, gender *enum.Gender, phone string, identity string, address string, postalCode string, profileImage string) (*model.User, error) {
	var err error
	var tx *sql.Tx
	tx, err = ur.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var user model.User
	query := "INSERT INTO users (email, password, role, phone, identity, first_name, last_name, gender, birthdate, postal_code, address, profile_image, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING *"
	if err = tx.QueryRowContext(ctx, query , email, password, enum.ROLE_TENANT, phone, identity, firstName, lastName, gender, birthdate, postalCode, address, profileImage, enum.USER_ACTIVE).Scan(
		&user.Id,
		&user.Email,
		&user.Password, 
		&user.Role,
		&user.Phone, 
		&user.Identity,
		&user.FirstName,
		&user.LastName,
		&user.Gender,
		&user.Birthdate,
		&user.PostalCode,
		&user.Address,
		&user.ProfileImage,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) CreateAdminUser(email string, password string, phone string, identity string, firstName string, lastName string, gender *enum.Gender, birthdate *time.Time, postalCode string, address string, profileImage string) error {
	var err error
	_, err = ur.Db.Exec("INSERT INTO users (email, password, role, phone, identity, first_name, last_name, gender, birthdate, postal_code, address, profile_image, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", email, password, enum.ROLE_ADMIN, phone, identity, firstName, lastName, gender, birthdate, postalCode, address, profileImage, enum.USER_ACTIVE)

	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) CheckIfAdminUserAlreadyExists() (bool, error) {
	rows, rowErr := ur.Db.Query("SELECT id FROM users WHERE role = $1", enum.ROLE_ADMIN)
	if rowErr != nil {
		return false, rowErr
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}


func (ur *UserRepository) ListAllUsers() ([]model.User, error) {
	var rows *sql.Rows
	var err error
	rows, err = ur.Db.Query("SELECT * FROM users ORDER BY id DESC")
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
		return nil, appErr.NotFoundError("No user found.")
	}
	return users, nil
}


func (ur *UserRepository) GetUserById(id int64) (*model.User, error) {
	var rows *sql.Rows
	var err error
	rows, err = ur.Db.Query("SELECT * FROM users WHERE id = $1 LIMIT 1", id)
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
		return nil, appErr.NotFoundError("No user found.")
	}
	return &users[0], nil
}


func (ur *UserRepository) UpdateUser(ctx context.Context, userId int64, firstName *string, lastName *string, birthdate *time.Time, gender *enum.Gender, phone *string, identity *string, address *string, postalCode *string, profileImage *string) (*model.User, error) {
	var err error
	var tx *sql.Tx

	tx, err = ur.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query := "UPDATE users SET first_name = $1, last_name = $2, birthdate = $3, gender = $4, phone = $5, identity = $6, address = $7, postal_code = $8, profile_image = $9 WHERE id = $10"

	var user model.User
	if err = tx.QueryRowContext(ctx, query, firstName, lastName, birthdate, gender, phone, identity, address, postalCode, profileImage, userId).Scan(
		&user.Id,
		&user.Email,
		&user.Password, 
		&user.Role,
		&user.Phone, 
		&user.Identity,
		&user.FirstName,
		&user.LastName,
		&user.Gender,
		&user.Birthdate,
		&user.PostalCode,
		&user.Address,
		&user.ProfileImage,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) ChangeUserPassword(ctx context.Context, userId int64, newPassword string) (error) {
	var err error
	var tx *sql.Tx
	tx, err = ur.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE users SET password = $1 WHERE id = $2", newPassword, userId)
	if err != nil {
		return nil
	}
	return tx.Commit()
}