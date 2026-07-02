package port

import (
	"net/http"
)

type TokenService interface {
	CreateToken(id string) (string, error)
	ParseFromRequest(r *http.Request) (string, error)
}
