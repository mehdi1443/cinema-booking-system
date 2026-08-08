package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/pkg/captcha"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
)

// TODO : FIX LOGGER
type captchaService struct {
	captchaRepo   cache.CaptchaInterface
	captchaConfig *bootstrap.CaptchaConfig
	captchaMan    captcha.CaptchaMan
	logger        logger.Logger
}

func NewCaptchaService(
	captchaRepo cache.CaptchaInterface,
	captchaConfig *bootstrap.CaptchaConfig,
	captchaMan captcha.CaptchaMan,
	log logger.Logger,
) service_contract.CaptchaServiceInterface {
	return &captchaService{
		captchaRepo:   captchaRepo,
		captchaConfig: captchaConfig,
		captchaMan:    captchaMan,
		logger:        log,
	}
}

func (cs *captchaService) Save(ctx context.Context, captchaId string, code string) error {
	if err := cs.captchaRepo.Save(ctx, captchaId, code, cs.captchaConfig.TTLInSeconds()); err != nil {
		cs.logger.Error(logger.Internal, logger.Captcha, "failed to save captcha to repository", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
			"captcha_id":        captchaId,
			"code":              code,
		})
		return fmt.Errorf("failed to save captcha: %w", err)
	}
	return nil
}

func (cs *captchaService) Get(ctx context.Context, captchaId string) (string, error) {
	code, err := cs.captchaRepo.Get(ctx, captchaId)
	if err != nil {
		cs.logger.Warn(logger.Internal, logger.Captcha, "failed to fetch captcha", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
			"captcha_id":        captchaId,
		})
		return "", fmt.Errorf("failed to get captcha: %w", err)
	}
	return code, nil
}

func (cs *captchaService) Delete(ctx context.Context, captchaId string) error {
	if err := cs.captchaRepo.Delete(ctx, captchaId); err != nil {
		cs.logger.Error(logger.Internal, logger.Captcha, "failed to delete captcha", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
			"captcha_id":        captchaId,
		})
		return fmt.Errorf("failed to delete captcha: %w", err)
	}
	return nil
}

func (cs *captchaService) Exist(ctx context.Context, captchaId string) (*bool, error) {
	exists, err := cs.captchaRepo.Exist(ctx, captchaId)
	if err != nil {
		cs.logger.Error(logger.Internal, logger.Captcha, "failed to check captcha existence", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
			"captcha_id":        captchaId,
		})
		return nil, fmt.Errorf("failed to check captcha existence: %w", err)
	}

	return exists, nil
}

func (cs *captchaService) Verify(ctx context.Context, captchaId string, claim string) (*bool, error) {
	expected, err := cs.Get(ctx, captchaId)
	if err != nil {

		cs.logger.Warn(logger.Internal, logger.Captcha, "failed to retrieve captcha", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
			"captcha_id":        captchaId,
		})

		isValid := false
		return &isValid, err
	}

	if expected != claim {

		cs.logger.Warn(logger.Internal, logger.Captcha, "code mismatch", map[logger.ExtraKey]interface{}{
			"captcha_id": captchaId,
			"expected":   expected,
			"claim":      claim,
		})

		isValid := false
		return &isValid, nil
	}

	if err := cs.Delete(ctx, captchaId); err != nil {

		cs.logger.Error(logger.Internal, logger.Captcha, "failed to delete captcha after successful verification", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
			"captcha_id":        captchaId,
		})
		return nil, err
	}

	cs.logger.Info(logger.Internal, logger.Captcha, "captcha verified successfully", map[logger.ExtraKey]interface{}{
		"captcha_id": captchaId,
		"expected":   expected,
		"claim":      claim,
	})

	isValid := true
	return &isValid, nil
}

func (cs *captchaService) GenerateAndSave(ctx context.Context) (*dto.GenerateCaptcha, error) {
	cid, b64, answer, err := cs.captchaMan.Generate()
	if err != nil {
		cs.logger.Error(logger.Internal, logger.Captcha, "failed to render captcha code on b64", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	err = cs.Save(ctx, cid, answer)
	if err != nil {

		cs.logger.Error(logger.Internal, logger.Captcha, "failed to save generated captcha to cache", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
			"captcha_id":        cid,
		})

		return nil, err
	}

	//TODO : refactor
	switch cs.captchaConfig.Type {
	case "string", "math":
		return &dto.GenerateCaptcha{
			CaptchaId: cid,
			ImageCode: b64,
			AudioCode: "",
			ExpireAt:  time.Time.Add(time.Now(), cs.captchaConfig.TTLInSeconds()),
		}, nil

	case "audio":
		return &dto.GenerateCaptcha{
			CaptchaId: cid,
			AudioCode: b64,
			ImageCode: "",
			ExpireAt:  time.Time.Add(time.Now(), cs.captchaConfig.TTLInSeconds()),
		}, nil
	}
	panic("WTF")

}
