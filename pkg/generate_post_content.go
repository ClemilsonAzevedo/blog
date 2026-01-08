package pkg

import "fmt"

// TODO: Transformar em arquivo markdown
func GeneratePostContent(content string, hashtags []string) string {
	postHashtags := []string{}
	for _, val := range hashtags {
		postHashtags = append(postHashtags, val)
	}

	postContent := fmt.Sprintf("%s\n\n%v", content, postHashtags)
	return postContent
}
