package dtos

type UserCreate struct {
	UserName     string
	UserPassword string
	UserEmail    string
}

type UserCreated struct {
	UserId int
}

type Authorize struct {
	Username string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UserInfo struct {
	UserID    int
	Username  string
	UserEmail string
}
