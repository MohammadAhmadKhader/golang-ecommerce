package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBName                 string
	DBAddress              string
	JWTExpirationInSeconds string
	JWTSecret              string
	Env 				   string
}

var Envs = initConfig()

func initConfig() Config {
	loadEnvFile()

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", ":8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "dbname"),
		JWTExpirationInSeconds: getEnv("JWTExpirationInSeconds", string(3600*24*7)),
		JWTSecret: getEnv("JWTSecret","fallback value"),
		Env: getEnv("env", "production"),
	}
}

var envs map[string]string

func getEnv(key, fallback string) string {
	if envs[key] != "" {
		return envs[key]
	}

	return fallback
}

func loadEnvFile() {
	godotenv.Load()
	envVariables, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}
	envs = envVariables
}
