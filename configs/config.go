package configs

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type HobbitConfig struct {
	Configs Configs
}

type Configs struct {
	Registry string
  	Images []string
}

func (hcfg *HobbitConfig) Parse() *string {
	// Receive config path as an argument via command line.
	var configPath = flag.String("c", "", "config file path.")
	flag.Parse()
	if *configPath == "" {
		log.Fatal("you forgot the config file path, please try again with -c <your yaml location>")
	}
	return configPath
}

func (hcfg *HobbitConfig) LoadConfig(path string) *HobbitConfig {
	yamlFile, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatal("[ Error ] - reading config yaml file. " + err.Error())
    }
    err = yaml.Unmarshal(yamlFile, &hcfg)
    if err != nil {
        log.Fatal("[ Error ] - format not valid for config file hobbit.yaml. " + err.Error())
    }
    return hcfg
}

func (hcfg *HobbitConfig) Validate() {
	if hcfg.Configs.Registry == "" {
		log.Fatal("[ ERROR ] - private registry not found in config file.")
	}
	if len(hcfg.Configs.Images) == 0 {
		log.Fatal("[ ERROR ] - you need at least one image to push to the remote registry.")
	}
}