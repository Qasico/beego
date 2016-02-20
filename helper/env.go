package helper

import (
	"os"
	"strconv"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv(envPath string) (err error) {
	configFile := filepath.Join(envPath, ".env")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		configFile = filepath.Join(envPath, "../", ".env")
	}

	e := godotenv.Load(configFile)
	if e != nil {
		return e
	}

	return nil
}

func EnvString(key string, defaultVal string) (val string)  {
	if val = os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func EnvInt(key string, defaultVal int) (val int)  {
	p, _ := strconv.ParseInt(os.Getenv(key), 10, 32)
	if val = int(p); val != 0 {
		return val
	}

	return defaultVal
}

func EnvBool(key string, defaultVal bool) (val bool)  {
	if v := os.Getenv(key); v != "" {
		if v == "true" {
			return true
		} else {
			return false
		}
	}

	return defaultVal
}