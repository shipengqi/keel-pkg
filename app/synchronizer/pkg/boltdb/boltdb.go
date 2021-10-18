package boltdb

import (
	"encoding/binary"
	"go/types"
	"time"

	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"

	"github.com/shipengqi/keel-pkg/lib/log"
)

const (
	_defaultTimeout  = 10 * time.Second
	_defaultFileMode = 0600
)

type Boltdb struct {
	db     *bolt.DB
	bucket string // bucket name
}

func New(dbUri string) (*Boltdb, error) {
	db, err := bolt.Open(dbUri, _defaultFileMode, &bolt.Options{Timeout: _defaultTimeout})
	if err != nil {
		return nil, err
	}
	return &Boltdb{db: db}, nil
}

func (b *Boltdb) Db() *bolt.DB {
	return b.db
}

func (b *Boltdb) Close() error {
	return b.db.Close()
}

func (b *Boltdb) Bucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte(b.bucket))
}

func (b *Boltdb) CreatBucket(domain string) error {
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

func (b *Boltdb) GetUint32(key string) (uint32, error) {
	var (
		err      error
		sumBytes []byte
	)
	err = b.db.View(func(tx *bolt.Tx) error {
		sumBytes = b.Bucket(tx).Get([]byte(key))
		return nil
	})
	if err != nil {
		return 0, err
	}
	if len(sumBytes) != int(types.Uint32) { // length not equal uint32
		log.Warnf("key: %s, sum length: %d", key, len(sumBytes))
		return 0, nil
	}
	sum := binary.LittleEndian.Uint32(sumBytes)
	return sum, nil
}

func (b *Boltdb) Diff(key string, sum uint32) (bool, error) {
	old, err := b.GetUint32(key)
	if err != nil {
		return false, err
	}
	if sum != old {
		return true, nil
	}
	return false, nil
}

func (b *Boltdb) SaveUint32(key string, value uint32) error {
	buf := make([]byte, types.Uint32)
	binary.LittleEndian.PutUint32(buf, value)
	return b.db.Update(func(tx *bolt.Tx) error {
		return b.Bucket(tx).Put([]byte(key), buf)
	})
}
