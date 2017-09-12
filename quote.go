package quotes

import (
	"bytes"
	"encoding/gob"

	"github.com/pkg/errors"
)

// Quote represents a quote, inlcuding its author and an optional source. The ID is a unique key.
type Quote struct {
	ID     int    `json:id`
	Author string `json:author`
	Text   string `json:text`
	Source string `json:source`
}

func (q Quote) Serialize() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(q)
	if err != nil {
		return nil, errors.Wrapf(err, "Serialize: encoding failed for %v", q)
	}
	return b.Bytes(), nil
}

func (q *Quote) Deserialize(b []byte) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(q)
	if err != nil {
		return errors.Wrapf(err, "Deserialize: decoding failed for %s", b)
	}
	return nil
}
