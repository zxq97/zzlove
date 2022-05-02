package rpc

import (
	"context"
	"strings"
	"sync"
	"time"
	"zzlove/internal/constant"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

type EtcdDiscover struct {
	addr       []string
	timeout    time.Duration
	keyPrefix  string
	done       chan struct{}
	etcdClient *clientv3.Client
	watchChan  clientv3.WatchChan
	cc         resolver.ClientConn
	svcAddrMap sync.Map
	lock       sync.Mutex
}

func newEtcdDiscover(etcdAddr []string, dialTimeout time.Duration, svcPrefix string) *EtcdDiscover {
	return &EtcdDiscover{
		addr:      etcdAddr,
		timeout:   dialTimeout,
		keyPrefix: svcPrefix,
	}
}

func (er *EtcdDiscover) Scheme() string {
	return constant.EtcdScheme
}

func (er *EtcdDiscover) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	er.cc = cc

	if _, err := er.start(); err != nil {
		return nil, err
	}
	return er, nil
}

func (er *EtcdDiscover) ResolveNow(o resolver.ResolveNowOptions) {
}

func (er *EtcdDiscover) Close() {
	er.done <- struct{}{}
}

func (er *EtcdDiscover) start() (chan<- struct{}, error) {
	var err error
	er.etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   er.addr,
		DialTimeout: er.timeout,
	})
	if err != nil {
		return nil, err
	}
	resolver.Register(er)

	er.done = make(chan struct{})

	err = er.sync()
	if err != nil {
		return nil, err
	}

	go er.watch()

	return er.done, nil
}

func (er *EtcdDiscover) watch() {
	t := time.NewTicker(constant.DefaultTicker)
	er.watchChan = er.etcdClient.Watch(context.Background(), er.keyPrefix, clientv3.WithPrefix())

	for {
		select {
		case <-er.done:
			return
		case res, ok := <-er.watchChan:
			if ok {
				er.update(res.Events)
			}
		case <-t.C:
			err := er.sync()
			if err != nil {
				excLogger.Println("etcd sync err", err)
			}
		}
	}
}

func (er *EtcdDiscover) updateState() {
	er.lock.Lock()
	defer er.lock.Unlock()
	svcAddr := make([]resolver.Address, 0)
	er.svcAddrMap.Range(func(k, v interface{}) bool {
		addr, ok := k.(string)
		if ok {
			svcAddr = append(svcAddr, resolver.Address{Addr: addr})
		}
		return true
	})
	er.cc.UpdateState(resolver.State{Addresses: svcAddr})
}

func (er *EtcdDiscover) update(events []*clientv3.Event) {
	for _, v := range events {
		switch v.Type {
		case mvccpb.PUT:
			er.svcAddrMap.Store(string(v.Kv.Value), struct{}{})
		case mvccpb.DELETE:
			ks := strings.Split(string(v.Kv.Key), "_")
			if len(ks) > 1 {
				er.svcAddrMap.Delete(ks[1])
			}
		}
	}
	er.updateState()

}

func (er *EtcdDiscover) sync() error {
	ctx, cancel := context.WithTimeout(context.Background(), constant.DefaultTimeout)
	defer cancel()

	res, err := er.etcdClient.Get(ctx, er.keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, v := range res.Kvs {
		er.svcAddrMap.Store(string(v.Value), struct{}{})
	}

	er.updateState()
	return nil
}
