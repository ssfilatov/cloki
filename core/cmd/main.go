package cmd

import (
	"flag"
	//"github.com/cortexproject/cortex/pkg/util/flagext"
	log "github.com/sirupsen/logrus"
	cloki "github.com/ssfilatov/cloki/core"
	"github.com/ssfilatov/cloki/core/server"
	"os"
)

func Main() {
	var (
		cfg        cloki.Config
		configFile = ""
	)
	flag.StringVar(&configFile, "config.file", "", "Configuration file to load.")
	//flagext.RegisterFlags(&cfg)
	flag.Parse()

	if configFile != "" {
		if err := cloki.LoadConfig(configFile, &cfg); err != nil {
			log.Errorf("error loading config %s: %s", configFile, err)
			os.Exit(1)
		}
	}

	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.Debugf("config: %s", cfg.LabelList)

	t, err := server.New(&cfg)
	if err != nil {
		log.Errorf("error initializing cloki: %s", err)
		os.Exit(1)
	}

	if err := t.Run(); err != nil {
		log.Errorf("error running cloki: %s", err)
		os.Exit(1)
	}
}
