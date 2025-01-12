package util

// Constants for all supported currencies
const (
	PENDING   = "pending"
	COMPLETED = "completed"
)

// IsSupportedStatus returns true if the status is supported
func IsSupportedStatus(status string) bool {
	switch status {
	case PENDING, COMPLETED:
		return true
	}

	return false
}
