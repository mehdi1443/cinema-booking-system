package service

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type smsService struct {
	filePath string
}

var mockSMSFileMu sync.Mutex

func NewSMSService() service_contract.SMSServiceProviderInterface {
	return NewMockSMSService("./SEND_SMS.txt")
}

func NewMockSMSService(filePath string) service_contract.SMSServiceProviderInterface {
	return &smsService{filePath: filePath}
}

func (s *smsService) SendSMS(req dto.SendSMSRequest) (*bool, error) {
	if s == nil || strings.TrimSpace(s.filePath) == "" || strings.TrimSpace(req.Phone) == "" || strings.TrimSpace(req.Text) == "" {
		return nil, fmt.Errorf("invalid mock SMS request")
	}
	text := fmt.Sprintf("[%s] Phone:%s - Text:%s\n", time.Now().Format(time.RFC3339), req.Phone, req.Text)

	mockSMSFileMu.Lock()
	defer mockSMSFileMu.Unlock()
	ffile, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("mock SMS file open failed: %v", err)
		return nil, err
	}
	defer ffile.Close()
	_, err = ffile.WriteString(text)
	if err != nil {
		return nil, err
	}

	success := true
	return &success, nil
}
