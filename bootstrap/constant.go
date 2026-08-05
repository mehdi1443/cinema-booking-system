package bootstrap

import (
	"encoding/json"
	"fmt"
	"time"
)

const DEBUG_MODE_FOR_DEVELOPING bool = true

type Constant struct {
	Database DBConst    `json:"database"`
	User     SuperUser  `json:"user"`
	Redis    RedisConst `json:"redis"`
}

type DBConst struct {
	MaxOpenDbConn int32         `json:"max_open_db_conn"`
	MaxIdleDbConn time.Duration `json:"max_idle_db_conn"`
	MaxDbLifeTime time.Duration `json:"max_db_life_time"`
}

type RedisConst struct {
	PhoneOtpSignature       string `json:"phone_otp_signature"`
	VerifiedOtpSignature    string `json:"verified_otp_signature"`
	RefreshTokenSignature   string `json:"refresh_token_signature"`
	UserSignature           string `json:"user"`
	UserPermissionSignature string `json:"user_permission_signature"`
}

type SuperUser struct {
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Wallet    int64  `json:"wallet"`
}

func NewConstant() *Constant {
	return &Constant{
		Database: DBConst{
			MaxOpenDbConn: 10,
			MaxIdleDbConn: 5 * time.Minute,
			MaxDbLifeTime: 5 * time.Minute,
		},
		User: SuperUser{
			Phone:     "09118991443",
			Email:     "mehdi1384sharifipoor@gmail.com",
			Wallet:    2_500_000,
			Password:  "123456",
			FirstName: "Mehdi",
			LastName:  "Sharifipour",
		},
		Redis: RedisConst{
			PhoneOtpSignature:       "otp:phone_email:",
			VerifiedOtpSignature:    "otp:phone_email:verified:",
			RefreshTokenSignature:   "token:refresh:",
			UserSignature:           "user",
			UserPermissionSignature: "permissions",
		},
	}
}

func (c *Constant) String() string {
	if c == nil {
		return "<nil>"
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("Constant{Database:%+v, User:%+v, Redis:%+v}",
			c.Database, c.User, c.Redis)
	}

	return string(data)
}
