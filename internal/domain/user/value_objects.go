package user

import (
	"net/mail"
	"strings"
	"unicode"
	"unicode/utf8"
)

type UserName struct{ value string }

func CreateUserName(value string) (UserName, error) {
	if utf8.RuneCountInString(value) < 4 {
		return UserName{}, &ValidationError{Field: "username", Reason: "Username is too short"}
	} else if utf8.RuneCountInString(value) > 30 {
		return UserName{}, &ValidationError{Field: "username", Reason: "Username is too long"}
	} else if strings.TrimSpace(value) == "" {
		return UserName{}, &ValidationError{Field: "username", Reason: "Empty username"}
	}

	return UserName{value: value}, nil
}

func (name UserName) String() string {
	return name.value
}

type UserPassword struct{ value string }

func CreatePassword(value string, hasher Hasher) (UserPassword, error) {
	if utf8.RuneCountInString(value) < 6 {
		return UserPassword{}, &ValidationError{Field: "password", Reason: "Password is too short"}
	} else if utf8.RuneCountInString(value) > 30 {
		return UserPassword{}, &ValidationError{Field: "password", Reason: "Password is too long"}
	} else if strings.ToLower(value) == value {
		return UserPassword{}, &ValidationError{Field: "password", Reason: "Password required upper symbol"}
	} else if strings.ToUpper(value) == value {
		return UserPassword{}, &ValidationError{Field: "password", Reason: "Password required lower symbol"}
	}
	hasDigit := func() bool {
		for _, r := range value {
			if unicode.IsDigit(r) {
				return true
			}
		}
		return false
	}
	if !hasDigit() {
		return UserPassword{}, &ValidationError{Field: "password", Reason: "Password required digit"}
	}

	hashedPass, err := hasher.Hash(value)
	if err != nil {
		return UserPassword{}, &ValidationError{Field: "password", Reason: "Password required digit"}
	}

	return UserPassword{value: hashedPass}, nil
}

func (pass UserPassword) String() string {
	return pass.value
}

type UserEmail struct{ value string }

func CreateUserEmail(value string) (UserEmail, error) {
	if utf8.RuneCountInString(value) > 200 {
		return UserEmail{}, &ValidationError{Field: "email", Reason: "Email address is too long"}
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		return UserEmail{}, &ValidationError{Field: "email", Reason: "Incorrect email address"}
	}

	return UserEmail{value: value}, nil
}

func (email UserEmail) String() string {
	return email.value
}
