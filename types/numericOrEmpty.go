package types

import (
	"encoding/json"
	"strconv"
)

type NumericOrEmpty float64

func (n *NumericOrEmpty) UnmarshalJSON(data []byte) error {
	// Try number directly
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*n = NumericOrEmpty(num)
		return nil
	}

	// Try string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		// "-" or empty → treat as zero
		if s == "-" {
			*n = 0
			return nil
		}

		// Try to parse string number
		if parsed, err := strconv.ParseFloat(s, 64); err == nil {
			*n = NumericOrEmpty(parsed)
			return nil
		}

		// Fallback
		*n = 0
		return nil
	}

	// Final fallback
	*n = 0
	return nil
}
