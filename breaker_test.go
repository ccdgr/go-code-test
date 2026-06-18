package main

import (
	"testing"
	"time"
)

func TestRush(t *testing.T) {
	initBreaker()
	for range 15 {
		resp := rush()
		t.Log(resp)
		time.Sleep(100 * time.Millisecond)
	}
}

func TestRushWithoutBreaker(t *testing.T) {
	initBreaker()
	for range 15 {
		resp := rushWithoutBreaker()
		t.Log(resp)
		time.Sleep(100 * time.Millisecond)
	}
}
