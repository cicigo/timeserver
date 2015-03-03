package utils

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {

	zeus := NewCounter("zeus")
	hera := NewCounter("hera")
	athena := NewCounter("athena")
	ares := NewCounter("ares")

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		zeus.Incr(2)
		hera.Incr(1)
		athena.Incr(1)
		ares.Incr(1)
		wg.Done()
	}()
	go func() {
		hera.Incr(21)
		zeus.Incr(6)
		ares.Incr(1)
		athena.Incr(4)
		wg.Done()
	}()
	go func() {
		zeus.Incr(2)
		hera.Incr(6)
		athena.Incr(1)
		ares.Incr(1)
		wg.Done()
	}()
	go func() {
		athena.Incr(2)
		hera.Incr(1)
		zeus.Incr(3)
		ares.Incr(1)
		wg.Done()
	}()
	wg.Wait()

	expected := map[string]int{
		"zeus":   13,
		"hera":   29,
		"ares":   4,
		"athena": 8,
	}
	actual := Dump()
	for k, v := range expected {
		if v != actual[k] {
			t.Errorf("counter %s: expected %d, got %d", k, v, actual[k])
		}
	}
}
