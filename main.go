package main

import (
	"hobbit-registry/clients"
	"hobbit-registry/configs"
	"log"
)

func main() {
	// Load configs from yaml with hobbit format specified on readme.
	var hobbit configs.HobbitConfig
	configPath := hobbit.Parse()
	hobbit.LoadConfig(*configPath).Validate()

	// Check connectivity with the private registry.
	clientHttp := clients.NewHttpClient(hobbit.Configs.Registry)
	clientHttp.CheckConnectivity()

	// Docker client
	docker := clients.NewDockerClient()
	for _, image := range hobbit.Configs.Images {
		doesImageExists := clientHttp.CheckImage(image)
		if doesImageExists {
			log.Printf("Image %s already exist in the private registry, I'll skip this %s.", image, image)
			continue
		}
		log.Printf("Pulling %s, please wait..", image)
		if err := docker.Pull(image); err != nil {
			log.Printf("[ Error ] pulling image %s\n", image)
			continue
		}
	}

}