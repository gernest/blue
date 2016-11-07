package blue

import (
	"encoding/json"
	"time"
)

type Options struct {
	KeyJoinFunc   func(a, b string) string
	IsTag         func(string) bool
	IsField       func(string) bool
	IsMeasurement func(string, interface{}) (string, bool)
	IsTimeStamp   func(string, interface{}) (time.Time, bool)
	Measurement   string
}

func Line(src []byte, opts Options) (*Measurement, error) {
	object := make(map[string]interface{})
	err := json.Unmarshal(src, &object)
	if err != nil {
		return nil, err
	}
	opts = getOpts(opts)
	ctx := newCtx(opts)
	err = process(ctx, "", object)
	if err != nil {
		return nil, err
	}
	m := processCollection(ctx.C, opts)
	return m, nil
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

func processCollection(c collector, opts Options) *Measurement {
	m := &Measurement{}
	m.Name = opts.Measurement
	for k, v := range c {
		if m.Name == "" {
			if msr, ok := opts.IsMeasurement(k, v); ok {
				m.Name = msr
			}
		}
		if opts.IsTag(k) {
			if m.Tags == nil {
				m.Tags = make(tags, 0)
			}
			m.Tags = append(m.Tags, &tag{key: k, value: v})
			continue
		}
		if ts, ok := opts.IsTimeStamp(k, v); ok {
			m.Timestamp = ts
		}
		if opts.IsField(k) {
			if m.Fields == nil {
				m.Fields = make(fields, 0)
			}
			m.Fields = append(m.Fields, &field{key: k, value: v})
		}
	}
	return m
}

type collector map[string]interface{}

func (c collector) set(key string, value interface{}) {
	c[key] = value
}

type Context struct {
	C       collector
	keyJoin func(a, b string) string
}

func newCtx(o Options) *Context {
	ctx := &Context{}
	if o.KeyJoinFunc != nil {
		ctx.keyJoin = o.KeyJoinFunc
	} else {
		ctx.keyJoin = joinkey
	}
	ctx.C = make(collector)
	return ctx
}

func process(ctx *Context, ns string, v interface{}) error {
	switch v.(type) {
	case bool, float64, string:
		ctx.C.set(ns, v)
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
