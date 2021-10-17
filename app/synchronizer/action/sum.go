package action

import (
	"encoding/binary"
	"go/types"

	bolt "go.etcd.io/bbolt"

	"github.com/shipengqi/keel-pkg/app/synchronizer/pkg/boltdb"
	"github.com/shipengqi/keel-pkg/lib/log"
)

const (
	NameSum = "sum"
)

type suma struct {
	*action

	dbUri string
	db    *boltdb.Boltdb
}

func NewSumAction(dbUri string) Interface {
	db, err := boltdb.New(dbUri)
	if err != nil {
		panic(err)
	}
	a := &suma{
		action: &action{
			name: NameSum,
			close: func() error {
				return db.Close()
			},
		},
		dbUri: dbUri,
		db:    db,
	}
	return a
}

func (s *suma) Run() error {
	log.Infof("list all check sum form database: %s", s.dbUri)
	if err := s.db.Db().View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(bName []byte, b *bolt.Bucket) error {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if len(v) != int(types.Uint32) {
					log.Warnf("wrong bucket:%s key=%s", bName, k)
					continue
				}
				log.Debugf("bucket:%-35s key=%-65s, value=%v", bName, k, binary.LittleEndian.Uint32(v))
			}
			return nil
		})
	}); err != nil {
		return err
	}
	return nil
}
