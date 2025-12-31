package request

type PostCreate struct {
	UserID  int    `json:"user_id" binding:"required,min=1"`
	Title   string `json:"title" binding:"required,min=2,max=100"`
	Content string `json:"content" binding:"required,min=2,max=1000"`
}

type PostUpdate struct {
	ID      int    `json:"id" binding:"required,min=1"`
	Title   string `json:"title" binding:"required,min=2,max=100"`
	Content string `json:"content" binding:"required,min=2,max=1000"`
}

type PostDelete struct {
	ID int `json:"id" binding:"required,min=1"`
}
