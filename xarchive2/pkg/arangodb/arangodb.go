package arangodb

import (
	"context"
	"encoding/json"

	driver "github.com/arangodb/go-driver"
	"github.com/golang/glog"
	"github.com/jalapeno/topology/pkg/dbclient"
	notifier "github.com/jalapeno/topology/pkg/kafkanotifier"
	"github.com/sbezverk/gobmp/pkg/bmp"
	"github.com/sbezverk/gobmp/pkg/tools"
)

type arangoDB struct {
	dbclient.DB
	*ArangoConn
	stop          chan struct{}
	l3vpnPrefixV4 driver.Collection
	l3vpnPrefixV6 driver.Collection
	vrfV4         driver.Collection
	vrfV6         driver.Collection
	srv6LocalSids driver.Collection
}

// NewDBSrvClient returns an instance of a DB server client process
func NewDBSrvClient(arangoSrv, user, pass, dbname, l3vpnPrefixV4 string, l3vpnPrefixV6 string, vrfV4 string, vrfV6 string, srv6LocalSids string) (dbclient.Srv, error) {
	if err := tools.URLAddrValidation(arangoSrv); err != nil {
		return nil, err
	}
	arangoConn, err := NewArango(ArangoConfig{
		URL:      arangoSrv,
		User:     user,
		Password: pass,
		Database: dbname,
	})
	if err != nil {
		return nil, err
	}
	arango := &arangoDB{
		stop: make(chan struct{}),
	}
	arango.DB = arango
	arango.ArangoConn = arangoConn

	// Check if peer collection exists, if not fail as Jalapeno topology is not running
	arango.l3vpnPrefixV4, err = arango.db.Collection(context.TODO(), l3vpnPrefixV4)
	if err != nil {
		return nil, err
	}
	// // Check if unicast_prefix_v4 collection exists, if not fail as Jalapeno topology is not running
	// arango.unicastprefixV4, err = arango.db.Collection(context.TODO(), unicastprefixV4)
	// if err != nil {
	// 	return nil, err
	// }
	// // Check if unicast_prefix_v4 collection exists, if not fail as Jalapeno ipv4_topology is not running
	// arango.unicastprefixV6, err = arango.db.Collection(context.TODO(), unicastprefixV6)
	// if err != nil {
	// 	return nil, err
	// }

	// check for vrf_v4 collection
	found, err := arango.db.CollectionExists(context.TODO(), vrfV4)
	if err != nil {
		return nil, err
	}
	if found {
		c, err := arango.db.Collection(context.TODO(), vrfV4)
		if err != nil {
			return nil, err
		}
		if err := c.Remove(context.TODO()); err != nil {
			return nil, err
		}
	}

	// check for vrf_v6 collection
	found, err = arango.db.CollectionExists(context.TODO(), vrfV6)
	if err != nil {
		return nil, err
	}
	if found {
		c, err := arango.db.Collection(context.TODO(), vrfV6)
		if err != nil {
			return nil, err
		}
		if err := c.Remove(context.TODO()); err != nil {
			return nil, err
		}
	}

	// check for srv6_localsids prefix collection
	found, err = arango.db.CollectionExists(context.TODO(), srv6LocalSids)
	if err != nil {
		return nil, err
	}
	if found {
		c, err := arango.db.Collection(context.TODO(), srv6LocalSids)
		if err != nil {
			return nil, err
		}
		if err := c.Remove(context.TODO()); err != nil {
			return nil, err
		}
	}

	// create vrf_v4 collection
	var vrfV4_options = &driver.CreateCollectionOptions{ /* ... */ }
	arango.vrfV4, err = arango.db.CreateCollection(context.TODO(), "vrf_v4", vrfV4_options)
	if err != nil {
		return nil, err
	}
	// check if collection exists, if not fail as processor has failed to create collection
	arango.vrfV4, err = arango.db.Collection(context.TODO(), vrfV4)
	if err != nil {
		return nil, err
	}

	// create vrf_v6 collection
	var vrfV6_options = &driver.CreateCollectionOptions{ /* ... */ }
	arango.vrfV6, err = arango.db.CreateCollection(context.TODO(), "vrf_v6", vrfV6_options)
	if err != nil {
		return nil, err
	}
	// check if collection exists, if not fail as processor has failed to create collection
	arango.vrfV6, err = arango.db.Collection(context.TODO(), vrfV6)
	if err != nil {
		return nil, err
	}

	// create srv6_localsids collection
	var srv6LocalSids_options = &driver.CreateCollectionOptions{ /* ... */ }
	arango.srv6LocalSids, err = arango.db.CreateCollection(context.TODO(), "srv6_localsids", srv6LocalSids_options)
	if err != nil {
		return nil, err
	}
	// check if collection exists, if not fail as processor has failed to create collection
	arango.srv6LocalSids, err = arango.db.Collection(context.TODO(), srv6LocalSids)
	if err != nil {
		return nil, err
	}
	return arango, nil
}

func (a *arangoDB) Start() error {
	if err := a.loadCollection(); err != nil {
		return err
	}
	glog.Infof("Connected to arango database, starting monitor")
	go a.monitor()

	return nil
}

func (a *arangoDB) Stop() error {
	close(a.stop)

	return nil
}

func (a *arangoDB) GetInterface() dbclient.DB {
	return a.DB
}

func (a *arangoDB) GetArangoDBInterface() *ArangoConn {
	return a.ArangoConn
}

func (a *arangoDB) StoreMessage(msgType dbclient.CollectionType, msg []byte) error {
	event := &notifier.EventMessage{}
	if err := json.Unmarshal(msg, event); err != nil {
		return err
	}
	event.TopicType = msgType
	switch msgType {
	// case bmp.PeerStateChangeMsg:
	// 	return a.peerHandler(event)
	case bmp.L3VPNV4Msg:
		return a.l3vpnV4Handler(event)
	case bmp.L3VPNV6Msg:
		return a.l3vpnV6Handler(event)
	}
	return nil
}

func (a *arangoDB) monitor() {
	for {
		select {
		case <-a.stop:
			return
		}
	}
}

// func (a *arangoDB) loadCollection() error {
// 	ctx := context.TODO()

// 	// copy l3vpn v4 prefixes into vrf_v4
// 	glog.Infof("copying l3vpn_v4 prefixes into vrf_v4 collection")
// 	lsn_query := "for l in " + a.lsnode.Name() + " insert l in " + a.lsnodeExt.Name() + ""
// 	cursor, err := a.db.Query(ctx, lsn_query, nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer cursor.Close()

// 	return nil
// }
