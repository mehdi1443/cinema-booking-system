package auth

// import (
// 	"time"

// 	"github.com/mehdi1443/cinema-booking-system/internal/domain/comunication"
// 	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/redis"
// )

// type MyAuthServiceImp struct {
// 	OtpRepo     redis.OTPInterFace
// 	smsProvider comunication.SMSInterFace
// }

// func NewMyAuthServiceImp() *MyAuthServiceImp {
// 	return &MyAuthServiceImp{}
// }

// func (s *MyAuthServiceImp) SendOTP(phone string, code string) (string, error) {
// 	err := s.OtpRepo.Save(phone, code, 2*time.Minute)
// 	if err != nil {
// 		return "", err
// 	}

// 	message := "your code : " + code

// 	err = s.smsProvider.SendOTP(phone, message)
// 	if err != nil {
// 		return "", err
// 	}

// 	return code, nil
// }
