package validation

import (
	"errors"

	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

func ValidateAccount(account *model.Account) (Validation, error) {
	var v Validation

	if account.Username == "" || account.Email == "" || account.Password == "" {
		return v, errors.New("Missing data in request")
	}

	v.minLength("Username", account.Username, 3)
	v.maxLength("Username", account.Username, 50)
	v.isDuplicate("Email", "account", "email", account.Email, 0)
	v.isDuplicate("Username", "account", "username", account.Username, 0)
	v.isEmail("Email", account.Email)
	v.password(account.Password)

	sanitizeText(&account.Username)
	sanitizeText(&account.Email)
	sanitizeText(&account.Password)

	return v, nil
}

func ValidateAccountUpdate(account *model.Account) (Validation, error) {
	var v Validation

	if account.Username == "" || account.Email == "" || account.Password == "" {
		return v, errors.New("Missing data in request")
	}

	v.minLength("Username", account.Username, 3)
	v.maxLength("Username", account.Username, 50)
	v.isDuplicate("Email", "account", "email", account.Email, 1)
	v.isDuplicate("Username", "account", "username", account.Username, 1)
	v.isEmail("Email", account.Email)
	if account.Password != "" {
		v.password(account.Password)
	}

	sanitizeText(&account.Username)
	sanitizeText(&account.Email)
	sanitizeText(&account.Email)
	sanitizeText(&account.Password)

	return v, nil
}

func ValidatePassword(password *string) Validation {
	var v Validation

	v.password(*password)

	sanitizeText(password)

	return v
}
