package bootstrap

import (
	"encoding/json"
	"time"
)

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
	PhoneEmailOtpSignature  string `json:"phone_otp_signature"`
	VerifiedOtpSignature    string `json:"verified_otp_signature"`
	RefreshTokenSignature   string `json:"refresh_token_signature"`
	UserSignature           string `json:"user"`
	UserPermissionSignature string `json:"user_permission_signature"`
}

type SuperUser struct {
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"-"`
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
			// Bootstrap credentials are supplied through SUPER_ADMIN_* environment
			// variables; never keep deployable credentials in source code.
		},
		Redis: RedisConst{
			PhoneEmailOtpSignature:  "otp:phone_email:",
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
		return "Constant{redacted: true}"
	}

	return string(data)
}
