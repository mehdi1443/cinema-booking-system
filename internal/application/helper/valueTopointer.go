package helper

import "time"

// timeToPtr returns a pointer to t if it is non-zero, otherwise nil.
func TimeToPtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func HideValue[T comparable](val T, hide bool) *T {
	if hide {
		return nil
	}
	return &val
}
