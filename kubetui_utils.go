package main

import "strings"

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

func FormatData(data string) [][]string {
	output := [][]string{}
	lines := strings.Split(data, "\n")
	for r, line := range lines {
		output = append(output, []string{})
		// output[r] = append(output[r], strings.Fields(line)...)
		output[r] = append(output[r], line)

	}
	return output
}
