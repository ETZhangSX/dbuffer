package dbuffer

import (
	"testing"
)

type testData struct {
	v string
}

func (i *testData) Load(d *map[string]string) (bool, error) {
	(*d)["test"] = i.v
	return true, nil
}

func TestNew(t *testing.T) {
	a := &testData{v: "testValue"}
	allo := func() map[string]string {
		obj := map[string]string{}
		return obj
	}
	buf := New[map[string]string](a, allo, WithInterval(Off))
	d := buf.Data()
	if a.v != d["test"] {
		t.Fatalf("want map[test:%s] but got %v", a.v, d)
	}
}
