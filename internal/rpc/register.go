package rpc

import (
	"context"
	"time"
	"zzlove/internal/constant"

	"go.etcd.io/etcd/client/v3"
)

type EtcdRegister struct {
	addr       []string
	timeout    time.Duration
	svcKey     string
	svcValue   string
	done       chan struct{}
	etcdClient *clientv3.Client
	leaseID    clientv3.LeaseID
	keepAlice  <-chan *clientv3.LeaseKeepAliveResponse
}

func newEtcdRegister(etcdAddr []string, dailTimeout time.Duration, key, value string) *EtcdRegister {
	return &EtcdRegister{
		addr:     etcdAddr,
		timeout:  dailTimeout,
		svcKey:   key,
		svcValue: value,
	}
}

func (er *EtcdRegister) Register() (chan<- struct{}, error) {
	var err error
	er.etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   er.addr,
		DialTimeout: er.timeout,
	})
	if err != nil {
		return nil, err
	}

	err = er.add()
	if err != nil {
		return nil, err
	}

	er.done = make(chan struct{})

	go er.keepAlive()

	return er.done, nil
}

func (er *EtcdRegister) Stop() {
	er.done <- struct{}{}
}

func (er *EtcdRegister) add() error {
	ctx, cancel := context.WithTimeout(context.Background(), constant.DefaultTimeout)
	defer cancel()

	res, err := er.etcdClient.Grant(ctx, constant.EtcdLeaseTTL)
	if err != nil {
		return err
	}
	er.leaseID = res.ID

	er.keepAlice, err = er.etcdClient.KeepAlive(context.Background(), res.ID)
	if err != nil {
		return err
	}

	_, err = er.etcdClient.Put(context.Background(), er.svcKey, er.svcValue, clientv3.WithLease(res.ID))
	return err
}

func (er *EtcdRegister) keepAlive() {
	var err error
	t := time.NewTicker(constant.DefaultTicker)
	for {
		select {
		case <-er.done:
			_, err = er.etcdClient.Delete(context.Background(), er.svcKey)
			if err != nil {
				excLogger.Println(err)
			}
			_, err = er.etcdClient.Revoke(context.Background(), er.leaseID)
			if err != nil {
				excLogger.Println(err)
			}
		case res := <-er.keepAlice:
			if res == nil {
				err = er.add()
				if err != nil {
					excLogger.Println(err)
				}
			}
		case <-t.C:
			err = er.add()
			if err != nil {
				excLogger.Println(err)
			}
		}
	}
}
