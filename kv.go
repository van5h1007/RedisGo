package main

import "sync"

type KV struct {
	mu sync.RWMutex
	data map[string]string
}

func NewKV() *KV {
	return &KV{
		data: make(map[string]string),
	}
}

func (kv *KV) Set(key, val string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key]= val
}

func (kv *KV) Get(key string) (string, bool){
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	val, ok := kv.data[key]
	return val, ok 
}

func(kv *KV) Del(key string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	_, existed := kv.data[key]
	delete(kv.data, key)
	return existed
}

func (kv *KV) Exists(key string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	_, ok:= kv.data[key]

	return ok
}