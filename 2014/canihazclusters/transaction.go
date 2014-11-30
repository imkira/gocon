package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/armon/consul-api"
)

func runNode(client *consulapi.Client, nodeID string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("%s starting...\n", nodeID)
	defer fmt.Printf("%s terminating...\n", nodeID)

	register(client, nodeID, port)
	defer unregister(client, nodeID)
	keepAlive(client, nodeID)
	runTransaction(client, nodeID)
}

//start block OMIT

func runTransaction(client *consulapi.Client, nodeID string) {
	se := &consulapi.SessionEntry{
		Name:   "keyLock",
		Checks: []string{"serfHealth", "service:" + nodeID},
	}
	sessionID, _, err := client.Session().Create(se, nil)
	if err != nil {
		return
	}
	defer client.Session().Destroy(sessionID, nil)

	for {
		p := &consulapi.KVPair{Key: "lockedKey", Value: nil, Session: sessionID}
		if ok, _, err := client.KV().Acquire(p, nil); err != nil || !ok {
			time.Sleep(1 * time.Millisecond)
			continue
		}
		defer client.KV().Release(p, nil)
		slowTransaction(nodeID)
		return
	}
}

//end block OMIT

func slowTransaction(nodeID string) {
	fmt.Printf("+ %s beginning transaction...\n", nodeID)
	time.Sleep(2 * time.Second)
	fmt.Printf("+ %s ending transaction...\n", nodeID)
}

func main() {
	wg := &sync.WaitGroup{}
	client, _ := consulapi.NewClient(consulapi.DefaultConfig())
	port := 1400
	for i := 1; i <= 5; i += 1 {
		wg.Add(1)
		go runNode(client, fmt.Sprintf("node%d", i), port+i, wg)
	}
	wg.Wait()
}

func register(client *consulapi.Client, id string, port int) error {
	service := &consulapi.AgentServiceRegistration{
		ID:    id,
		Name:  "canihazclusters",
		Tags:  []string{"gocon", "autumn"},
		Port:  port,
		Check: &consulapi.AgentServiceCheck{TTL: "5s"},
	}
	return client.Agent().ServiceRegister(service)
}

func keepAlive(client *consulapi.Client, id string) error {
	return client.Agent().PassTTL("service:"+id, "some status")
}

func unregister(client *consulapi.Client, id string) error {
	return client.Agent().ServiceDeregister(id)
}
