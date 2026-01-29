package test

import (
	"fmt"
	"mcp_service/pkg/memcache"
	"testing"
)

func TestMemcache(t *testing.T) {
	memcache.InitMemcache()
	// memcache.SetMemcache("key1", "value1")
	value := memcache.GetMemcache("key1")
	fmt.Println(value)
	memcache.SetMemcacheWithExpiration("key2", "value2", 10)
	value, expiration := memcache.GetMemcacheWithExpiration("key2")
	fmt.Println(value, expiration)
	ok := memcache.DeleteMemcache("key1")
	fmt.Println(ok)
}
