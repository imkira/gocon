package main

import (
	"fmt"
	"time"

	"github.com/armon/consul-api"
)

//start block OMIT

func register(client *consulapi.Client) error {
	service := &consulapi.AgentServiceRegistration{
		ID:    "node1",
		Name:  "canihazclusters",
		Tags:  []string{"gocon", "autumn"},
		Port:  1337,
		Check: &consulapi.AgentServiceCheck{TTL: "30s"},
	}
	return client.Agent().ServiceRegister(service)
}

func keepAlive(client *consulapi.Client) error {
	return client.Agent().PassTTL("service:node1", "some status")
}

func unregister(client *consulapi.Client) error {
	return client.Agent().ServiceDeregister("node1")
}

//end block OMIT

func main() {
	client, _ := consulapi.NewClient(consulapi.DefaultConfig())

	// register service
	fmt.Println("+ Registering Service...")
	err := register(client)
	if err != nil {
		panic(err)
	}

	// dont forget to unregister
	defer func() {
		fmt.Println("- Unregistering Service...")
		unregister(client)
	}()

	// keep service alive...
	for i := 0; i < 5; i += 1 {
		time.Sleep(time.Second)
		fmt.Println("> Keeping Service Alive...")
		keepAlive(client)
	}
}
