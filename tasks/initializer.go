package tasks

import (
	"database/sql"
	"os"
	"fmt"
	"github.com/nhatflash/fbchain/enum"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/security"
	"time"
)

func CreateAdminUserIfNotExists(db *sql.DB) error {
	var err error
	var exist bool

	userRepository := repository.NewUserRepository(db)
	exist, err = userRepository.CheckIfAdminUserAlreadyExists()
	if err != nil {
		return err
	}

	if exist {
		fmt.Println("Admin account exists, skip creation!")
		return nil
	}

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	phone := os.Getenv("ADMIN_PHONE")
	identity := os.Getenv("ADMIN_IDENTITY")
	firstName := os.Getenv("ADMIN_FIRSTNAME")
	lastName := os.Getenv("ADMIN_LASTNAME")
	gender := enum.GENDER_MALE
	postalCode := os.Getenv("ADMIN_POSTALCODE")
	address := os.Getenv("ADMIN_ADDRESS")
	profileImage := os.Getenv("ADMIN_PROFILEIMAGE")

	var birthdate *time.Time
	birthdate, err = helper.ConvertToDate(os.Getenv("ADMIN_BIRTHDATE"))
	if err != nil {
		return err
	}
	var hashedPassword string
	hashedPassword, err = security.GenerateHashedPassword(password)
	if err != nil {
		return err
	}

	err = userRepository.CreateAdminUser(email, hashedPassword, phone, identity, firstName, lastName, &gender, birthdate, postalCode, address, profileImage)
	if err != nil {
		return err
	}
	return nil
}
