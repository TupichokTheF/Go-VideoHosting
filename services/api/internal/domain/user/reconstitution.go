package user

type UserState struct {
	ID           int
	UserName     string
	UserEmail    string
	UserPassword string
}

func Reconstitute(state *UserState) *User {
	return &User{
		id:           state.ID,
		userName:     UserName{state.UserName},
		userEmail:    UserEmail{state.UserEmail},
		userPassword: UserPassword{state.UserPassword},
	}
}

func (u *User) State() UserState {
	return UserState{
		ID:           u.id,
		UserEmail:    u.userEmail.String(),
		UserName:     u.userName.String(),
		UserPassword: u.userPassword.String(),
	}
}
