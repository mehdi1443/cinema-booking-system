package comunication

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/mehdi1443/cinema-booking-system/internal/domain/comunication"
// )

// type fileSMS struct {
// 	filePath string
// }

// func NewFileSMS(path string) comunication.SMSInterFace {
// 	return &fileSMS{
// 		filePath: path,
// 	}
// }

// func (s *fileSMS) SendOTP(phone string, code string) error {
// 	message := fmt.Sprintf("Time: %s | Phone: %s | OTP: %s\n",
// 		time.Now().Format("2006-01-02 15:04:05"),
// 		phone,
// 		code,
// 	)

// 	f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return fmt.Errorf("could not open sms log file: %w", err)
// 	}
// 	defer f.Close()

// 	if _, err := f.WriteString(message); err != nil {
// 		return fmt.Errorf("could not write to sms log file: %w", err)
// 	}

// 	return nil
// }
