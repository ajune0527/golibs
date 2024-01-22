package process

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Bar struct {
	barWidth     int
	processChar  string
	arrow        bool
	pad          bool
	enableTiming bool
	start        time.Time
}

func NewBar(options ...Option) *Bar {
	b := &Bar{
		barWidth:    50,
		processChar: "-",
		arrow:       true,
		pad:         true,
	}

	for _, o := range options {
		o(b)
	}

	if b.enableTiming {
		b.start = time.Now()
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

func WithEnableTiming(enableTiming bool) Option {
	return func(b *Bar) {
		b.enableTiming = enableTiming
	}
}

func (p Bar) Print(current, total int) {
	if p.arrow {
		p.arrowProgressBar(current, total)
		return
	}

	p.noArrowProgressBar(current, total)
}

func (p Bar) Finish() {
	fmt.Printf("\n%s [bar] Total time consumption：%s\n", time.Now().Format(time.DateTime), time.Since(p.start))
}
func (p Bar) noArrowProgressBar(current, total int) {
	progress := current * p.barWidth / total
	bar := strings.Repeat(p.char(), progress) + strings.Repeat(" ", p.barWidth-progress)

	a, b := p.formatNumber(current, total)
	fmt.Printf("\r[%s] %s/%s", bar, a, b)
}

func (p Bar) arrowProgressBar(current, total int) {
	progress := current * (p.barWidth - 1) / total // 减去2个字符的空间，用于包含0和100
	bar := strings.Repeat(p.char(), progress) + ">" + strings.Repeat(" ", p.barWidth-progress-1)

	a, b := p.formatNumber(current, total)
	fmt.Printf("\r[%s] %s/%s", bar, a, b)
}

func (p Bar) char() string {
	if p.processChar == "" {
		return "-"
	}

	return p.processChar
}

func (p Bar) formatNumber(current, total int) (string, string) {
	if p.pad {
		return padNumber(current, total)
	}

	return strconv.Itoa(current), strconv.Itoa(total)
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
