package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func InitializeLogger() {
	currentTime := time.Now()
	date := currentTime.Format("20060102")
	path := fmt.Sprintf("%s/%s-%s.%s", os.Getenv("LOG_PATH"), os.Getenv("LOG_PREFIX"), date, os.Getenv("LOG_EXT"))

	log.SetFormatter(&log.JSONFormatter{})

	if err := os.MkdirAll(filepath.Dir(path), 0770); err == nil {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.SetOutput(os.Stdout)
			return
		}
		mw := io.MultiWriter(os.Stdout, file)
		log.SetOutput(mw)
	} else {
		log.SetOutput(os.Stdout)
	}

}

func LogInfo(logtype, message string) {
	log.WithFields(log.Fields{
		"app_name":    os.Getenv("APP_NAME"),
		"app_version": os.Getenv("APP_VERSION"),
		"log_type":    logtype,
	}).Info(message)
}

func LogInfoWithData(logType string, body interface{}, method, route, msg string) {
	log.WithFields(log.Fields{
		"app_name":    os.Getenv("APP_NAME"),
		"app_version": os.Getenv("APP_VERSION"),
		"method":      method,
		"request":     body,
		"route":       route,
		"log_type":    logType,
	}).Info(msg)
}

func LogError(logtype, message string, request interface{}, route, method string) {
	log.WithFields(log.Fields{
		"app_name":    os.Getenv("APP_NAME"),
		"app_version": os.Getenv("APP_VERSION"),
		"log_type":    logtype,
		"request":     request,
		"method":      method,
		"route":       route,
	}).Error(message)
}
