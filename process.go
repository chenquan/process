//   Copyright 2022 chenquan
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package progress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Process struct {
	Total  float64
	value  float64
	Len    int
	Format string
	W      io.Writer
	Prefix string
	Cancel func()
	mutex  sync.Mutex

	startTime time.Time
	once      sync.Once
}

func (p *Process) SetLen(n int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Len = n
}

func (p *Process) SetPrefix(prefix string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Prefix = prefix
}

func (p *Process) SetTotal(v float64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Total = v
	p.render()
}

func (p *Process) SetValue(v float64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.setValue(v)
}

func (p *Process) Increment() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.setValue(p.value + 1)
}

func (p *Process) setValue(v float64) {
	p.value = v
	p.render()
}

func (p *Process) render() {
	p.once.Do(func() {
		p.startTime = time.Now()
	})

	total := p.getTotal()
	format := p.getFormat()
	w := p.getWriter()
	n := p.getLen()
	rate := p.getRate()

	done := int(float64(n) * rate)
	_, _ = fmt.Fprintf(
		w,
		"\r%s%5.1f%% ["+strings.Repeat("#", done)+strings.Repeat(".", n-done)+"] ["+format+"/"+format+" in %s]",
		p.Prefix, rate*100, p.value, total, time.Now().Sub(p.startTime),
	)

	if p.value >= total {
		if p.Cancel != nil {
			p.Cancel()
		}
	}
}

func (p *Process) getTotal() float64 {
	if p.Total > 0 {
		return p.Total
	}

	return 100
}

func (p *Process) getFormat() string {
	if p.Format != "" {
		return p.Format
	}

	return "%1.0f"
}

func (p *Process) getWriter() io.Writer {
	if p.W != nil {
		return p.W
	}

	return os.Stdout
}

func (p *Process) getLen() int {
	if p.Len > 0 {
		return p.Len
	}

	return 24
}

func (p *Process) getRate() float64 {
	total := p.getTotal()
	if p.value >= total {
		return 1
	}

	return p.value / total
}

func (p *Process) Done() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.setValue(p.getTotal())
}

func (p *Process) SetFormat(format string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Format = format
}

func (p *Process) SetWriter(w io.Writer) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.W = w
}
