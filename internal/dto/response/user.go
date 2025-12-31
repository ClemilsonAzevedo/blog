package response

type UserLogin struct {
	Token   string `json:"token"`
	Expires int64  `json:"exp"`
}

type UserLogout struct {
	Message string `json:"message"`
}
