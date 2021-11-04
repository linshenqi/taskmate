package cluster

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/linshenqi/sptty"
	"github.com/linshenqi/taskmate/src/base"

	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Service struct {
	sptty.BaseService

	cfg      Config
	clientID atomic.Value
	isMaster atomic.Value

	etcd    *v3.Client
	session *concurrency.Session
	lease   *v3.LeaseGrantResponse
}

func (s *Service) Init(app sptty.ISptty) error {
	if err := app.GetConfig(s.ServiceName(), &s.cfg); err != nil {
		return nil
	}

	s.clientID.Store(sptty.GenerateUID())
	s.isMaster.Store(false)

	if !s.cfg.Enable {
		sptty.Log(sptty.InfoLevel, "Cluster Mode Disabled", s.ServiceName())
		return nil
	}

	if err := s.initCluster(); err != nil {
		return nil
	}

	return nil
}

func (s *Service) ServiceName() string {
	return base.ServiceCluster
}

func (s *Service) initCluster() error {

	var err error

	s.etcd, err = v3.New(v3.Config{
		Endpoints:   s.cfg.EtcdURLs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}

	s.lease, err = s.etcd.Grant(context.TODO(), base.LeaseTTL)
	if err != nil {
		return err
	}

	s.session, err = concurrency.NewSession(s.etcd, concurrency.WithLease(s.lease.ID))
	if err != nil {
		return err
	}

	go s.asyncWatch()
	go s.asyncClientAlive()

	s.tryBecomeMaster()

	return nil
}

func (s *Service) asyncClientAlive() {

	ctx, cancel := context.WithTimeout(context.Background(), base.LeaseTTL*time.Second)
	defer cancel()

	_, err := s.etcd.Put(ctx, fmt.Sprintf("%s/%s", base.KeyClients, s.clientID.Load().(string)), s.clientID.Load().(string), v3.WithLease(s.lease.ID))
	if err != nil {
		sptty.Log(sptty.ErrorLevel, fmt.Sprintf("%s.Put Failed: %s", base.CurrentFuncName(), err.Error()), s.ServiceName())
		return
	}

	for {
		_, err := s.etcd.KeepAliveOnce(context.Background(), s.lease.ID)
		if err != nil {
			sptty.Log(sptty.ErrorLevel, fmt.Sprintf("%s.KeepAliveOnce Failed: %s", base.CurrentFuncName(), err.Error()), s.ServiceName())
		}

		time.Sleep(1 * time.Second)
	}
}

func (s *Service) asyncWatch() {
	for {
		rch := s.etcd.Watch(context.Background(), base.KeyClients, v3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				sptty.Log(sptty.DebugLevel, fmt.Sprintf("%s: %s %s %s", base.CurrentFuncName(), ev.Type, string(ev.Kv.Key), string(ev.Kv.Value)), s.ServiceName())

				switch ev.Type {
				case v3.EventTypeDelete:
					// client offline
					s.tryBecomeMaster()
				}
			}
		}
	}
}

func (s *Service) tryBecomeMaster() {
	if s.isMaster.Load().(bool) {
		return
	}

	ele := concurrency.NewElection(s.session, base.KeyMaster)
	ctx, cancel := context.WithTimeout(context.Background(), base.LeaseTTL*time.Second)
	defer cancel()

	if err := ele.Campaign(ctx, s.clientID.Load().(string)); err != nil {
		s.isMaster.Store(false)
		sptty.Log(sptty.InfoLevel, "Master Election Failed", s.ServiceName())
		return
	}

	s.isMaster.Store(true)
	sptty.Log(sptty.InfoLevel, fmt.Sprintf("I Am Master: %s", s.clientID.Load().(string)), s.ServiceName())
}
