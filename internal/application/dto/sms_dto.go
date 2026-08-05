package dto

type SendSMSRequest struct {
	Phone string `json:"phone" binding:"required,iranian-phone"`
	Text  string `json:"text"`
}
