package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Opts struct {
	Server   string
	Port     string
	Password string
}

type OptFunc func(*Opts)

func defaultOpts() *Opts {
	return &Opts{
		Server:   "redis",
		Port:     "6379",
		Password: "password",
	}
}

type RedisDatabase struct {
	Opts     *Opts
	Instance *redis.Client
}

// Create a Redis database instance
func New(
	optFuncs ...OptFunc,
) *RedisDatabase {

	// Instantiate options with default values
	opts := defaultOpts()

	// Apply external options
	for _, f := range optFuncs {
		f(opts)
	}

	return &RedisDatabase{
		Opts: opts,
	}
}

// Configure Redis server
func WithServer(server string) OptFunc {
	return func(opts *Opts) {
		opts.Server = server
	}
}

// Configure MySQL port
func WithPort(port string) OptFunc {
	return func(opts *Opts) {
		opts.Port = port
	}
}

// Configure MySQL password
func WithPassword(password string) OptFunc {
	return func(opts *Opts) {
		opts.Password = password
	}
}

// Creates Redis database connection
func (r *RedisDatabase) CreateDatabaseConnection() {

	// Connect to Redis
	db := redis.NewClient(
		&redis.Options{
			Addr:     r.Opts.Server + ":" + r.Opts.Port,
			Password: r.Opts.Password,
			DB:       0,
		})

	// Ping the Redis server to ensure the connection is established
	pong, err := db.Ping().Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}

	fmt.Println("Connected to Redis:", pong)
	r.Instance = db
}
