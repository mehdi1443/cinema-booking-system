package comunication

type SMSInterFace interface {
	SendOTP(phone string, code string) error
}
