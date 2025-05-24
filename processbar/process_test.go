package processbar

import (
	"testing"
	"time"
)

func TestPrint(t *testing.T) {
	total := 10
	p := NewBar(WithTotal(total))
	for i := 0; i <= total; i++ {
		p.Refresh()
		time.Sleep(100 * time.Millisecond)
	}
}
