package etcd_dlm

import (
	"context"
	"errors"
	"log"

	"github.com/ICoder0/etcd-dlm/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type DlmConfig struct {
	Key         string // 租约绑定的key
	Ttl         int64  // 租约过期时间
	IsKeepaLive bool   // 是否开启 keepalive
	CloseFunc   func() // 租约关闭时的回调函数
}

// Lock 抢锁
func Lock(ctx context.Context, cli *clientv3.Client, cfg DlmConfig) (err error, revoke func(), succeed bool) {
	grant, err := etcd.NewLease(cli).Grant(ctx, cfg.Ttl)
	if err != nil {
		return err, nil, false
	}

	// 创建交易
	txn := etcd.NewKV(cli).Txn(ctx)
	txn.If(etcd.Compare(etcd.CreateRevision(cfg.Key), "=", 0)).
		Then(etcd.OpPut(cfg.Key, "", etcd.WithLease(grant.ID))).
		Else()
	txnResp, err := txn.Commit()
	if err != nil {
		return err, nil, false
	}

	// 判断 txn.if 是否成立
	if !txnResp.Succeeded {
		return errors.New("抢锁失败"), nil, false
	}

	// 续租
	if cfg.IsKeepaLive {
		var leaseRespChan <-chan *clientv3.LeaseKeepAliveResponse
		leaseRespChan, err = etcd.KeepAlive(ctx, cli, grant.ID)
		if err != nil {
			return err, nil, false
		}
		// 抢锁成功，设置续租
		go func() {

			defer func() {
				cfg.CloseFunc()
				_ = etcd.Revoke(ctx, cli, grant.ID)
			}()

			for {
				select {
				case leaseKeepResp := <-leaseRespChan:
					if leaseKeepResp == nil {
						log.Println("已经关闭续租,lease id :", grant.ID)
						return
					}
				}
			}
		}()
	}

	revokeFc := func() {
		_ = etcd.Revoke(ctx, cli, grant.ID)
	}
	return nil, revokeFc, true
}
