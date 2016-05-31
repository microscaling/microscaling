package utils

import (
	"os"
	"strconv"
)

// EnvFl64 gets an environment variable and converts it into a float64
func EnvFl64(envVar string, defaultVal float64) (val float64) {
	s := os.Getenv(envVar)
	if s == "" {
		val = defaultVal
		return val
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Warningf("Bad value for %s, using default %f", envVar, defaultVal)
		val = defaultVal
	}

	return val
}
