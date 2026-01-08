package response

type AiResponse struct {
	Title    string   `json:"title"`
	Hashtags []string `json:"hashtags"`
}
