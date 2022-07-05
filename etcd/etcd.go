package etcd

import (
	"context"

	"go.etcd.io/etcd/client/v3"
)

// Revoke 租约回收
func Revoke(ctx context.Context, cli *clientv3.Client, leaseID clientv3.LeaseID) error {
	_, err := cli.Revoke(ctx, leaseID)
	if err != nil {
		return err
	}
	return nil
}

// KeepAlive 续租
func KeepAlive(ctx context.Context, cli *clientv3.Client, leaseID clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	return cli.KeepAlive(ctx, leaseID)
}

// NewLease 创建租约
func NewLease(cli *clientv3.Client) clientv3.Lease {
	return clientv3.NewLease(cli)
}

// NewKV 创建 kv
func NewKV(cli *clientv3.Client) clientv3.KV {
	return clientv3.NewKV(cli)
}

// Compare 创建比较器
func Compare(cmp clientv3.Cmp, result string, v interface{}) clientv3.Cmp {
	return clientv3.Compare(cmp, result, v)
}

// OpPut 创建 put 操作
func OpPut(key string, value string, opts ...clientv3.OpOption) clientv3.Op {
	return clientv3.OpPut(key, value, opts...)
}

// CreateRevision 创建比较器
func CreateRevision(key string) clientv3.Cmp {
	return clientv3.CreateRevision(key)
}

// WithLease 创建租约
func WithLease(lease clientv3.LeaseID) clientv3.OpOption {
	return clientv3.WithLease(lease)
}
