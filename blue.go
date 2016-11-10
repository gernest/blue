package blue

import (
	"encoding/json"
	"time"
)

//Options provides fine grained options for processing the json input.
type Options struct {
	//This is the function that joins the keys when flattering the json object.
	//The first argument is the top level key(although this might be not the
	//case for deeply nested objects) and the scond is the current key.
	//
	// The returned string is the key that will be used. This is implementation
	// specific, you can do whatever the hell you want with this.
	KeyJoinFunc func(a, b string) string

	// Checks  if the key is a tag. If the second returned value is true then
	// the first returned value is going to be used as key.
	//
	// Something to consider here the key might be a joint keys for the deeply
	// nested objects. The keys will be supposedly joined by the KeyJoinFunc
	// implementation.
	IsTag func(string) (string, bool)

	// Checks if the aregumen( a key) belongs to a Field. Implementations should
	// return true if the key is a field and false otherwise.
	IsField func(string) bool

	//IsMeasurement this determines if the given key is a key with the
	//measurement name.
	//
	//It is up to the implementaion to return the measurement name, true. If the
	//second returned value is false then the key is assumend to be not
	//measurement name.
	//
	// If the Measurement field is set i.e not empty then this function is nver
	// going to be used.
	IsMeasurement func(string, interface{}) (string, bool)

	// Determines the timestamp of the measurement. Timestamp is set only once.
	IsTimeStamp func(string, interface{}) (time.Time, bool)
	Measurement string
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
