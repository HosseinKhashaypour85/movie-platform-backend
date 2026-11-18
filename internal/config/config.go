package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port string
	AllowedOrigins []string
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	DBMaxOpen int
	DBMaxIdle int
	DBConnMaxLifetime time.Duration
	JWTSecret []byte
	JWTIssuer string
	JWTTTL time.Duration
}

func Load() (Config,error) {
	_ = godotenv.Load()
	c:=Config{
		Environment: value("APP_ENV","development"),
		Port:value("APP_PORT","8080"),
		AllowedOrigins:csv(value("APP_ALLOWED_ORIGINS","http://localhost:3000")),
		DBHost:value("DB_HOST","127.0.0.1"), DBPort:value("DB_PORT","3306"),
		DBUser:strings.TrimSpace(os.Getenv("DB_USER")), DBPassword:os.Getenv("DB_PASSWORD"),
		DBName:strings.TrimSpace(os.Getenv("DB_NAME")),
		DBMaxOpen:intValue("DB_MAX_OPEN_CONNS",25), DBMaxIdle:intValue("DB_MAX_IDLE_CONNS",10),
		DBConnMaxLifetime:durationValue("DB_CONN_MAX_LIFETIME",30*time.Minute),
		JWTSecret:[]byte(os.Getenv("JWT_SECRET")), JWTIssuer:value("JWT_ISSUER","movie-platform"),
		JWTTTL:durationValue("JWT_TTL",15*time.Minute),
	}
	if c.Environment!="development" && c.Environment!="test" && c.Environment!="production" {return c,errors.New("invalid APP_ENV")}
	if strings.ContainsAny(c.Port,"/: ") || c.Port=="" {return c,errors.New("invalid APP_PORT")}
	if c.DBUser=="" || c.DBName=="" || c.DBHost=="" || c.DBPort=="" {return c,errors.New("database configuration is incomplete")}
	if len(c.JWTSecret)<32 {return c,errors.New("JWT_SECRET must be at least 32 random bytes")}
	if c.JWTTTL<time.Minute || c.JWTTTL>time.Hour {return c,errors.New("JWT_TTL must be between 1m and 1h")}
	if c.DBMaxOpen<1 || c.DBMaxIdle<0 || c.DBMaxIdle>c.DBMaxOpen {return c,errors.New("invalid database pool settings")}
	if c.Environment=="production" && (c.DBPassword=="" || len(c.AllowedOrigins)==0) {return c,errors.New("production requires DB_PASSWORD and explicit CORS origins")}
	return c,nil
}
func value(k,d string)string{if v:=strings.TrimSpace(os.Getenv(k));v!=""{return v};return d}
func intValue(k string,d int)int{v,e:=strconv.Atoi(os.Getenv(k));if e!=nil{return d};return v}
func durationValue(k string,d time.Duration)time.Duration{v,e:=time.ParseDuration(os.Getenv(k));if e!=nil{return d};return v}
func csv(s string)[]string{var out []string;for _,v:=range strings.Split(s,","){if v=strings.TrimSpace(v);v!=""{out=append(out,v)}};return out}
