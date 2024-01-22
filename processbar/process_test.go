package process

import (
	"testing"
	"time"
)

func TestPrint(t *testing.T) {
	total := 10
	p := NewBar()
	for i := 0; i <= total; i++ {
		p.Print(i, total)
		time.Sleep(100 * time.Millisecond)
	}
}
