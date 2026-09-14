package schemas

type UserInfo struct {
	UserID    int    `json:"id"`
	Username  string `json:"username"`
	UserEmail string `json:"email"`
}
