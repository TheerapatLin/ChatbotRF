package utils

// ValidatePagination validates and normalizes pagination parameters
// Returns normalized (limit, offset) values
func ValidatePagination(limit, offset, maxLimit int) (int, int) {
	// Set default limit if not provided
	if limit <= 0 {
		limit = 20
	}

	// Enforce max limit
	if limit > maxLimit {
		limit = maxLimit
	}

	// Ensure offset is non-negative
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}