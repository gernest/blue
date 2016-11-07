package blue

import (
	"bytes"
	"fmt"
	"sort"
	"time"
)

type Measurement struct {
	name      string
	Tags      tags
	Fields    fields
	Timestamp time.Time
}

func (m *Measurement) String() string {
	return m.line()
}

func (m *Measurement) line() string {
	var buf bytes.Buffer
	buf.WriteString(escape(m.name))
	if m.Tags != nil {
		buf.WriteRune(',')
		buf.WriteString(m.Tags.line())
	}
	if m.Fields != nil {
		buf.WriteRune(' ')
		buf.WriteString(m.Fields.line())
	}
	if !m.Timestamp.IsZero() {
		buf.WriteRune(' ')
		buf.WriteString(fmt.Sprint(m.Timestamp.UnixNano()))
	}
	return buf.String()
}

type unit struct {
	key   string
	value interface{}
}

func newUnit(key string, value interface{}) *unit {
	return &unit{key: key, value: value}
}

func (u *unit) escape() *unit {
	n := &unit{key: u.key, value: u.value}
	n.key = escape(n.key)
	if k, ok := n.value.(string); ok {
		n.value = escape(k)
	}
	return n
}

func (u unit) line() string {
	switch u.value.(type) {
	case int, int64, int32:
		u.value = fmt.Sprintf("%vi", u.value)
	}
	return u.key + "=" + fmt.Sprint(u.value)
}

// escapes the string to suit the influxdb line protocol. It returns the string
// with space,comma and equal sign escaped by the \ character
func escape(src string) string {
	var buf bytes.Buffer
	for _, v := range src {
		switch v {
		case ' ', ',', '=', '"':
			buf.WriteString(`\`)
			buf.WriteRune(v)
		default:
			buf.WriteRune(v)
		}
	}
	return buf.String()
}

type field unit

func (f *field) line() string {
	f.key = escape(f.key)
	switch f.value.(type) {
	case int, int64, int32:
		f.value = fmt.Sprintf("%vi", f.value)
	case string:
		return f.key + "=" + fmt.Sprintf("\"%s\"", f.value)
	}
	return f.key + "=" + fmt.Sprint(f.value)
}

type fields []*field

func (f fields) Len() int {
	return len(f)
}
func (f fields) Less(i, j int) bool {
	return f[i].key < f[j].key
}

func (f fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f fields) line() string {
	sort.Sort(f)
	var buf bytes.Buffer
	for _, v := range f {
		if buf.Len() == 0 {
			buf.WriteString(v.line())
			continue
		}
		buf.WriteRune(',')
		buf.WriteString(v.line())
	}
	return buf.String()
}

type tag unit

func (t tag) line() string {
	u := unit(t)
	return u.escape().line()
}

type tags []*tag

func (t tags) Len() int {
	return len(t)
}
func (t tags) Less(i, j int) bool {
	return t[i].key < t[j].key
}

func (t tags) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t tags) line() string {
	sort.Sort(t)
	var buf bytes.Buffer
	for _, v := range t {
		if buf.Len() == 0 {
			buf.WriteString(v.line())
			continue
		}
		buf.WriteRune(',')
		buf.WriteString(v.line())
	}
	return buf.String()
}
