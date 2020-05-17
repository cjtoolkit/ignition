package cache

import "time"

func GetAndCheckExpiration(c Core, key string, expiration time.Duration) ([]byte, error) {
	if c, ok := c.(CoreGetCheck); ok {
		return c.GetBytesCheck(key, expiration)
	}
	return c.GetBytes(key)
}
