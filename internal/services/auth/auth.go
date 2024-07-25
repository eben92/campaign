package authservice

import (
	"campaign/internal/database"
	"campaign/internal/models"
	"campaign/internal/utils/jwt"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	Login(email, password string) (LoginRes, error)
	Register(name, email, password, msisdn string) error
}

type service struct {
	db database.Database
}

func NewService(db database.Database) Service {
	return &service{db: db}
}

type LoginRes struct {
	models.User
	Token string `json:"access_token"`
}

func (s *service) Login(email, password string) (LoginRes, error) {
	result := LoginRes{}
	user := models.User{}

	s.db.SetCollection(models.UsersCollection)

	err := s.db.FindOne(bson.M{"email": email}, &user)

	if err != nil {
		slog.Error("Error finding user", "error", err)

		return result, errors.New("invalid email or password")
	}

	if user.Email == "" {
		slog.Error("User not found", "error", err)

		return result, errors.New("invalid email or password")
	}

	result.User = user
	token, err := jwt.GenereteJWT(jwt.AuthContext{
		Sub:  result.ID,
		Name: result.Name,
	})

	if err != nil {
		slog.Error("Error generating token", "error", err)

		return result, errors.New("error generating token")
	}

	result.Token = string(token)

	return result, nil

}

func (s *service) Register(name, email, password, msisdn string) error {
	s.db.SetCollection(models.UsersCollection)

	err := s.db.InsertOne(bson.M{"name": name, "email": email, "password": password, "msisdn": msisdn})

	if err != nil {
		slog.Error("Error inserting user", "error", err)

		return errors.New("error creating user. email already exists")
	}

	return nil
}
