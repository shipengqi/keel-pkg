package action

import (
	"encoding/binary"
	"fmt"
	"go/types"
	"strings"

	bolt "go.etcd.io/bbolt"

	"github.com/shipengqi/keel-pkg/app/synchronizer/pkg/boltdb"
	gcrc "github.com/shipengqi/keel-pkg/app/synchronizer/pkg/registry/gcr/client"
	"github.com/shipengqi/keel-pkg/lib/log"
)

const (
	NameCheck = "check"
)

type CheckOptions struct {
	*gcrc.Options

	Db       string
	CheckSum string
}

type checka struct {
	*action

	opts *CheckOptions
	db   *boltdb.Boltdb
	gcr  *gcrc.Client
}

func NewCheckAction(opts *CheckOptions) Interface {
	db, err := boltdb.New(opts.Db)
	if err != nil {
		panic(err)
	}
	a := &checka{
		action: &action{
			name: NameCheck,
		},
		opts: opts,
		db:   db,
		gcr:  gcrc.New(opts.Options),
	}
	return a
}

func (c *checka) Close() error {
	log.Debugf("action [%s] closing ...", c.name)
	return c.db.Close()
}

func (c *checka) Run() error {
	log.Infof("check if [%s] needs to be synchronized", c.opts.CheckSum)
	if err := c.db.Db().View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(bName []byte, b *bolt.Bucket) error {
			cur := b.Cursor()
			for k, v := cur.First(); k != nil; k, v = cur.Next() {
				if len(v) != int(types.Uint32) {
					log.Warnf("wrong bucket [%s] key=%s", bName, k)
					continue
				}
				dbKey := fmt.Sprintf("%s/%s", bName, k)
				if strings.Compare(dbKey, c.opts.CheckSum) == 0 {
					lev := binary.LittleEndian.Uint32(v)
					rv, err := c.gcr.ManifestCheckSum(dbKey)
					if err != nil {
						return err
					}
					fmt.Printf("%s/%s local:%d remote:%d\n", bName, k, lev, rv)
					break
				}
			}
			return nil
		})
	}); err != nil {
		return err
	}
	return nil
}
