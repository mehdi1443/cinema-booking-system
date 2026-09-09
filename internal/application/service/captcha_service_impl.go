package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/pkg/captcha"
)

type captchaService struct {
	captchaRepo   cache.CaptchaInterface
	captchaConfig *bootstrap.CaptchaConfig
	captchaMan    captcha.CaptchaMan
	logger        domainlogger.Logger
}

func NewCaptchaService(
	captchaRepo cache.CaptchaInterface,
	captchaConfig *bootstrap.CaptchaConfig,
	captchaMan captcha.CaptchaMan,
	log domainlogger.Logger,
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
		// Never write the captcha answer to logs; it is a short-lived authentication secret.
		cs.logger.Error(ctx, "failed to save captcha to repository", domainlogger.Error(err), domainlogger.String("captcha_id", captchaId))
		return fmt.Errorf("failed to save captcha: %w", err)
	}
	return nil
}

func (cs *captchaService) Get(ctx context.Context, captchaId string) (string, error) {
	code, err := cs.captchaRepo.Get(ctx, captchaId)
	if err != nil {
		cs.logger.Warn(ctx, "failed to fetch captcha", domainlogger.Error(err), domainlogger.String("captcha_id", captchaId))
		return "", fmt.Errorf("failed to get captcha: %w", err)
	}
	return code, nil
}

func (cs *captchaService) Delete(ctx context.Context, captchaId string) error {
	if err := cs.captchaRepo.Delete(ctx, captchaId); err != nil {
		cs.logger.Error(ctx, "failed to delete captcha", domainlogger.Error(err), domainlogger.String("captcha_id", captchaId))
		return fmt.Errorf("failed to delete captcha: %w", err)
	}
	return nil
}

func (cs *captchaService) Exist(ctx context.Context, captchaId string) (*bool, error) {
	exists, err := cs.captchaRepo.Exist(ctx, captchaId)
	if err != nil {
		cs.logger.Error(ctx, "failed to check captcha existence", domainlogger.Error(err), domainlogger.String("captcha_id", captchaId))
		return nil, fmt.Errorf("failed to check captcha existence: %w", err)
	}

	return exists, nil
}

func (cs *captchaService) Verify(ctx context.Context, captchaId string, claim string) (*bool, error) {
	expected, err := cs.Get(ctx, captchaId)
	if err != nil {
		cs.logger.Warn(ctx, "failed to retrieve captcha", domainlogger.Error(err), domainlogger.String("captcha_id", captchaId))

		isValid := false
		return &isValid, err
	}

	if expected != claim {
		cs.logger.Warn(ctx, "captcha code mismatch", domainlogger.String("captcha_id", captchaId))

		isValid := false
		return &isValid, nil
	}

	if err := cs.Delete(ctx, captchaId); err != nil {
		cs.logger.Error(ctx, "failed to delete captcha after successful verification", domainlogger.Error(err), domainlogger.String("captcha_id", captchaId))
		return nil, err
	}

	cs.logger.Info(ctx, "captcha verified successfully", domainlogger.String("captcha_id", captchaId))

	isValid := true
	return &isValid, nil
}

func (cs *captchaService) GenerateAndSave(ctx context.Context) (*dto.GenerateCaptcha, error) {
	cid, b64, answer, err := cs.captchaMan.Generate()
	if err != nil {
		cs.logger.Error(ctx, "failed to render captcha code on b64", domainlogger.Error(err))
		return nil, err
	}

	err = cs.Save(ctx, cid, answer)
	if err != nil {
		cs.logger.Error(ctx, "failed to save generated captcha to cache", domainlogger.Error(err), domainlogger.String("captcha_id", cid))
		return nil, err
	}

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
	return nil, fmt.Errorf("unsupported captcha type %q", cs.captchaConfig.Type)
}
