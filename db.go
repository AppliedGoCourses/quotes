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

func (d *DB) Close() error {
	err := d.db.Close()
	if err != nil {
		return errors.Wrap(err, "Close: cannot close database")
	}
	return nil
}

// Put takes a quote and saves it to the database, using the author name
// as the key.
func (d *DB) Put(q *Quote) error {
	err := d.db.Update(func(tx *bolt.Tx) error {

		bucket, err := tx.CreateBucketIfNotExists([]byte(quoteBucket))
		if err != nil {
			return errors.Wrapf(err, "Put: cannot open or create bucket %s", []byte(quoteBucket))
		}

		b, err := q.Serialize()
		err = bucket.Put([]byte(q.Author), b)
		if err != nil {
			return errors.Wrapf(err, "Put: cannot put quote from %s into bucket", q.Author)
		}
		return nil
	})

	if err != nil {
		return errors.Wrapf(err, "Put: cannot create record for (%s|%s|%s)", q.Author, q.Text, q.Source)
	}
	return nil
}

// Get takes an author name and retrieves the corresponding quote from the DB.
func (d *DB) Get(author string) (*Quote, error) {
	q := &Quote{}
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(quoteBucket))
		if b == nil {
			return errors.Errorf("Cannot get %s - there is no bucket named %s", author, quoteBucket)
		}
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

// List lists all records in the DB.
func (d *DB) List() ([]*Quote, error) {
	// The database returns byte slices that we need to de-serialize
	// into Quote structures.
	structList := []*Quote{}

	// We use a view as we don't update anything.
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(quoteBucket))
		if b == nil {
			// It is valid to attempt listing an empty database.
			// Hence no error is returned in this case.
			return nil
		}
		// ForEach iterates over all elements of a bucket and
		// executes the passed-in function for each element.
		err := b.ForEach(func(k []byte, v []byte) error {
			q := &Quote{}
			err := q.Deserialize(v)
			if err != nil {
				return errors.Wrapf(err, "List: cannot deserialize data of author %s", k)
			}
			// We're inside a closure, so we can access structList
			// that resides in the outer scope.
			structList = append(structList, q)
			return nil
		})
		if err != nil {
			return errors.Wrapf(err, "List: failed iterating over bucket %s", quoteBucket)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "List: view failed")
	}
	return structList, nil
}
