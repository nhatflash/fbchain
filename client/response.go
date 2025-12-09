package client

import (
	enum "github.com/nhatflash/fbchain/enum"
	sql "database/sql"
	"time"
)

type UserResponse struct {
	Id				int64				`json:"id"`
	Email			string				`json:"email"`
	Role			*enum.Role			`json:"role"`
	Phone			sql.NullString		`json:"phone"`
	Identity		sql.NullString		`json:"identity"`
	FirstName		string				`json:"firstName"`
	LastName		string				`json:"lastName"`
	Gender			*enum.Gender		`json:"gender"`
	Birthdate		time.Time			`json:"birthdate"`
	PostalCode		sql.NullString		`json:"postalCode"`
	Address 		sql.NullString		`json:"address"`
	ProfileImage	sql.NullString		`json:"profileImage"`
	Status			*enum.UserStatus	`json:"status"`
}