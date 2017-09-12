package quotes

import (
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

func TestDB_Create(t *testing.T) {
	tests := []struct {
		name  string
		quote Quote
	}{
		{"01", Quote{ID: 7, Author: "007", Text: "Shaken, not stirred", Source: "Diamonds Are Forever"}},
	}
	path := "testdata/db"
	d, err := Open(path)
	if err != nil {
		t.Errorf("Open(): Cannot open %s", path)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = d.Create(&tt.quote); err != nil {
				t.Errorf("DB.Create() error = %v", err)
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
	d.db.Close()
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
