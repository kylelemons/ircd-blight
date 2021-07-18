package main

import (
	"flag"
	"os"

	"github.com/kylelemons/ircd-blight/old/ircd/core"
	"github.com/kylelemons/ircd-blight/old/ircd/log"
)

var (
	// Normal execution parameters
	config  = flag.String("config", "/etc/ircd.conf", "The configuration file to use")
	logfile = flag.String("log", "/var/log/ircd.log", "The file to which logs are written")

	// Flags
	silent = flag.Bool("silent", false, "Don't write logs to the console")

	// Other execution modes
	genconf   = flag.Bool("genconf", false, "Genereate a configuration file and exit")
	checkconf = flag.Bool("checkconf", false, "Check the configuration file and exit")
)

func main() {
	flag.Parse()

	if *genconf {
		conf, err := os.Create(*config)
		if err != nil {
			log.Error.Fatalf("Opening config file %q for writing: %s", *config, err)
		}
		_, err = conf.WriteString(core.DefaultXML)
		if err != nil {
			log.Error.Fatalf("Writing default configuration to %q: %s", *config, err)
		}
		log.Info.Printf("Configuration file written to %q", *config)
		os.Exit(0)
	}

	if err := log.SetFile(*logfile); err != nil {
		log.Error.Fatalf("Opening logfile: %s", err)
	}
	if !*silent {
		log.ShowInConsole()
	}

	if err := core.LoadConfigFile(*config); err != nil {
		log.Error.Fatalf("Loading config: %s", err)
	}

	if *checkconf {
		if !core.CheckConfig() {
			log.Error.Fatalf("Invalid configuration")
		}
		log.Info.Printf("Configuration successfully checked.")
	}

	core.Start()
}
