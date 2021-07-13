package main

import (
	"fmt"
	"hobbit-registry/clients"
	"hobbit-registry/configs"
)

func main() {
	// Load configs from yaml with hobbit format specified on readme.
	var hobbit configs.HobbitConfig
	configPath := hobbit.Parse()
	hobbit.LoadConfig(*configPath).Validate()

	// Check connectivity with the private registry.
	clientHttp := clients.NewHttpClient(hobbit.Configs.Registry)
	clientHttp.CheckConnectivity()

	fmt.Println(hobbit.Configs)
}