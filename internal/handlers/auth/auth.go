package auth

import (
	"campaign/internal/database"
	"campaign/internal/models"
	authservice "campaign/internal/services/auth"
	"campaign/internal/utils"
	"encoding/json"
	"net/http"
	"net/mail"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController interface {
	Signin(w http.ResponseWriter, r *http.Request)
	Signup(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	db *mongo.Database
}

func NewAuthHandler(db *mongo.Database) AuthController {
	return &authHandler{db: db}
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *authHandler) Signin(w http.ResponseWriter, r *http.Request) {
	reqBody := Login{}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	defer r.Body.Close()

	if err != nil {
		res := utils.WrapInResponse("error decoding request body", nil)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)
		return
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		res := utils.WrapInResponse("email and password are required", nil)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)
		return
	}

	dbM := database.NewDatabaseService(r.Context(), a.db, models.UsersCollection)

	authService := authservice.NewService(dbM)

	result, err := authService.Login(reqBody.Email, reqBody.Password)

	if err != nil {
		res := utils.WrapInResponse(err.Error(), nil)

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write(res)
		return
	}

	res := utils.WrapInResponse("login successful", result)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)

}

type Register struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Msisdn   string `json:"msisdn"`
	Name     string `json:"name"`
}

func (a *authHandler) Signup(w http.ResponseWriter, r *http.Request) {
	reqBody := Register{}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	defer r.Body.Close()

	if err != nil {
		res := utils.WrapInResponse("error decoding request body", nil)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)
		return
	}

	if reqBody.Email == "" || reqBody.Password == "" || reqBody.Msisdn == "" || reqBody.Name == "" {
		res := utils.WrapInResponse("email, password and msisdn is required", nil)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)
		return
	}

	if len(reqBody.Password) < 6 {
		res := utils.WrapInResponse("password must be at least 6 characters", nil)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)
		return
	}

	if len(reqBody.Msisdn) < 10 {
		res := utils.WrapInResponse("msisdn must be at least 10 characters", nil)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)
		return
	}

	_, err = mail.ParseAddress(reqBody.Email)

	if err != nil {

		res := utils.WrapInResponse("invalid email address", nil)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)
		return
	}

	if len(reqBody.Name) < 3 {
		res := utils.WrapInResponse("name must be at least 3 characters", nil)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)
		return
	}

	dbM := database.NewDatabaseService(r.Context(), a.db, models.UsersCollection)

	authService := authservice.NewService(dbM)

	err = authService.Register(reqBody.Name, reqBody.Email, reqBody.Password, reqBody.Msisdn)

	if err != nil {
		res := utils.WrapInResponse(err.Error(), nil)

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write(res)
		return
	}

	res := utils.WrapInResponse("account created successfully", nil)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}
