package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/quells/rsrp"
)

func main() {
	config, err := parseArgs(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(usage())
		os.Exit(1)
	}

	if config == nil {
		fmt.Println("could not parse config")
		os.Exit(1)
	}

	routes, err := rsrp.ConvertRules(config.Routes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// This is the only line which uses rsrp;
	// everything else is just parsing command line arguments and setting up the HTTP server.
	http.HandleFunc("/", rsrp.RouteAll(*routes))

	httpServer := &http.Server{
		Addr: ":5000",
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func usage() string {
	return "rsrp <CONFIG_FILE>"
}

func parseArgs(args []string) (*rsrp.Config, error) {
	switch len(args) {
	case 2:
		configFilename := args[1]
		data, err := ioutil.ReadFile(configFilename)
		if err != nil {
			return nil, err
		}

		config := &rsrp.Config{}
		err = json.Unmarshal(data, config)

		return config, err
	default:
		return nil, fmt.Errorf("exactly 1 argument expected")
	}
}
