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
	if Mc == nil {
		return ""
	}
	item, err := Mc.Get(key)
	if err != nil {
		return ""
	}
	return string(item.Value)
}

func GetMemcacheFloat(key string) float64 {
	if Mc == nil {
		return 0
	}
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
	if Mc == nil {
		return false
	}
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
	if Mc == nil {
		return false
	}
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
	if Mc == nil {
		return
	}
	Mc.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: expiration,
	})
}

func GetMemcacheWithExpiration(key string) (string, int32) {
	if Mc == nil {
		return "", 0
	}
	item, err := Mc.Get(key)
	if err != nil {
		return "", 0
	}
	return string(item.Value), item.Expiration
}

func DeleteMemcache(key string) bool {
	if Mc == nil {
		return false
	}
	err := Mc.Delete(key)
	if err != nil {
		return false
	}
	return true
}
