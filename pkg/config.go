package logiapkg

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

var (
	// Host --> Host for run application without protocol
	Host string

	// HostFull --> Host with protocol
	HostFull string

	// DevMode --> Dev mode for use .env or kubernetes configmap
	DevMode bool

	// RPCDialTimeout --> gRPC dial timout to another services
	RPCDialTimeout time.Duration

	// LogiaValidate --> Validation configuration
	LogiaValidate *validator.Validate

	// RedisPool --> Redis pool for open connection
	RedisPool *redis.Pool
)

func InitHost() {
	protocol := "http"
	ssl, _ := strconv.ParseBool(os.Getenv("USE_SSL"))
	if ssl == true {
		protocol = "https"
	}

	Host = os.Getenv("DOMAIN")
	port := os.Getenv("PORT")

	HostFull = protocol + "://" + Host
	if ssl == false {
		HostFull += ":" + port
	}

	Host += ":" + port
}

func InitDevMode() {
	if DevMode {
		fmt.Println("Running in development mode..")
		err := godotenv.Load()
		if err != nil {
			panic(err.Error())
		}
	}
}

func InitRedisPool() {
	RedisPool = &redis.Pool{
		MaxIdle:     100,
		MaxActive:   500,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")))
			if err != nil {
				return nil, err
			}

			if os.Getenv("REDIS_PASSWORD") != "" {
				if _, err = c.Do("AUTH", os.Getenv("REDIS_PASSWORD")); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
	}
}
