package etcd

import (
	"context"
	"fmt"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	client *etcd.Client
}

var EtcdKV *Etcd = nil

func init() {
	EtcdKV = GetEtcdClient(EtcdDefaultConfig().EtcdEndpoints,
		EtcdDefaultConfig().EtcdDialTimeout)
}
func GetEtcdClient(endpoint []string, timeout time.Duration) *Etcd {
	config := etcd.Config{
		Endpoints:   endpoint,
		DialTimeout: timeout,
	}
	cli, err := etcd.New(config)
	if err != nil {
		return nil
	}
	return &Etcd{client: cli}
}
func (e *Etcd) Get(key string) ([]byte, error) {
	resp, err := e.client.Get(context.TODO(), key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		fmt.Println("key not found")
		return nil, nil
	}
	return resp.Kvs[0].Value, nil
}
func (e *Etcd) GetPrefix(key string) ([]string, error) {
	resp, err := e.client.Get(context.TODO(), key, etcd.WithPrefix())
	if err != nil {
		return nil, err
	}
	var res []string
	for _, kv := range resp.Kvs {
		res = append(res, string(kv.Value))
	}
	return res, nil
}
func (e *Etcd) Put(key string, value []byte) error {
	_, err := e.client.Put(context.TODO(), key, string(value))
	return err
}

func (e *Etcd) Delete(key string) error {
	_, err := e.client.Delete(context.TODO(), key)
	return err
}
