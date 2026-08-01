package jwtsrv

import (
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application"
)

type jwtService struct {
	cfg *bootstrap.JWT
}

func NewTokenService(cfg *bootstrap.JWT) service_contract.JWTServiceInterface {
	return &jwtService{
		cfg: cfg,
	}
}

func (s *jwtService) GenerateAccessToken() {

}

func (s *jwtService) GenerateRefreshToken() {

}

func (s *jwtService) GetPayloads() {

}

func (s *jwtService) generateToken() {

}

func (s *jwtService) ValidateToken() {

}
