package utils

import "os"

func GetEnv(key string) string {
	if env, exists := os.LookupEnv(key); exists {
		return env
	}

	return ""
}
