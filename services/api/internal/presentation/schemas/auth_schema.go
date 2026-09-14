package schemas

type UserCreated struct {
	UserID int `json:"user_id" example:"1"`
}

type CreateUser struct {
	Username string `json:"username" example:"maximEZ"`
	Email    string `json:"email" example:"maxim@mail.ru"`
	Password string `json:"password" example:"1Q2w3e"`
}

type Authorize struct {
	Username string `json:"username" example:"maximEZ"`
	Password string `json:"password" example:"1Q2w3e"`
}

type Token struct {
	AccessToken string `json:"access_token"`
}
