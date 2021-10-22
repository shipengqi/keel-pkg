package cmd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/shipengqi/keel-pkg/lib/utils/fsutil"
	"path/filepath"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/shipengqi/keel-pkg/lib/log"
	"github.com/spf13/cobra"
)

type pushOptions struct {
	accessKey string
	secretKey string
	bucket    string
	pkgUri    string
}

func newPushCommand() *cobra.Command {
	o := &pushOptions{}
	c := &cobra.Command{
		Use:   "push [options]",
		Short: "push package",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !fsutil.PathExists(o.pkgUri) {
				return errors.Errorf("%s not found", o.pkgUri)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return push(o)
		},
	}
	c.DisableFlagsInUseLine = true
	c.Flags().StringVarP(&o.accessKey, "access-key", "k", "", "The AccessKey")
	c.Flags().StringVarP(&o.secretKey, "secret-key", "s", "","The SecretKey")
	c.Flags().StringVarP(&o.bucket, "bucket", "b", "keel","The bucket of Storage")
	c.Flags().StringVar(&o.pkgUri, "pkg-uri", "", "The location of package")

	return c
}

func push(opts *pushOptions) error {
	bucket := opts.bucket
	key := filepath.Base(opts.pkgUri)
	log.Debugf("pushing [%s] to [%s] ...", opts.pkgUri, bucket)

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseHTTPS = true
	cfg.UseCdnDomains = false

	po := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(opts.accessKey, opts.secretKey)
	token := po.UploadToken(mac)

	uploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := uploader.PutFile(context.Background(), &ret, token, key, opts.pkgUri, &storage.PutExtra{})
	if err != nil {
		return err
	}
	log.Debugf("key [%s], hash [%s]", ret.Key,ret.Hash)

	return nil
}
