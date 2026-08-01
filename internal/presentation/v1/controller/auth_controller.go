package controller

import "github.com/gin-gonic/gin"

type AuthController struct {
	// authService *service.AuthService
}

func NewAuthController() AuthController {
	return AuthController{}
}

func (ac *AuthController) SendOTP(ctx *gin.Context) {

}

func (ac *AuthController) VerifyOTP(ctx *gin.Context) {

}

func (ac *AuthController) Register(ctx *gin.Context) {

}

func (ac *AuthController) Refresh(ctx *gin.Context) {

}

func (ac *AuthController) Login(ctx *gin.Context) {

}

func (ac *AuthController) Logout(ctx *gin.Context) {

}
