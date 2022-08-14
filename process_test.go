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
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess_Done(t *testing.T) {
	buffer := &bytes.Buffer{}
	process := Process{
		W: buffer,
	}
	process.Done()
	assert.Contains(t, buffer.String(), "\r100.0%")
}

func TestProcess_SetFormat(t *testing.T) {
	buffer := &bytes.Buffer{}
	process := Process{
		W: buffer,
	}
	process.SetFormat("%0.3f")
	process.Done()
	assert.Contains(t, buffer.String(), "100.000/100.000")
}

func TestProcess_SetLen(t *testing.T) {
	buffer := &bytes.Buffer{}
	process := Process{
		W:     buffer,
		Total: 101,
	}

	process.SetFormat("%0.3f")
	process.SetLen(40)
	process.SetValue(3)
	strs := split(buffer.String())
	assert.Equal(t, "3.0%", strs[0])
	assert.Equal(t, "3.000/101.000", strs[2])
}

func TestProcess_SetWriter(t *testing.T) {
	process := Process{
		Total: 101,
	}
	buffer := &bytes.Buffer{}
	process.SetWriter(buffer)
	process.SetValue(3)
	assert.NotEmpty(t, buffer.String())
}

func TestProcess_IncValue(t *testing.T) {
	process := Process{
		Total: 101,
	}
	buffer := &bytes.Buffer{}
	process.SetWriter(buffer)
	process.IncValue()
	strs := split(buffer.String())
	assert.Contains(t, strs, "1.0%")
	assert.Contains(t, strs, "1/101")

	buffer.Reset()
	process.IncValue()
	strs = split(buffer.String())
	assert.Contains(t, strs, "2.0%")
	assert.Contains(t, strs, "2/101")
}

func TestProcess_IncTotal(t *testing.T) {
	process := Process{
		Total: 101,
	}
	buffer := &bytes.Buffer{}
	process.SetWriter(buffer)
	process.IncTotal()
	strs := split(buffer.String())
	assert.Contains(t, strs, "0.0%")
	assert.Contains(t, strs, "0/102")

	buffer.Reset()
	process.IncTotal()
	strs = split(buffer.String())
	assert.Contains(t, strs, "0.0%")
	assert.Contains(t, strs, "0/103")
}

func TestProcess_SetPrefix(t *testing.T) {
	process := Process{
		Total: 101,
	}
	buffer := &bytes.Buffer{}
	process.SetWriter(buffer)
	process.SetPrefix("foo")
	process.IncValue()
	strs := split(buffer.String())
	assert.Contains(t, strs, "1.0%")
	assert.Contains(t, strs, "foo")

	buffer.Reset()
	process.SetPrefix("bar")
	process.IncValue()
	strs = split(buffer.String())
	assert.Contains(t, strs, "2.0%")
	assert.Contains(t, strs, "bar")

}

func TestProcess_SetTotal(t *testing.T) {
	process := Process{
		Total: 101,
	}
	buffer := &bytes.Buffer{}
	process.SetWriter(buffer)
	process.IncValue()
	strs := split(buffer.String())
	assert.Contains(t, strs, "1/101")

	buffer.Reset()
	process.SetTotal(102)
	strs = split(buffer.String())
	assert.Contains(t, strs, "1/102")
}

func TestProcess_SetValue(t *testing.T) {
	process := Process{
		Total: 101,
	}
	buffer := &bytes.Buffer{}
	process.SetWriter(buffer)
	process.SetValue(99)
	strs := split(buffer.String())

	assert.Contains(t, strs, "99/101")

	buffer.Reset()
	process.SetValue(101)
	strs = split(buffer.String())
	assert.Contains(t, strs, "101/101")
	assert.Contains(t, strs, "100.0%")
}

func TestProcess_Cancel(t *testing.T) {
	i := 1
	process := Process{
		Total: 101,
		Cancel: func() {
			i++
		},
	}
	process.Done()
	assert.Equal(t, 2, i)
}

func split(s string) []string {
	//\r 3.0% [#.......................................] [3.000/101.000 in 0s]
	return strings.FieldsFunc(s, func(r rune) bool {
		return r == '\r' || r == ' ' || r == '[' || r == ']'
	})
}
