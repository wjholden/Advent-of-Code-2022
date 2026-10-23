package main

import "testing"

func TestJetPeek(t *testing.T) {
	jg := JetGenerator{}

	for range 1000000 {
		d1, k := jg.peek()
		d2 := jg.next()
		if d1 != d2 {
			t.Errorf("Peek returned %d at %d but next returned %d", d1, k, d2)
		}
	}
}
