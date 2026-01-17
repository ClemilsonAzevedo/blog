package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/domain/exceptions"
	"github.com/clemilsonazevedo/blog/internal/dto/request"
	"github.com/clemilsonazevedo/blog/internal/dto/response"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/clemilsonazevedo/blog/pkg"
	"github.com/clemilsonazevedo/blog/tools"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type PostController struct {
	service *service.PostService
}

func NewPostController(service *service.PostService) *PostController {
	return &PostController{
		service: service,
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Creates a new blog post (Author role required)
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body request.PostCreate true "Post creation data"
// @Success 201 {string} string "Post created"
// @Failure 400 {string} string "You need set Content and AuthorId to create a post"
// @Failure 500 {string} string "Cannot create post"
// @Security CookieAuth
// @Router /post [post]
func (pc *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var postDTO request.PostCreate
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&postDTO); err != nil {
		exceptions.BadRequest(w, err, "Cannot Decode Body", nil)
		return
	}

	if dec.More() {
		exceptions.BadRequest(w, errors.New("Request Error"), "Multiple JSON values not allowed", &postDTO)
		return
	}

	authorId, err := pkg.ParseULID(postDTO.AuthorId.String())
	if err != nil {
		exceptions.BadRequest(w, errors.New("Request Error"), "Cannot parse Author ID", &postDTO)
	}

	contextUser, ok := r.Context().Value("user").(*entities.User)
	if !ok {
		exceptions.Unauthorized(w, "unauthorized")
		return
	}

	if contextUser.ID != authorId {
		exceptions.Unauthorized(w, "Cannot create post with other author id")
		return
	}

	if postDTO.AuthorId.String() == "" || postDTO.Content == "" || postDTO.Title == "" {
		exceptions.BadRequest(w, errors.New("Request Error"), "You need set all params to create a post", postDTO)
		return
	}

	postId, err := pkg.NewULID()
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot Generate ULID to this Post", reqId)
		return
	}

	slug, err := pc.service.GenerateUniqueSlug(postDTO.Title)
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot Generate Slug to this Post", reqId)
		return
	}

	Post := entities.Post{
		ID:       postId,
		Title:    postDTO.Title,
		Slug:     slug,
		Content:  postDTO.Content,
		AuthorId: authorId,
	}

	if err := pc.service.CreatePost(&Post); err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot create Post", reqId)
		return
	}

	response.CreatedPost(w, postId, authorId)
}

// CreatePostWithAi godoc
// @Summary Create a post with AI-generated title and hashtags
// @Description Creates a new blog post with AI-generated title and hashtags from content (Author role required)
// @Tags Posts, AI
// @Accept json
// @Produce json
// @Param request body request.AiPostCreate true "AI post creation data"
// @Success 201 {string} string "Post created"
// @Failure 400 {string} string "You need set Content and AuthorId to create a post"
// @Failure 500 {string} string "Cannot create post"
// @Security CookieAuth
// @Router /post-with-ai [post]
func (pc *PostController) CreatePostWithAi(w http.ResponseWriter, r *http.Request) {
	var aiPostDTO request.AiPostCreate
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&aiPostDTO); err != nil {
		exceptions.BadRequest(w, err, "Cannot Decode Body", nil)
		return
	}

	if dec.More() {
		exceptions.BadRequest(w, errors.New("Request Error"), "Multiple JSON values not allowed", &aiPostDTO)
		return
	}

	authorId, err := pkg.ParseULID(aiPostDTO.AuthorId.String())
	if err != nil {
		exceptions.BadRequest(w, errors.New("Request Error"), "Cannot parse Author ID", &aiPostDTO)
	}

	contextUser, ok := r.Context().Value("user").(*entities.User)
	if !ok {
		exceptions.Unauthorized(w, "unauthorized")
		return
	}

	if contextUser.ID != authorId {
		exceptions.Unauthorized(w, "Cannot create post with other author id")
		return
	}

	if aiPostDTO.AuthorId.String() == "" || aiPostDTO.Content == "" {
		exceptions.BadRequest(w, errors.New("Request Error"), "You need set all params to create a post", aiPostDTO)
		return
	}

	aiPostId, err := pkg.NewULID()
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot Generate ULID to this Post", reqId)
		return
	}

	aiRes := tools.GeneratePropsOfContent(aiPostDTO.Content)

	slug, err := pc.service.GenerateUniqueSlug(aiRes.Title)
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot Generate Slug to this Post", reqId)
		return
	}

	Post := entities.Post{
		ID:       aiPostId,
		Title:    aiRes.Title,
		Slug:     slug,
		Content:  pkg.GeneratePostContent(aiPostDTO.Content, aiRes.Hashtags),
		AuthorId: authorId,
	}

	if err := pc.service.CreatePost(&Post); err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot create this post", reqId)
		return
	}

	response.CreatedPost(w, aiPostId, authorId)
}

// GetPostById godoc
// @Summary Get post by ID
// @Description Retrieves a single post by its ULID
// @Tags Posts
// @Produce json
// @Param id path string true "Post ULID"
// @Success 200 {object} response.PostResponse
// @Failure 400 {string} string "ID is required"
// @Failure 400 {string} string "Cannot parse Id of Post"
// @Failure 500 {string} string "Error retrieving post"
// @Router /post/{id} [get]
func (pc *PostController) GetPostById(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("postId")
	if postIdStr == "" {
		exceptions.BadRequest(w, errors.New("Request Error"), "You need to provide Post Id", postIdStr)
		return
	}

	postId, err := pkg.ParseULID(postIdStr)
	if err != nil {
		exceptions.BadRequest(w, err, "Cannot parse Id of Post", postId)
		return
	}

	post, err := pc.service.GetPostByID(postId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			exceptions.NotFound(w, err, fmt.Sprintf("Post with id %v not found", postId))
			return
		}
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot get this post", reqId)
		return
	}

	postObj := response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Slug:      post.Slug,
		Content:   post.Content,
		AuthorId:  post.AuthorId,
		Likes:     post.Likes,
		Dislikes:  post.Dislikes,
		Views:     post.Views,
		CreatedAt: post.CreatedAt,
	}

	response.ShowPost(w, postObj)
}

// GetAllPosts godoc
// @Summary Get all posts
// @Description Retrieves all blog posts

// GetPaginatedPosts godoc
// @Summary Get paginated posts
// @Description Retrieves a paginated list of blog posts
// @Tags Posts
// @Produce json
// @Param page query int true "Page number (default: 1)"
// @Param limit query int true "Number of posts per page (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Returns data array and meta object with pagination info"
// @Failure 400 {string} string "Page is required"
// @Failure 400 {string} string "Limit is required"
// @Failure 500 {string} string "Error retrieving posts"
// @Router /posts [get]
func (c *PostController) GetPaginatedPosts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageStr == "" || limitStr == "" {
		exceptions.BadRequest(w, errors.New("Request Error"), "You need provide values to page and limit", map[string]string{
			"page":  pageStr,
			"limit": limitStr,
		})
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot convert page into int number", reqId)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot convert limit into int number", reqId)
		return
	}

	if page <= 0 || limit <= 0 || limit > 25 {
		page = 1
		limit = 10
	}

	posts, total, err := c.service.GetPaginatedPosts(page, limit)
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot get paginated Posts", reqId)
		return
	}

	postsObj := make([]response.PostResponse, len(posts))
	for i := range len(posts) {
		p := posts[i]
		postsObj[i] = response.PostResponse{
			ID:        p.ID,
			Title:     p.Title,
			Slug:      p.Slug,
			Content:   p.Content,
			Views:     p.Views,
			AuthorId:  p.AuthorId,
			CreatedAt: p.CreatedAt,
		}
	}

	response.ListPosts(w, postsObj, page, limit, int(total))
}

// UpdatePost godoc
// @Summary Update a post
// @Description Updates an existing blog post (Author role required)
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ULID"
// @Param request body request.PostUpdate true "Post update data"
// @Success 200 {object} response.PostResponse
// @Failure 400 {string} string "ID is required"
// @Failure 500 {string} string "Error updating post"
// @Security CookieAuth
// @Router /post/{id} [put]
func (pc *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("postId")
	if postIdStr == "" {
		exceptions.BadRequest(w, errors.New("Request Error"), "You need Provide Post Id on route", postIdStr)
		return
	}

	postId, err := pkg.ParseULID(postIdStr)
	if err != nil {
		exceptions.BadRequest(w, err, "Cannot Parse Comment Id", postId)
		return
	}

	existingPost, err := pc.service.GetPostByID(postId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			exceptions.NotFound(w, err, "Post Does not exists")
			return
		}

		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot get post", reqId)
		return
	}

	var updatePostDTO request.PostUpdate
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&updatePostDTO); err != nil {
		exceptions.BadRequest(w, err, "Cannot Decode Body", nil)
		return
	}

	if dec.More() {
		exceptions.BadRequest(w, errors.New("Request Error"), "Multiple JSON values not allowed", &updatePostDTO)
		return
	}

	if updatePostDTO.Title == "" {
		updatePostDTO.Title = existingPost.Title
	}

	if len(updatePostDTO.Title) >= 1 {
		slug, err := pc.service.GenerateUniqueSlug(updatePostDTO.Title)
		if err != nil {
			reqId := middleware.GetReqID(r.Context())
			exceptions.InternalError(w, err, "Cannot Generate Slug to this Post", reqId)
			return
		}

		existingPost.Slug = slug
	}

	if updatePostDTO.Content == "" {
		updatePostDTO.Content = existingPost.Content
	}

	postObj := entities.Post{
		ID:        existingPost.ID,
		AuthorId:  existingPost.AuthorId,
		Title:     updatePostDTO.Title,
		Content:   updatePostDTO.Content,
		Slug:      existingPost.Slug,
		Author:    existingPost.Author,
		Likes:     existingPost.Likes,
		Views:     existingPost.Views,
		Dislikes:  existingPost.Dislikes,
		CreatedAt: existingPost.CreatedAt,
	}

	if err := pc.service.UpdatePost(&postObj); err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot update this post", reqId)
		return
	}

	response.OK(w, "Comment updated with success", postObj)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Deletes an existing blog post (Author role required)
// @Tags Posts
// @Param id path string true "Post ULID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "ID is required"
// @Failure 500 {string} string "Error deleting post"
// @Security CookieAuth
// @Router /post/{id} [delete]
func (pc *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("postId")
	if postIdStr == "" {
		exceptions.BadRequest(w, errors.New("Request Error"), "You need Provide Post Id on route", postIdStr)
		return
	}

	postId, err := pkg.ParseULID(postIdStr)
	if err != nil {
		exceptions.BadRequest(w, err, "Cannot Parse Comment Id", postId)
		return
	}

	existingPost, err := pc.service.GetPostByID(postId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			exceptions.NotFound(w, err, "Post Does not exists")
			return
		}

		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot get post", reqId)
		return
	}

	if err := pc.service.DeletePost(existingPost.ID); err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot delete this post", reqId)
		return
	}

	response.DeletedPost(w, existingPost.ID)
}
