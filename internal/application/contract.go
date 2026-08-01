package service_contract

type (
	JWTServiceInterface interface {
		GenerateAccessToken()
		ValidateToken()
		GetPayloads()
		GenerateRefreshToken()
	}

	UserServiceInterface interface {
	}
)
