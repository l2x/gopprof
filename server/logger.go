package server

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/cihub/seelog"
)

var logconf = `
<seelog type="asynctimer" asyncinterval="100">
    <outputs formatid="main">
		<filter levels="#debug_output#">
			<console/>
		</filter>
		<filter levels="trace,debug,info,warn,error,critical">
		    <rollingfile type="date" filename="#filename#" datepattern="2006.01.02" maxrolls="7" />
		</filter>
    </outputs>
    <formats>
		<format id="main" format="%Date %Time|%LEVEL|%RelFile:%Line|%FuncShort|%Msg%n"/>
    </formats>
</seelog>
`

var (
	logger seelog.LoggerInterface
)

func initLogger(logPath string, debug bool) error {
	if debug {
		logconf = strings.Replace(logconf, "#debug_output#", "trace,debug,info,warn,error,critical", 1)
	} else {
		logconf = strings.Replace(logconf, "#debug_output#", "off", 1)
	}
	logconf = strings.Replace(logconf, "#filename#", filepath.Join(logPath, "gopprof.log"), 1)

	var err error
	logger, err = seelog.LoggerFromConfigAsString(logconf)
	if err != nil {
		log.Println(err)
		return err
	}
	seelog.ReplaceLogger(logger)
	return nil
}
