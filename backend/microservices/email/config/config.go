package config

import (
	"errors"
	"os"
	"strconv"
)

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func GetConfig() (*SMTPConfig, error) {
	host, exist := os.LookupEnv("SMTPHOST")
	if !exist {
		return nil, errors.New("ENV SMTPHOST do not exist")
	}

	portStr, exist := os.LookupEnv("SMTPPORT")
	if !exist {
		return nil, errors.New("ENV SMTPPORT do not exist")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	user, exist := os.LookupEnv("SMTPUSER")
	if !exist {
		return nil, errors.New("ENV SMTPUSER do not exist")
	}

	password, exist := os.LookupEnv("SMTPPASSWORD")
	if !exist {
		return nil, errors.New("ENV SMTPPASSWORD do not exist")
	}

	cfg := SMTPConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}

	return &cfg, nil
}
