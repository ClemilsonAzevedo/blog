package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/clemilsonazevedo/blog/internal/dto/request"
	"github.com/clemilsonazevedo/blog/internal/dto/response"
	"github.com/clemilsonazevedo/blog/internal/http/auth"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var data request.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Cannot decode body", http.StatusInternalServerError)
		return
	}

	if data.Email == "" || data.UserName == "" || data.Password == "" {
		http.Error(w, "You need provide all credentials", http.StatusBadRequest)
		return
	}

	if len(data.Password) < 8 {
		http.Error(w, "Password need 8 or more characters", http.StatusBadRequest)
		return
	}

	existingUser, err := uc.service.GetUserByEmail(data.Email)
	if err != nil {
		http.Error(w, "Cannot search per user", http.StatusInternalServerError)
		return
	}

	if existingUser.Email == data.Email {
		http.Error(w, "User Already Exists", http.StatusConflict)
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		http.Error(w, "Cannot hash password", http.StatusInternalServerError)
		return
	}

	user := entities.User{
		UserName: data.UserName,
		Email:    data.Email,
		Password: hashedPassword,
		Role:     enums.Reader,
	}
	if err := uc.service.CreateUser(&user); err != nil {
		http.Error(w, "Cannot create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User Created has successfully"))
}

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
		http.Error(w, "Email or Password is incorrect", http.StatusBadRequest)
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

func (uc *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println(id)
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if err := uuid.Validate(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uc.service.GetUserByID(uuid.MustParse(id))
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

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	userId := uuid.MustParse(id)

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

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	userId := uuid.MustParse(id)
	if err := uc.service.DeleteUser(userId); err != nil {
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
