package api

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Dorrrke/golang-test-task/internal/domain/models"
	"github.com/Dorrrke/golang-test-task/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/pelletier/go-toml/v2"
)

type Storage interface {
	InsertUser(ctx context.Context, user models.User) (int, error)
	GetUserByID(ctx context.Context, userID int) (models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type Server struct {
	log       *slog.Logger
	storage   Storage
	timeout   time.Duration
	errorChan chan error
	respChan  chan []byte
}

func New(log *slog.Logger, storage Storage, timeout time.Duration) *Server {
	errorChan := make(chan error)
	respChan := make(chan []byte)
	return &Server{
		log:       log,
		storage:   storage,
		timeout:   timeout,
		errorChan: errorChan,
		respChan:  respChan,
	}
}

func (s *Server) AddUserHandler(res http.ResponseWriter, req *http.Request) {
	const op = "server.AddUserHandler"
	log := s.log.With(slog.String("op", op))
	log.Debug("Add user handler")

	dec := json.NewDecoder(req.Body)
	var user models.User
	if err := dec.Decode(&user); err != nil {
		http.Error(res, "request decoding error", http.StatusInternalServerError)
		return
	}
	log.Debug("Decode value:", slog.Any("user", user))
	uID, err := s.insertUser(user)
	if err != nil {
		log.Error("write data in db error")
		log.Debug("Error: " + err.Error())
		http.Error(res, "write data in db error", http.StatusInternalServerError)
		return
	}
	res.Write(uID)
	res.WriteHeader(http.StatusAccepted)
}

func (s *Server) GetUserHandler(res http.ResponseWriter, req *http.Request) {
	const op = "server.GetUserHandler"
	log := s.log.With(slog.String("op", op))

	log.Debug("Get user handler")

	userIDParam := chi.URLParam(req, "id")

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		http.Error(res, "Error getting user id", http.StatusInternalServerError)
		return
	}

	user, err := s.getUser(int(userID))
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("User not found")
			http.Error(res, "user not found", http.StatusNotFound)
			return
		}
		log.Error("get from db error: " + err.Error())
		http.Error(res, "get from db error", http.StatusInternalServerError)
		return
	}
	log.Debug("Run json enc")
	go s.jsonDecoder(user)
	log.Debug("Run xml enc")
	go s.xmlDecoder(user)
	log.Debug("Run toml enc")
	go s.tomlDecoder(user)

	temp := 0
	for {
		select {
		case resp := <-s.respChan:
			log.Debug("encode value:", slog.String("Encode", string(resp)))
			log.Debug("temp value", slog.Int("temp", temp))
			res.Write(resp)
			temp++
		case err := <-s.errorChan:
			log.Error("Error encoding to json: " + err.Error())
			http.Error(res, "Error encoding to json", http.StatusInternalServerError)
			return
		default:
			if temp >= 2 {
				res.Header().Set("Content-Type", "application/json")
				return
			}
		}
	}
}

func (s *Server) GetAllUsersHandler(res http.ResponseWriter, req *http.Request) {
	const op = "server.GetAllUsersHandler"
	log := s.log.With(slog.String("op", op))

	users, err := s.getAllUsers()
	if err != nil {
		if errors.Is(err, storage.ErrDataNotFound) {
			log.Error("No data in db")
			http.Error(res, "no users in db", http.StatusNoContent)
			return
		}
		log.Error("get from db error: " + err.Error())
		http.Error(res, "get from db error", http.StatusInternalServerError)
		return
	}
	for i := range users {
		log.Debug("Run json enc")
		go s.jsonDecoder(users[i])
		log.Debug("Run xml enc")
		go s.xmlDecoder(users[i])
		log.Debug("Run toml enc")
		go s.tomlDecoder(users[i])
	}

	temp := 0
	for {
		select {
		case resp := <-s.respChan:
			log.Debug("encode value:", slog.String("Encode", string(resp)))
			log.Debug("temp value", slog.Int("temp", temp))
			res.Write(resp)
			temp++
		case err := <-s.errorChan:
			log.Error("Error encoding to json: " + err.Error())
			http.Error(res, "Error encoding to json", http.StatusInternalServerError)
			return
		default:
			if temp >= len(users)-1 {
				return
			}
		}
	}
}

func (s *Server) insertUser(user models.User) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	uID, err := s.storage.InsertUser(ctx, user)
	if err != nil {
		return nil, err
	}
	s.log.Debug("User id from db ", slog.Int("uID", uID))
	userID := strconv.Itoa(uID)

	return []byte(userID), nil
}

func (s *Server) getUser(userID int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		// TODO: обработать ошибки ( отсутсвие данных )
		return models.User{
			Name:       "",
			Age:        0,
			Salary:     0,
			Occupation: "",
		}, err
	}
	return user, nil
}

func (s *Server) getAllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	users, err := s.storage.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Server) jsonDecoder(user models.User) {
	resp, err := json.MarshalIndent(user, "", "   ")
	if err != nil {
		s.errorChan <- err
		return
	}
	s.respChan <- resp
}

func (s *Server) xmlDecoder(user models.User) {
	resp, err := xml.MarshalIndent(user, "", "   ")
	if err != nil {
		s.errorChan <- err
		return
	}
	s.respChan <- resp
}

func (s *Server) tomlDecoder(user models.User) {
	resp, err := toml.Marshal(user)
	if err != nil {
		s.errorChan <- err
		return
	}
	s.respChan <- resp
}
