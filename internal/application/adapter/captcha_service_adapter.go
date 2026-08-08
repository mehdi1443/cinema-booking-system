package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type CaptchaServiceInterface interface {
	Save(ctx context.Context, captchaId string, code string) error
	Get(ctx context.Context, captchaId string) (string, error)
	Delete(ctx context.Context, captchaId string) error
	Exist(ctx context.Context, captchaId string) (*bool, error)
	GenerateAndSave(ctx context.Context) (*dto.GenerateCaptcha, error)
	Verify(ctx context.Context, captchaId string, claim string) (*bool, error)
}
