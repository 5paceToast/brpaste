package storage

import "github.com/go-redis/redis/v7"

var _ CHR = &Redis{}

// Redis storage engine
type Redis redis.Client

// Create an entry in redis
func (r *Redis) Create(key, value string, checkcollision bool) error {
	if !r.Healthy() {
		return Unhealthy
	}
	if checkcollision {
		col, err := r.Exists(key).Result()
		if err != nil {
			return Unhealthy
		}
		if col > 0 {
			return Collision
		}
	}
	_, err := r.Set(key, value, 0).Result()
	return err
}

func (r *Redis) Read(key string) (string, error) {
	if !r.Healthy() {
		return "", Unhealthy
	}
	res, err := r.Get(key).Result()
	if err == redis.Nil {
		return res, NotFound
	}
	return res, err
}

// Healthy determines whether redis is responding to pings
func (r *Redis) Healthy() bool {
	_, err := r.Ping().Result()
	if err != nil {
		return false
	}
	return true
}
