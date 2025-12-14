package initializer

import (
	"database/sql"
	"os"

	"github.com/nhatflash/fbchain/enum"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/security"
)

func CreateAdminUserIfNotExists(db *sql.DB) error {
	exist, existErr := repository.CheckIfAdminUserAlreadyExists(db)
	if existErr != nil {
		return existErr
	}
	if exist {
		return nil
	}
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	phone := os.Getenv("ADMIN_PHONE")
	identity := os.Getenv("ADMIN_IDENTITY")
	firstName := os.Getenv("ADMIN_FIRSTNAME")
	lastName := os.Getenv("ADMIN_LASTNAME")
	gender := enum.GENDER_MALE
	birthdate, bdErr := helper.ConvertToDate(os.Getenv("ADMIN_BIRTHDATE"))
	postalCode := os.Getenv("ADMIN_POSTALCODE")
	address := os.Getenv("ADMIN_ADDRESS")
	profileImage := os.Getenv("ADMIN_PROFILEIMAGE")
	if bdErr != nil {
		return bdErr
	}
	hashedPassword, hashErr := security.GenerateHashedPassword(password)
	if hashErr != nil {
		return hashErr
	}
	dbErr := repository.CreateAdminUser(email, hashedPassword, phone, identity, firstName, lastName, &gender, birthdate, postalCode, address, profileImage, db)
	if dbErr != nil {
		return dbErr
	}
	return nil
}
