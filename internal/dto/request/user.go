package request

type UserRegister struct {
	Username string `json:"username" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
	Role     int    `json:"role" binding:"required,oneof=0 1 2"`
}
