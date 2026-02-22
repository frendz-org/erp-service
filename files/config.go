package files

import "time"

type Config struct {
	BatchSize     int
	StaleClaimAge time.Duration
}

func DefaultConfig() Config {
	return Config{
		BatchSize:     50,
		StaleClaimAge: 5 * time.Minute,
	}
}
