package cmd

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	argv struct {
		help            bool
		version         bool
		flunetbitConfig string
		nProc           uint
		log             string
		pprofHostPort   string
		debug           bool
	}

	logFd *os.File
)

var (
	BuildTime    string
	BuildOSUname string
	BuildCommit  string
	buildVersion string // concatination of Build* into a single string
)

func init() {
	buildVersion = fmt.Sprintf(`Dpp server compiled at %s by %s after %s on %s`, BuildTime, runtime.Version(),
		BuildCommit, BuildOSUname,
	)

	//log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	// actions
	flag.BoolVar(&argv.help, `h`, false, `show this help`)
	flag.BoolVar(&argv.version, `version`, false, `show version`)
	flag.BoolVar(&argv.debug, "debug", false, "show debug messages")

	flag.UintVar(&argv.nProc, `cores`, uint(0), `max cpu cores usage`)
	flag.StringVar(&argv.log, `l`, "", "log file (if needed)")
	flag.StringVar(&argv.pprofHostPort, `pprof`, ``, `host:port for http pprof`)

	// Fluent bit proxy options
	flag.StringVar(&argv.flunetbitConfig, `config`, ``, `path to fluentbit proxy config`)

	flag.Parse()
}

func updateThread(ch chan os.Signal) {
	for range ch {
		reopenLog()
	}
}

func reopenLog() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	if argv.log == "" {
		return
	}

	var err error

	logFd, err = os.OpenFile(argv.log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf(`Cannot log to file "%s": %s`, argv.log, err.Error()))
		return
	}
	log.SetOutput(logFd)
}

func Main() {
	if argv.version {
		fmt.Fprint(os.Stderr, buildVersion, "\n")
		return
	} else if argv.help {
		flag.Usage()
		return
	}

	if argv.debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if argv.nProc > 0 {
		runtime.GOMAXPROCS(int(argv.nProc))
	} else {
		argv.nProc = uint(runtime.NumCPU())
	}

	if argv.pprofHostPort != `` {
		go func() {
			if err := http.ListenAndServe(argv.pprofHostPort, nil); err != nil {
				log.Printf(`pprof listen fail: %s`, err.Error())
			}
		}()
	}

	updCh := make(chan os.Signal, 10)
	signal.Notify(updCh, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)
	reopenLog()
	go updateThread(updCh)

	err := fluentbit.RunDPPProxy(argv.flunetbitConfig)
	if err != nil {
		log.Fatalf("Could not run fluentbit proxy: %s", err.Error())
	}
}
