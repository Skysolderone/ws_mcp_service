package memcache

import (
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

var Mc *memcache.Client

func InitMemcache() {
	Mc = memcache.New("54.65.152.26:11211")
}

func GetMemcache(key string) string {
	item, err := Mc.Get(key)
	if err != nil {
		return ""
	}
	return string(item.Value)
}

func GetMemcacheFloat(key string) float64 {
	item, err := Mc.Get(key)
	if err != nil {
		return 0
	}
	value, err := strconv.ParseFloat(string(item.Value), 64)
	if err != nil {
		return 0
	}
	return value
}

func SetMemcache(key string, value string) bool {
	err := Mc.Set(&memcache.Item{
		Key:   key,
		Value: []byte(value),
	})
	if err != nil {
		return false
	}
	return true
}

func SetMemcacheFloat(key string, value float64) bool {
	err := Mc.Set(&memcache.Item{
		Key:   key,
		Value: []byte(strconv.FormatFloat(value, 'f', -1, 64)),
	})
	if err != nil {
		return false
	}
	return true
}
func SetMemcacheWithExpiration(key string, value string, expiration int32) {
	Mc.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: expiration,
	})
}

func GetMemcacheWithExpiration(key string) (string, int32) {
	item, err := Mc.Get(key)
	if err != nil {
		return "", 0
	}
	return string(item.Value), item.Expiration
}

func DeleteMemcache(key string) bool {
	err := Mc.Delete(key)
	if err != nil {
		return false
	}
	return true
}
