package enum

import "errors"

type CinemaStatus string

const (
	CinemaStatusActive   CinemaStatus = "ACTIVE"
	CinemaStatusInactive CinemaStatus = "INACTIVE"
)

func (s CinemaStatus) IsValid() bool {
	switch s {
	case CinemaStatusActive, CinemaStatusInactive:
		return true
	}
	return false
}
