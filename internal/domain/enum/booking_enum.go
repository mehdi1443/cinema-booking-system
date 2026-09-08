package enum

type HallStatus string

const (
	HallStatusActive   HallStatus = "ACTIVE"
	HallStatusInactive HallStatus = "INACTIVE"
)

func (s HallStatus) IsValid() bool {
	return s == HallStatusActive || s == HallStatusInactive
}

type SeatType string

const (
	SeatTypeStandard   SeatType = "STANDARD"
	SeatTypeVIP        SeatType = "VIP"
	SeatTypeWheelchair SeatType = "WHEELCHAIR"
)

func (s SeatType) IsValid() bool {
	return s == SeatTypeStandard || s == SeatTypeVIP || s == SeatTypeWheelchair
}

type SeatStatus string

const (
	SeatStatusActive   SeatStatus = "ACTIVE"
	SeatStatusInactive SeatStatus = "INACTIVE"
)

func (s SeatStatus) IsValid() bool {
	return s == SeatStatusActive || s == SeatStatusInactive
}

type ShowtimeStatus string

const (
	ShowtimeStatusDraft     ShowtimeStatus = "DRAFT"
	ShowtimeStatusPublished ShowtimeStatus = "PUBLISHED"
	ShowtimeStatusCancelled ShowtimeStatus = "CANCELLED"
	ShowtimeStatusCompleted ShowtimeStatus = "COMPLETED"
)

func (s ShowtimeStatus) IsValid() bool {
	return s == ShowtimeStatusDraft || s == ShowtimeStatusPublished ||
		s == ShowtimeStatusCancelled || s == ShowtimeStatusCompleted
}

type ShowtimeSeatStatus string

const (
	ShowtimeSeatAvailable ShowtimeSeatStatus = "AVAILABLE"
	ShowtimeSeatHeld      ShowtimeSeatStatus = "HELD"
	ShowtimeSeatSold      ShowtimeSeatStatus = "SOLD"
	ShowtimeSeatBlocked   ShowtimeSeatStatus = "BLOCKED"
)

type ReservationStatus string

const (
	ReservationPendingPayment ReservationStatus = "PENDING_PAYMENT"
	ReservationPaid           ReservationStatus = "PAID"
	ReservationExpired        ReservationStatus = "EXPIRED"
	ReservationCancelled      ReservationStatus = "CANCELLED"
)
