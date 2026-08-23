package schemas

type UserInfoSchema struct {
	UserID    int    `json:"id"`
	Username  string `json:"username"`
	UserEmail string `json:"email"`
}
