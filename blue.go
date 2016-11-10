package blue

import (
	"encoding/json"
	"time"
)

//Options provides fine grained options for processing the json input.
type Options struct {
	KeyJoinFunc   func(a, b string) string
	IsTag         func(string) (string, bool)
	IsField       func(string) bool
	IsMeasurement func(string, interface{}) (string, bool)
	IsTimeStamp   func(string, interface{}) (time.Time, bool)
	Measurement   string
}

//Line generates *Measurement object from the src. src is expected to be a valid
//json input.
func Line(src []byte, opts Options) (*Measurement, error) {
	object := make(map[string]interface{})
	err := json.Unmarshal(src, &object)
	if err != nil {
		return nil, err
	}
	opts = getOpts(opts)
	ctx := newCtx(opts)
	err = collect(ctx, "", object)
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
		if tg, ok := opts.IsTag(k); ok {
			if m.Tags == nil {
				m.Tags = make(Tags, 0)
			}
			if tg == "" {
				tg = k
			}
			m.Tags = append(m.Tags, &Tag{Key: tg, Value: v})
			continue
		}
		if ts, ok := opts.IsTimeStamp(k, v); ok {
			m.Timestamp = ts
		}
		if opts.IsField(k) {
			if m.Fields == nil {
				m.Fields = make(Fields, 0)
			}
			m.Fields = append(m.Fields, &Field{Key: k, Value: v})
		}
	}
	return m
}

type collector map[string]interface{}

func (c collector) set(key string, value interface{}) {
	c[key] = value
}

//Context is the procesing context for json input. It is where the key, value
//pairs are collected.
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

func collect(ctx *Context, ns string, v interface{}) error {
	switch v.(type) {
	case bool, float64, string:
		ctx.C.set(ns, v)
	case []interface{}:
		for _, i := range v.([]interface{}) {
			err := collect(ctx, ns, i)
			if err != nil {
				return err
			}
		}
	case map[string]interface{}:
		for key, value := range v.(map[string]interface{}) {
			err := collect(ctx, ctx.keyJoin(ns, key), value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
