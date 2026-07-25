package bootstrap

import "time"

type Constant struct {
	Database DBConst
	User     SuperUser
}

type DBConst struct {
	MaxOpenDbConn int32
	MaxIdleDbConn time.Duration
	MaxDbLifeTime time.Duration
}

type SuperUser struct {
	Phone     string
	Email     string
	Password  string
	FirstName string
	LastName  string
	Wallet    int64
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
	}
}
