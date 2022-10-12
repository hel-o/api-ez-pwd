package settings

import (
	"app-ez-pwd/internal/logger"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"os"
)

var Settings struct {
	SecretHex       string `json:"SECRET_HEX"`
	ApiHostPort     string `json:"API_HOST_PORT"`
	DatabaseURL     string `json:"DATABASE_URL"`
	CookieWebDomain string `json:"COOKIE_WEB_DOMAIN"`
	Debug           bool   `json:"DEBUG"`
}

func LoadConfiguration() {
	file, err := os.Open(os.Getenv("FILE_CONFIG"))
	if err != nil {
		logger.Logger.Error("invalid file", zap.Error(err))
		os.Exit(1)
	}
	defer file.Close()

	byteValues, _ := io.ReadAll(file)
	err = json.Unmarshal(byteValues, &Settings)
	if err != nil {
		logger.Logger.Error("invalid FILE_CONFIG file", zap.Error(err))
		os.Exit(1)
	}
}
