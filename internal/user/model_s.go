package user

type User struct {
	ID       string `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
