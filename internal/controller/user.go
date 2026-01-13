package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/clemilsonazevedo/blog/internal/domain/exceptions"
	"github.com/clemilsonazevedo/blog/internal/dto/request"
	"github.com/clemilsonazevedo/blog/internal/dto/response"
	"github.com/clemilsonazevedo/blog/internal/http/auth"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/clemilsonazevedo/blog/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type User = entities.User
type UserService = service.UserService

type UserController struct {
	service *UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// CreateUser godoc
// @Summary Register a new user
// @Description Creates a new user account with the provided credentials
// @Tags Auth
// @Accept json
// @Produce plain
// @Param request body request.UserRegister true "User registration data"
// @Success 201 {string} string "User Created has successfully"
// @Failure 400 {string} string "You need provide all credentials"
// @Failure 400 {string} string "Password need 8 or more characters"
// @Failure 409 {string} string "User Already Exists"
// @Failure 500 {string} string "Cannot create user"
// @Router /register [post]
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var data request.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot decode body", reqId)
		return
	}

	if data.Email == "" || data.UserName == "" || data.Password == "" {
		exceptions.BadRequest(w, errors.New("Request Failed"), "You need provide all credentials", "")
		return
	}

	email, err := mail.ParseAddress(data.Email)
	if err != nil {
		exceptions.BadRequest(w, err, "You need Provide a valid email", email.Address)
		return
	}

	if len(data.Password) < 8 {
		exceptions.BadRequest(w, errors.New("Request Failed"), "Password need 8 or more characters", data.Password)
		return
	}

	existingUser, err := uc.service.GetUserByEmail(email.Address)
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot search per user", reqId)
		return
	}

	if existingUser.Email == data.Email {
		exceptions.Conflict(w, errors.New("Conflict"), "User Already Exists")
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot hash password", reqId)
		return
	}

	userId, err := pkg.NewULID()
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot Generate ULID to this User", reqId)
		return
	}

	user := entities.User{
		ID:       userId,
		UserName: data.UserName,
		Email:    data.Email,
		Password: hashedPassword,
		Role:     enums.Reader,
	}

	if err := uc.service.CreateUser(&user); err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot Create this user", reqId)
		return
	}

	response.CreatedUser(w, userId.String())
}

// LoginUser godoc
// @Summary Login user
// @Description Authenticates a user and sets a JWT token in an HTTP-only cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.UserLogin true "User login credentials"
// @Success 200 {string} string "Login successful - token set in cookie"
// @Failure 400 {string} string "Email and Password are Required"
// @Failure 400 {string} string "Email or Password is incorrect"
// @Header 200 {string} Set-Cookie "token=<jwt>; Path=/; HttpOnly; SameSite=Lax"
// @Router /login [post]
func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var data request.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data.Email == "" || data.Password == "" {
		http.Error(w, "Email and Password are Required", http.StatusBadRequest)
		return
	}

	authUser, err := uc.service.GetUserByEmail(data.Email)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		http.Error(w, "Email or Password is incorrectu", http.StatusBadRequest)
		return
	}

	isPasswordEquals := auth.CheckPassword(authUser.Password, data.Password)
	if !isPasswordEquals {
		http.Error(w, "Email or Password is incorrect", http.StatusBadRequest)
		return
	}

	token, exp, err := auth.GenerateJWT(*authUser, 24*7*time.Hour)
	if err != nil {
		http.Error(w, "Email or Password is incorrect", http.StatusBadRequest)
		return
	}

	age := int(exp - time.Now().Unix())
	age = max(age, 0)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(exp, 0),
		MaxAge:   age,
	})

	w.WriteHeader(http.StatusOK)
}

// Logout godoc
// @Summary Logout user
// @Description Clears the authentication cookie and logs out the user
// @Tags Auth
// @Produce json
// @Success 200 {object} response.UserLogout "Logout successful"
// @Security CookieAuth
// @Router /logout [post]
func (c *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"logout successful"}`))
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Retrieves a user by their ULID
// @Tags Users
// @Produce json
// @Param id path string true "User ULID"
// @Success 200 {object} response.UserByID
// @Failure 400 {string} string "ID is required"
// @Failure 500 {string} string "Error retrieving user"
// @Security CookieAuth
// @Router /user/{id} [get]
func (uc *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "id")
	if userIdStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	userId, err := pkg.ParseULID(userIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uc.service.GetUserByID(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response.UserByID{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Role:     user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Profile godoc
// @Summary Get current user profile
// @Description Retrieves the profile of the currently authenticated user
// @Tags Users
// @Produce json
// @Success 200 {object} response.UserByID
// @Failure 401 {string} string "unauthorized"
// @Security CookieAuth
// @Router /profile [get]
func (uc *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*entities.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response := response.UserByID{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Role:     user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetUserByEmail godoc
// @Summary Get user by email
// @Description Retrieves a user by their email address
// @Tags Users
// @Produce json
// @Param email query string true "User email"
// @Success 200 {object} response.UserByID
// @Failure 400 {string} string "Email is required"
// @Failure 500 {string} string "Error retrieving user"
// @Security CookieAuth
// @Router /user [get]
func (uc *UserController) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	user, err := uc.service.GetUserByEmail(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response.UserByID{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetUserByName godoc
// @Summary Get user by name
// @Description Retrieves a user by their username
// @Tags Users
// @Produce json
// @Param name query string true "Username"
// @Success 200 {object} response.UserByID
// @Failure 400 {string} string "Name is required"
// @Failure 500 {string} string "Error retrieving user"
// @Security CookieAuth
// @Router /user/name [get]
func (uc *UserController) GetUserByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	user, err := uc.service.GetUserByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response.UserByID{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieves a list of all users (Author role required)
// @Tags Users
// @Produce json
// @Success 200 {array} entities.User
// @Failure 500 {string} string "Error retrieving users"
// @Security CookieAuth
// @Router /users [get]
func (uc *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uc.service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// UpdateUser godoc
// @Summary Update user profile
// @Description Updates the currently authenticated user's profile
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ULID"
// @Param request body request.UserUpdate true "User update data"
// @Success 200 {object} response.UserByID
// @Failure 400 {string} string "ID is required"
// @Failure 400 {string} string "Cannot Parse String to ULID"
// @Failure 500 {string} string "Error updating user"
// @Security CookieAuth
// @Router /profile [put]
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "id")
	if userIdStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	userId, err := pkg.ParseULID(userIdStr)
	if err != nil {
		http.Error(w, "Cannot Parse String to ULID", http.StatusBadRequest)
		return
	}

	var dto request.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := entities.User{
		ID:       userId,
		UserName: dto.UserName,
		Email:    dto.Email,
		Role:     dto.Role,
	}

	if err := uc.service.UpdateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.UserByID{
		ID:       userId,
		UserName: user.UserName,
		Email:    user.Email,
	})
}

// DeleteUser godoc
// @Summary Delete user profile
// @Description Deletes the currently authenticated user's account
// @Tags Users
// @Produce json
// @Param id path string true "User ULID"
// @Success 200 {object} response.UserByID
// @Failure 400 {string} string "ID is required"
// @Failure 500 {string} string "Error deleting user"
// @Security CookieAuth
// @Router /profile [delete]
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "id")
	if userIdStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	userId, err := pkg.ParseULID(userIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.UserByID{
		ID:       userId,
		UserName: "",
		Email:    "",
	})
}
