package rpc

import (
	"context"
	"strings"
	"sync"
	"time"
	"zzlove/global"
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

func (ed *EtcdDiscover) Scheme() string {
	return constant.EtcdScheme
}

func (ed *EtcdDiscover) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ed.cc = cc

	if _, err := ed.start(); err != nil {
		return nil, err
	}
	return ed, nil
}

func (ed *EtcdDiscover) ResolveNow(o resolver.ResolveNowOptions) {
}

func (ed *EtcdDiscover) Close() {
	ed.done <- struct{}{}
}

func (ed *EtcdDiscover) start() (chan<- struct{}, error) {
	var err error
	ed.etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   ed.addr,
		DialTimeout: ed.timeout,
	})
	if err != nil {
		return nil, err
	}
	resolver.Register(ed)

	ed.done = make(chan struct{})

	err = ed.sync()
	if err != nil {
		return nil, err
	}

	go ed.watch()

	return ed.done, nil
}

func (ed *EtcdDiscover) watch() {
	t := time.NewTicker(constant.DefaultTicker)
	ed.watchChan = ed.etcdClient.Watch(context.Background(), ed.keyPrefix, clientv3.WithPrefix())

	for {
		select {
		case <-ed.done:
			return
		case res, ok := <-ed.watchChan:
			if ok {
				ed.update(res.Events)
			}
		case <-t.C:
			err := ed.sync()
			if err != nil {
				global.ExcLogger.Println("etcd sync err", err)
			}
		}
	}
}

func (ed *EtcdDiscover) updateState() {
	ed.lock.Lock()
	defer ed.lock.Unlock()
	svcAddr := make([]resolver.Address, 0)
	ed.svcAddrMap.Range(func(k, v interface{}) bool {
		addr, ok := k.(string)
		if ok {
			svcAddr = append(svcAddr, resolver.Address{Addr: addr})
		}
		return true
	})
	ed.cc.UpdateState(resolver.State{Addresses: svcAddr})
}

func (ed *EtcdDiscover) update(events []*clientv3.Event) {
	for _, v := range events {
		switch v.Type {
		case mvccpb.PUT:
			ed.svcAddrMap.Store(string(v.Kv.Value), struct{}{})
		case mvccpb.DELETE:
			ks := strings.Split(string(v.Kv.Key), "_")
			if len(ks) > 1 {
				ed.svcAddrMap.Delete(ks[1])
			}
		}
	}
	ed.updateState()
}

func (ed *EtcdDiscover) sync() error {
	ctx, cancel := context.WithTimeout(context.Background(), constant.DefaultTimeout)
	defer cancel()

	res, err := ed.etcdClient.Get(ctx, ed.keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, v := range res.Kvs {
		ed.svcAddrMap.Store(string(v.Value), struct{}{})
	}

	ed.updateState()
	return nil
}
