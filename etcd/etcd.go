package etcd

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caiotedim/crawler-tags/config"
	"github.com/golang/glog"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
)

func verifySlashKey(key string) string {
	if !strings.HasPrefix(key, "/") {
		key = fmt.Sprintf("/%s", key)
	}
	return key
}

// EtcdPut function to insert data on etcd
func EtcdPut(c *config.Config, key string, data []byte) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("%s:%s", c.Db.Host, c.Db.Port)},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		glog.Error(err)
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	dataString := base64.StdEncoding.EncodeToString(data)
	_, err = cli.Put(ctx, verifySlashKey(key), dataString)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			glog.Errorf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			glog.Errorf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			glog.Errorf("client-side error: %v\n", err)
		default:
			glog.Errorf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
		return err
	}

	return nil
}

// EtcdGet function to get object on etcd
func EtcdGet(c *config.Config, key, id string) ([]byte, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("%s:%s", c.Db.Host, c.Db.Port)},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		glog.Error(fmt.Sprintf("Id:[%s] %v", id, err))
		return nil, err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	getResp, err := cli.Get(ctx, verifySlashKey(key))
	cancel()
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	var data string
	if len(getResp.Kvs) > 1 {
		glog.Error(fmt.Sprintf("Id:[%s] More than 1 response on etcd key:[%s]", id, key))
		return nil, fmt.Errorf("Id:[%s] More than 1 response on etcd key:[%s]", id, key)
	} else if len(getResp.Kvs) == 0 {
		glog.Infof("Id:[%s] Not found key:[%s]", id, key)
		return nil, nil
	}
	for _, ev := range getResp.Kvs {
		data = string(ev.Value)
	}

	obj, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		glog.Error(fmt.Sprintf("Id:[%s] %v", id, err))
	}
	return obj, err
}

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Set("logtostderr", "false")
	flag.Set("log_dir", "/tmp/crawler-tags")
	flag.Parse()
}
