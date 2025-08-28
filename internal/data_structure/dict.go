package data_structure

import (
	"time"
)

type Dict struct {
	keyValStore    map[string]interface{}
	keyExpiryStore map[string]uint64
}

func CreateDict() *Dict {
	return &Dict{
		keyValStore:    make(map[string]interface{}),
		keyExpiryStore: make(map[string]uint64),
	}
}

func (d *Dict) GetKeyExpiredStore() map[string]uint64 {
	return d.keyExpiryStore
}

func (d *Dict) Set(key string, val interface{}) {
	d.keyValStore[key] = val

}

func (d *Dict) SetExpiry(key string, ttlMs int64) {
	if ttlMs > 0 {
		d.keyExpiryStore[key] = uint64(time.Now().UnixMilli()) + uint64(ttlMs)
	}

}

func (d *Dict) Get(key string) interface{} {
	val := d.keyValStore[key]
	if val != nil {
		if d.HasExpired(key) {
			d.Del(key)
			return nil
		}
	}
	return val
}

func (d *Dict) GetExpiry(key string) (uint64, bool) {
	exp, existed := d.keyExpiryStore[key]
	return exp, existed
}

func (d *Dict) HasExpired(key string) bool {
	exp, existed := d.keyExpiryStore[key]
	if !existed {
		return false
	}
	return exp <= uint64(time.Now().UnixMilli())
}

func (d *Dict) Del(key string) bool {
	if _, exist := d.keyValStore[key]; exist {
		delete(d.keyValStore, key)
		delete(d.keyExpiryStore, key)
		return true
	}
	return false
}
