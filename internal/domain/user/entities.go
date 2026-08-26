package user

type User struct {
	id           int
	userName     UserName
	userEmail    UserEmail
	userPassword UserPassword
}

func New(inputUserName, inputUserEmail, inputUserPassword string, hasher Hasher) (*User, error) {
	userName, err := CreateUserName(inputUserName)
	if err != nil {
		return nil, err
	}

	pass, err := CreatePassword(inputUserPassword, hasher)
	if err != nil {
		return nil, err
	}

	userEmail, err := CreateUserEmail(inputUserEmail)
	if err != nil {
		return nil, err
	}

	return &User{
		userName:     userName,
		userPassword: pass,
		userEmail:    userEmail,
	}, nil
}

func (u *User) Username() UserName {
	return u.userName
}

func (u *User) Email() UserEmail {
	return u.userEmail
}

func (u *User) Password() UserPassword {
	return u.userPassword
}

func (u *User) ID() int {
	return u.id
}

func (u *User) VerifyPassword(raw string, hasher Hasher) bool {
	return hasher.Verify(raw, u.userPassword.value)
}
