package boltdb

import (
	"encoding/binary"
	"go/types"

	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

type Boltdb struct {
	db         *bolt.DB
	bucket     string // bucket name
}

func New(db *bolt.DB) *Boltdb {
	return &Boltdb{db: db}
}

func (b *Boltdb) Bucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte(b.bucket))
}

func (b *Boltdb) SetBucket(domain string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		var err error
		_, err = tx.CreateBucketIfNotExists([]byte(domain))
		if err != nil {
			return errors.Wrap(err, "set bucket")
		}
		b.bucket = domain
		return nil
	})
}

func (b *Boltdb) SaveUint32(key string, value uint32) error {
	buf := make([]byte, types.Uint32)
	binary.LittleEndian.PutUint32(buf, value)
	return b.db.Update(func(tx *bolt.Tx) error {
		return b.Bucket(tx).Put([]byte(key), buf)
	})
}
