package main

import (
	"fmt"
	"time"

	"github.com/armon/consul-api"
)

//start block OMIT
func listNew(client *consulapi.Client, opts *consulapi.QueryOptions) {
	entries, meta, err := client.Health().Service("canihazclusters", "", false, opts)
	if err != nil {
		return
	}
	if opts.WaitIndex != meta.LastIndex {
		opts.WaitIndex = meta.LastIndex
		fmt.Printf("%d services\n", len(entries))
		for _, e := range entries {
			printService(e)
		}
	}
}

func printService(e *consulapi.ServiceEntry) {
	for _, c := range e.Checks {
		if c.ServiceName == "canihazclusters" {
			fmt.Printf("%s %s:%d = %s\n", e.Service.ID, e.Node.Address, e.Service.Port, c.Status)
			return
		}
	}
}

//end block OMIT

func discover(stop, done chan int) {
	fmt.Printf("Detecting services...\n")
	defer close(done)

	client, _ := consulapi.NewClient(consulapi.DefaultConfig())
	opts := &consulapi.QueryOptions{WaitIndex: 0, WaitTime: time.Second}
	for {
		select {
		case <-stop:
			return
		default:
			listNew(client, opts)
		}
	}
}

func main() {
	// discover
	stop := make(chan int)
	done := make(chan int)
	go discover(stop, done)
	defer func() {
		close(stop)
		<-done
	}()

	// providers
	time.Sleep(1 * time.Second)
	ensureProvider("node1", 1337)
	time.Sleep(1 * time.Second)
	ensureProvider("node2", 1338)

	time.Sleep(5 * time.Second)
}

func ensureProvider(id string, port int) {
	stop := make(chan int)
	done := make(chan int)
	go provider(id, port, stop, done)
	time.Sleep(3 * time.Second)
	close(stop)
	<-done
}

func provider(id string, port int, stop, done chan int) {
	defer func() {
		close(done)
		fmt.Printf("- Stopping service %s...\n", id)
	}()

	fmt.Printf("+ Starting service %s...\n", id)
	client, _ := consulapi.NewClient(consulapi.DefaultConfig())
	err := register(client, id, port)
	if err != nil {
		panic(err)
	}
	defer unregister(client, id)
	for {
		select {
		case <-stop:
			return
		case <-time.After(time.Second):
			keepAlive(client, id)
		}
	}
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
