package main

type Kv struct {
	Key   string
	Value string
}

func NewKv(key, value string) *Kv {
	return &Kv{
		Key:   key,
		Value: value,
	}
}
