package etcd_dlm

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func TestLock(t *testing.T) {

	ctx := context.TODO()

	conf := clientv3.Config{
		Endpoints:   []string{"http://10.10.101.75:32379"},
		DialTimeout: 3 * time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}

	cli, err := clientv3.New(conf)
	if err != nil {
		t.Error(err)
	}

	type args struct {
		ctx context.Context
		cli *clientv3.Client
		cfg DlmConfig
	}
	tests := []struct {
		name        string
		args        args
		wantErr     error
		wantRevoke  func()
		wantSucceed bool
	}{
		{name: "TestLock1", args: struct {
			ctx context.Context
			cli *clientv3.Client
			cfg DlmConfig
		}{ctx: ctx, cli: cli, cfg: struct {
			Key         string
			Ttl         int64
			IsKeepaLive bool
			CloseFunc   func()
		}{Key: "test", Ttl: 10, IsKeepaLive: true, CloseFunc: func() {
			t.Log("close")
		}}}, wantErr: nil, wantRevoke: nil, wantSucceed: true},
		{name: "TestLock2", args: struct {
			ctx context.Context
			cli *clientv3.Client
			cfg DlmConfig
		}{ctx: ctx, cli: cli, cfg: struct {
			Key         string
			Ttl         int64
			IsKeepaLive bool
			CloseFunc   func()
		}{Key: "test", Ttl: 10, IsKeepaLive: true, CloseFunc: func() {
			t.Log("close")
		}}}, wantErr: nil, wantRevoke: nil, wantSucceed: true},
		{name: "TestLock3", args: struct {
			ctx context.Context
			cli *clientv3.Client
			cfg DlmConfig
		}{ctx: ctx, cli: cli, cfg: struct {
			Key         string
			Ttl         int64
			IsKeepaLive bool
			CloseFunc   func()
		}{Key: "test", Ttl: 10, IsKeepaLive: true, CloseFunc: func() {
			t.Log("close")
		}}}, wantErr: nil, wantRevoke: nil, wantSucceed: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(10)
			for i := 0; i < 10; i++ {
				go func() {
					gotErr, gotRevoke, gotSucceed := Lock(tt.args.ctx, tt.args.cli, tt.args.cfg)
					if !reflect.DeepEqual(gotErr, tt.wantErr) {
						t.Errorf("Lock() gotErr = %v, want %v", gotErr, tt.wantErr)
					}
					if gotSucceed {
						t.Log("抢锁成功")
						time.Sleep(2 * time.Second)
						gotRevoke()
						time.Sleep(2 * time.Second)
					}
					wg.Done()
				}()
			}
			wg.Wait()
		})
	}
}
