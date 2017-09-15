package quotes

import (
	"os"
	"reflect"
	"testing"
)

func TestOpenClose(t *testing.T) {
	type args struct {
		path string
	}
	path := "testdata/db"
	tests := []struct {
		name string
		args args
	}{
		{"Test01", args{path}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Open(tt.args.path)
			if err != nil {
				t.Errorf("Open() error = %v", err)
				return
			}
			if got == nil {
				t.Error("Open(): got == nil")
			}
			if got.db == nil {
				t.Error("Open(): got.db == nil")
			}
			if got.db.Path() != path {
				t.Error("Open(): path = %s, expected = %s", got.db.Path(), path)
			}
			err = got.Close()
			if err != nil {
				t.Errorf("Cannot close database %s", path)
			}
		})
	}
}

func TestDB_PutAndGet(t *testing.T) {
	tests := []struct {
		name  string
		quote Quote
	}{
		{"Create", Quote{Author: "007", Text: "Stirred, not shaken", Source: "Graphite is ephemeral"}},
		{"Update", Quote{Author: "007", Text: "Shaken, not stirred", Source: "Diamonds Are Forever"}},
	}
	path := "testdata/putgetdb"

	// Setup
	d, err := Open(path)
	if err != nil {
		t.Fatalf("Open(): Cannot open %s", path)
	}

	// Test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = d.Put(&tt.quote); err != nil {
				t.Errorf("DB.Put() error = %v", err)
			}
			q, err := d.Get(tt.quote.Author)
			if err != nil {
				t.Errorf("DB.Get(): Cannot get record for %s", tt.quote.Author)
			}
			if !reflect.DeepEqual(*q, tt.quote) {
				t.Errorf("Expected: %#v, got: %#v", *q, tt.quote)
			}
		})
	}

	// Teardown
	err = d.Close()
	if err != nil {
		t.Errorf("Cannot close %s", path)
	}
	err = os.Remove(path)
	if err != nil {
		t.Errorf("Cannot remove %s", path)
	}
}

func TestDB_List(t *testing.T) {
	tt := struct {
		data, want []*Quote
	}{
		[]*Quote{
			&Quote{Author: "Leo Babauta", Text: "This statement is completely mis-quoted."},
			&Quote{Author: "Leo Babauta", Text: "The value of doing is so much greater than the value of being safe and doing nothing."},
			&Quote{Author: "Albert Szent-Györgyi", Text: "Discovery consists of seeing what everybody has seen and thinking what nobody has thought."},
		},
		[]*Quote{
			// The items in the DB are sorted, due to the
			// internal B+tree data structure. Hence the output of List()
			// is expected to be sorted by author.
			&Quote{Author: "Albert Szent-Györgyi", Text: "Discovery consists of seeing what everybody has seen and thinking what nobody has thought."},
			&Quote{Author: "Leo Babauta", Text: "The value of doing is so much greater than the value of being safe and doing nothing."},
		},
	}
	path := "testdata/listdb"

	// Setup
	d, err := Open(path)
	if err != nil {
		t.Fatalf("Open(): Cannot open %s", path)
	}
	// Fill the DB
	for _, q := range tt.data {
		err := d.Put(q)
		if err != nil {
			t.Fatalf("Cannot fill test database: " + err.Error())
		}
	}

	// Test
	got, err := d.List()
	if err != nil {
		t.Errorf("DB.List() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, tt.want) {
		t.Errorf("DB.List() = %v, want %v", got, tt.want)
	}

	// Teardown
	err = d.Close()
	if err != nil {
		t.Errorf("Cannot close %s", path)
	}
	err = os.Remove(path)
	if err != nil {
		t.Errorf("Cannot remove %s", path)
	}
}
