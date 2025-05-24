package processbar

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Bar struct {
	barWidth    int
	processChar string
	arrow       bool
	pad         bool
	start       time.Time
	current     *int64
	total       int64
	step        int64
	onceStart   sync.Once
	timeElapsed time.Duration
	mtx         sync.RWMutex
}

func NewBar(options ...Option) *Bar {
	current := int64(0)
	b := &Bar{
		barWidth:    50,
		processChar: "-",
		arrow:       true,
		pad:         true,
		current:     &current,
		step:        1,
		onceStart:   sync.Once{},
		mtx:         sync.RWMutex{},
	}

	for _, o := range options {
		o(b)
	}

	return b
}

type Option func(b *Bar)

func WithProcessChar(ch string) Option {
	return func(b *Bar) {
		b.processChar = ch
	}
}

func WithBarWidth(num int) Option {
	return func(b *Bar) {
		b.barWidth = num
	}
}

func WithArrow(arrow bool) Option {
	return func(b *Bar) {
		b.arrow = arrow
	}
}

func WithPad(pad bool) Option {
	return func(b *Bar) {
		b.pad = pad
	}
}

func WithTotal(total int) Option {
	return func(b *Bar) {
		b.total = int64(total)
	}
}

func WithStep(step int) Option {
	return func(b *Bar) {
		b.total = int64(step)
	}
}

func (p *Bar) Print(current, total int) {
	p.onceStart.Do(func() {
		p.start = time.Now()
	})
	p.timeElapsed = time.Since(p.start)

	if p.arrow {
		p.arrowProgressBar(current, total)
		return
	}

	p.noArrowProgressBar(current, total)
}

func (p *Bar) Finish() {
	fmt.Printf("\n%s [bar] Total time consumption：%s\n", time.Now().Format(time.DateTime), PrettyTime(p.TimeElapsed()))
}

func (p *Bar) noArrowProgressBar(current, total int) {
	progress := current * p.barWidth / total
	if p.barWidth-progress < 0 {
		return
	}
	bar := strings.Repeat(p.char(), progress) + strings.Repeat(" ", p.barWidth-progress)

	a, b := p.formatNumber(current, total)

	fmt.Printf("\r[%s] %s/%s", bar, a, b)
}

func (p *Bar) arrowProgressBar(current, total int) {
	progress := current * (p.barWidth - 1) / total // 减去2个字符的空间，用于包含0和100
	if p.barWidth-progress-1 < 0 {
		return
	}
	bar := strings.Repeat(p.char(), progress) + ">" + strings.Repeat(" ", p.barWidth-progress-1)

	a, b := p.formatNumber(current, total)

	fmt.Printf("\r[%s] %s/%s", bar, a, b)
}

func (p *Bar) char() string {
	if p.processChar == "" {
		return "-"
	}

	return p.processChar
}

func (p *Bar) formatNumber(current, total int) (string, string) {
	if p.pad {
		return padNumber(current, total)
	}

	return strconv.Itoa(current), strconv.Itoa(total)
}

func (p *Bar) Add(incr int64) {
	atomic.AddInt64(p.current, incr)
}

func (p *Bar) Refresh() {
	p.onceStart.Do(func() {
		p.start = time.Now()
	})
	p.timeElapsed = time.Since(p.start)

	// 必须设置total
	atomic.AddInt64(p.current, p.step)
	current := atomic.LoadInt64(p.current)

	if p.arrow {
		p.arrowProgressBar(int(current), int(p.total))
		return
	}

	p.noArrowProgressBar(int(current), int(p.total))
}

func padNumber(num1, num2 int) (string, string) {
	str1 := strconv.Itoa(num1)
	str2 := strconv.Itoa(num2)

	maxLen := len(str1)
	if len(str2) > maxLen {
		maxLen = len(str2)
	}

	str1 = padString(str1, maxLen)
	str2 = padString(str2, maxLen)

	return str1, str2
}

func padString(str string, length int) string {
	for len(str) < length {
		str = " " + str
	}
	return str
}

func (p *Bar) TimeElapsed() time.Duration {
	p.mtx.RLock()
	defer p.mtx.RUnlock()
	return p.timeElapsed
}

func (p *Bar) TimeElapsedString() string {
	return PrettyTime(p.TimeElapsed()).String()
}

func PrettyTime(t time.Duration) time.Duration {
	if t == 0 {
		return 0
	}
	return t - (t % time.Second)
}
