package port

import (
	"net/http"
)

type TokenService interface {
	CreateToken(id, role string) (string, error)
	ParseFromRequest(r *http.Request) (map[string]string, error)
}
