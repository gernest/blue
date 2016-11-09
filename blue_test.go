package blue

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestLine(t *testing.T) {
	inputDir := "fixture/input"
	outputDir := "fixture/output"
	sample :=
		[]string{
			"1-measurement_only",
			"2-measurement_with_tag",
			"2-measurement_with_tag_and_comma",
			"2-measurement_with_tag_and_space",
			"3-measurement_with_field",
			"3-measurement_with_time",
		}
	for _, v := range sample {
		i, err := ioutil.ReadFile(
			filepath.Join(inputDir, v+".json"),
		)
		if err != nil {
			t.Fatal(err)
		}
		ss, err := Line(i, Options{
			IsTag: tagFunc("host", "region"),
		})
		if err != nil {
			t.Fatal(err)
		}
		s := ss.line()
		e, err := ioutil.ReadFile(filepath.Join(outputDir, v))
		if err != nil {
			t.Fatal(err)
		}
		e = bytes.TrimSpace(e)
		if s != string(e) {
			t.Errorf("expected %s got %s", string(e), s)
		}
	}
}

func tagFunc(tg ...string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		low := strings.ToLower(key)
		for _, v := range tg {
			t := strings.ToLower(v)
			if t == low {
				return key, true
			}
		}
		return "", false
	}
}
