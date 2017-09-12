package quotes

import (
	"os"
	"reflect"
	"testing"

	"github.com/coreos/bbolt"
)

func TestOpen(t *testing.T) {
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
			got.db.Close()
		})
	}
}

func TestDB_Put(t *testing.T) {
	tests := []struct {
		name  string
		quote Quote
	}{
		{"Create", Quote{ID: 7, Author: "007", Text: "Shaken, not stirred", Source: "Diamonds Are Forever"}},
		{"Update", Quote{ID: 7, Author: "007", Text: "Shaken, not stirred", Source: "Diamonds Are Forever"}},
	}
	path := "testdata/db"

	// Setup
	d, err := Open(path)
	if err != nil {
		t.Errorf("Open(): Cannot open %s", path)
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
	err = d.db.Close()
	if err != nil {
		t.Fatalf("Cannot close %s", path)
	}
	err = os.Remove(path)
	if err != nil {
		t.Fatalf("Cannot remove %s", path)
	}
}

func TestDB_Get(t *testing.T) {
	type fields struct {
		db *bolt.DB
	}
	type args struct {
		author string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Quote
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				db: tt.fields.db,
			}
			got, err := d.Get(tt.args.author)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
