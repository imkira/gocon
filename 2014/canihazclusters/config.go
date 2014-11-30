package main

import (
	"fmt"
	"time"

	"github.com/armon/consul-api"
)

//start block OMIT

func get(client *consulapi.Client) ([]byte, error) {
	pair, _, err := client.KV().Get("config", nil)
	if err != nil || pair == nil {
		return nil, err
	}
	return pair.Value, nil
}

func put(client *consulapi.Client, b []byte) error {
	pair := &consulapi.KVPair{Key: "config", Value: b}
	_, err := client.KV().Put(pair, nil)
	return err
}

//end block OMIT

func main() {
	client, _ := consulapi.NewClient(consulapi.DefaultConfig())
	var config []byte
	config, _ = get(client)
	fmt.Printf("get: %s\n", config)
	config = []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	fmt.Printf("put: %s\n", config)
	put(client, config)
	config, _ = get(client)
	fmt.Printf("get: %s\n", config)
}
