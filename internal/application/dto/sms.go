package dto

type SendSMSRequest struct {
	Phone string `json:"phone"`
	Text  string `json:"text"`
}
