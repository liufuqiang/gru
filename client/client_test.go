// +build integration

package client

import (
	"time"
	"testing"

	"github.com/dnaeon/gru/minion"

	"golang.org/x/net/context"
	etcdclient "github.com/coreos/etcd/client"
)

// Default config for etcd minions and clients
var defaultEtcdConfig = etcdclient.Config{
	Endpoints:               []string{"http://127.0.0.1:2379", "http:127.0.0.1:4001"},
	Transport:               etcdclient.DefaultTransport,
	HeaderTimeoutPerRequest: time.Second,
}

// Cleans up the minion space in etcd after tests
func cleanupAfterTest(t *testing.T) {
	c, err := etcdclient.New(defaultEtcdConfig)
	if err != nil {
		t.Fatal(err)
	}

	kapi := etcdclient.NewKeysAPI(c)
	deleteOpts := &etcdclient.DeleteOptions{
		Recursive: true,
		Dir: true,
	}

	_, err = kapi.Delete(context.Background(), minion.EtcdMinionSpace, deleteOpts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMinionList(t *testing.T) {
	defer cleanupAfterTest(t)

	minions := []minion.Minion{
		minion.NewEtcdMinion("Bob", defaultEtcdConfig),
		minion.NewEtcdMinion("Kevin", defaultEtcdConfig),
		minion.NewEtcdMinion("Stuart", defaultEtcdConfig),
	}

	// Start our minions
	for _, m := range minions {
		m.Serve()
		defer m.Stop()
	}

	klient := NewEtcdMinionClient(defaultEtcdConfig)
	minionList, err := klient.MinionList()
	if err != nil {
		t.Fatal(err)
	}

	want := len(minions)
	got := len(minionList)

	if want != got {
		t.Errorf("want %d minion, got %d minion(s)", want, got)
	}
}
