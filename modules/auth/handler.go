package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"movie-platform/internal/httpx"
)

type Handler struct{ Service *Service }

func NewHandler(s *Service) *Handler { return &Handler{Service: s} }

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func decode(r *http.Request, dst any) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if e := d.Decode(dst); e != nil {
		return e
	}
	var extra any
	if e := d.Decode(&extra); e != io.EOF {
		return errors.New("request must contain one json object")
	}
	return nil
}
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var q registerRequest
	if decode(r, &q) != nil {
		httpx.Error(w, 400, "invalid_json", "invalid request body")
		return
	}
	u, e := h.Service.Register(q.Name, q.Email, q.Password)
	if e != nil {
		switch {
		case errors.Is(e, ErrEmailAlreadyExists):
			httpx.Error(w, 409, "account_conflict", "unable to create account with supplied details")
		case errors.Is(e, ErrInvalidInput):
			httpx.Error(w, 400, "invalid_input", "name, email or password is invalid")
		default:
			httpx.Error(w, 500, "internal_error", "internal server error")
		}
		return
	}
	httpx.JSON(w, 201, map[string]any{"user": u})
}
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var q loginRequest
	if decode(r, &q) != nil || strings.TrimSpace(q.Email) == "" || q.Password == "" {
		httpx.Error(w, 400, "invalid_input", "email and password are required")
		return
	}
	token, u, e := h.Service.Login(q.Email, q.Password)
	if e != nil {
		httpx.Error(w, 401, "invalid_credentials", "email or password is incorrect")
		return
	}
	httpx.JSON(w, 200, map[string]any{"access_token": token, "token_type": "Bearer", "expires_in": int64(h.Service.Config.JWTTTL.Seconds()), "user": u})
}
