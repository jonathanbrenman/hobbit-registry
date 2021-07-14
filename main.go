package main

import (
	"fmt"
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
	clientHttp := clients.NewHttpClient(fmt.Sprintf("%s://%s:%d",
		hobbit.Configs.Registry.Scheme,
		hobbit.Configs.Registry.URL,
		hobbit.Configs.Registry.Port,
	))
	clientHttp.CheckConnectivity()

	// Docker client
	docker := clients.NewDockerClient(
		hobbit.Configs.Registry.Username,
		hobbit.Configs.Registry.Password,
		fmt.Sprintf("%s:%d",
			hobbit.Configs.Registry.URL,
			hobbit.Configs.Registry.Port,
		),
	)

	for _, image := range hobbit.Configs.Images {
		doesImageExists := clientHttp.CheckImage(image)
		if doesImageExists {
			log.Printf("Image %s already exist in the private registry, I'll skip this %s.\n", image, image)
			continue
		}

		// Pulling image
		log.Printf("Pulling %s, please wait..\n", image)
		if err := docker.Pull(image); err != nil {
			log.Printf("[ Error ] pulling image %s %s\n", image, err.Error())
			continue
		}

		// Tagging image
		log.Printf("Tagging %s, please wait..\n", image)
		if err := docker.Tag(image); err != nil {
			log.Printf("[ Error ] tagging image %s %s\n", image, err.Error())
			continue
		}
		log.Printf("Tagged %s ok!\n", image)

		// Pushing image
		log.Printf("Pushing %s, please wait..\n", image)
		if err := docker.Push(image); err != nil {
			log.Printf("[ Error ] pushing image %s %s\n", image, err.Error())
			continue
		}

		// Delete image
		if err := docker.Delete(image); err != nil {
			log.Printf("[ Error ] deleting image %s %s\n", image, err.Error())
			continue
		}
	}

}