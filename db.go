package quotes

import (
	"github.com/coreos/bbolt"
	"github.com/pkg/errors"
)

// DB is a quote database.
type DB struct {
	db *bolt.DB
}

const (
	quoteBucket = "standard"
)

// Open opens the database file at path and returns a DB or an error.
func Open(path string) (*DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Open: cannot open DB file "+path)
	}
	return &DB{
		db: db,
	}, nil
}

func (d *DB) Put(q *Quote) error {
	err := d.db.Update(func(tx *bolt.Tx) error {

		bucket, err := tx.CreateBucketIfNotExists([]byte(quoteBucket))
		if err != nil {
			return errors.Wrapf(err, "Put: cannot open or create bucket %s", []byte(quoteBucket))
		}

		b, err := q.Serialize()
		err = bucket.Put([]byte(q.Author), b)
		if err != nil {
			return errors.Wrapf(err, "Put: cannot put quote %d into bucket", q.ID)
		}
		return nil
	})

	if err != nil {
		return errors.Wrapf(err, "Put: cannot create record for (%s|%s|%s)", q.Author, q.Text, q.Source)
	}
	return nil
}

func (d *DB) Get(author string) (*Quote, error) {
	q := &Quote{}
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(quoteBucket))
		v := b.Get([]byte(author))
		err := q.Deserialize(v)
		if err != nil {
			return errors.Wrapf(err, "Get: cannot deserialize %s", v)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Get: DB.View() failed")
	}
	return q, nil
}
