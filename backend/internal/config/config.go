package config

import (
	"os"
	"strconv"
	"time"
)

// Config 集中管理全部服务配置，所有配置均通过环境变量注入。
type Config struct {
	ServerPort   string
	RunMode      string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret    string
	JWTExpire    time.Duration
	APIKeySecret string
	RedisAddr    string
	RedisPass    string
	RedisDB      int
	RateLimit    int
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOUseSSL    bool
	MinIOBucket    string
}

// Load 从环境变量读取配置，并为缺失项提供默认值。
func Load() *Config {
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		RunMode:        getEnv("GIN_MODE", "release"),
		DBHost:         getEnv("DB_HOST", "127.0.0.1"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "rxcheck_user"),
		DBPassword:     getEnv("DB_PASSWORD", "rxcheck_pwd"),
		DBName:         getEnv("DB_NAME", "rxcheck_db"),
		JWTSecret:      getEnv("JWT_SECRET", "change_me_to_a_long_random_string"),
		JWTExpire:      time.Duration(getEnvInt("JWT_EXPIRE_HOURS", 24)) * time.Hour,
		APIKeySecret:   getEnv("API_KEY_SECRET", "change_me_to_a_long_random_string"),
		RedisAddr:      getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPass:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getEnvInt("REDIS_DB", 0),
		RateLimit:      getEnvInt("RATE_LIMIT_PER_SECOND", 200),
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "127.0.0.1:47020"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "rxcheck_minio"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "rxcheck_minio_pwd"),
		MinIOUseSSL:    getEnvBool("MINIO_USE_SSL", false),
		MinIOBucket:    getEnv("MINIO_BUCKET", "rxcheck-reports"),
	}
}

// DSN 返回 PostgreSQL 连接串。
func (c *Config) DSN() string {
	return "host=" + c.DBHost + " port=" + c.DBPort + " user=" + c.DBUser +
		" password=" + c.DBPassword + " dbname=" + c.DBName +
		" sslmode=disable TimeZone=Asia/Shanghai"
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}
