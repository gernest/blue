package blue

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"time"
)

type Options struct {
	KeyJoinFunc   func(a, b string) string
	IsTag         func(string) bool
	IsField       func(string) bool
	IsMeasurement func(string, interface{}) (string, bool)
	IsTimeStamp   func(string, interface{}) bool
	Measurement   string
}

//Line accepts JSON input and returns influxdb line compatible string.
func Line(src io.Reader, opts Options) (string, error) {
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return "", err
	}
	return line(data, opts)
}

func line(src []byte, opts Options) (string, error) {
	object := make(map[string]interface{})
	err := json.Unmarshal(src, &object)
	if err != nil {
		return "", err
	}
	opts = getOpts(opts)
	ctx := newCtx(opts)
	err = process(ctx, "", object)
	if err != nil {
		return "", err
	}
	m := processCollection(ctx.c, opts)
	return m.line(), nil
}

func getOpts(opts Options) Options {
	if opts.IsTag == nil {
		opts.IsTag = IsTag
	}
	if opts.IsField == nil {
		opts.IsField = IsField
	}
	if opts.IsTimeStamp == nil {
		opts.IsTimeStamp = IsTimeStamp
	}
	if opts.IsMeasurement == nil {
		opts.IsMeasurement = IsMeasurement
	}
	return opts
}

func processCollection(c collector, opts Options) *measurement {
	m := &measurement{}
	m.name = opts.Measurement
	for k, v := range c {
		if m.name == "" {
			if msr, ok := opts.IsMeasurement(k, v); ok {
				m.name = msr
				continue
			}
		}
		if opts.IsTag(k) {
			if m.tags == nil {
				m.tags = make(tags, 0)
			}
			m.tags = append(m.tags, &tag{key: k, value: v})
			continue
		}
		if opts.IsTimeStamp(k, v) {
			switch v.(type) {
			case float64:
				fv := v.(float64)
				m.timestamp = time.Unix(0, int64(fv))
				continue
			}
		}
		if m.fields == nil {
			m.fields = make(fields, 0)
		}
		m.fields = append(m.fields, &field{key: k, value: v})
	}
	return m
}

type collector map[string]interface{}

func (c collector) set(key string, value interface{}) {
	c[key] = value
}

type context struct {
	c       collector
	keyJoin func(a, b string) string
}

func newCtx(o Options) *context {
	ctx := &context{}
	if o.KeyJoinFunc != nil {
		ctx.keyJoin = o.KeyJoinFunc
	} else {
		ctx.keyJoin = joinkey
	}
	ctx.c = make(collector)
	return ctx
}

func process(ctx *context, ns string, v interface{}) error {
	switch v.(type) {
	case bool, float64, string:
		ctx.c.set(ns, v)
	case []interface{}:
		for _, i := range v.([]interface{}) {
			err := process(ctx, ns, i)
			if err != nil {
				return err
			}
		}
	case map[string]interface{}:
		for key, value := range v.(map[string]interface{}) {
			err := process(ctx, ctx.keyJoin(ns, key), value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
