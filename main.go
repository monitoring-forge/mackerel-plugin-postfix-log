package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/mackerelio/golib/pluginutil"
	"github.com/monitoring-forge/followparser"
	"github.com/monitoring-forge/mackerel-plugin-postfix-log/postfixlog"
)

var version string
var commit string

const (
	OK = iota
	WARNING
	CRITICAL
	UNKNOWN
)

type Opt struct {
	LogFile       string `long:"logfile" default:"/var/log/maillog" description:"path to postfix/maillog logfile" required:"true"`
	PosFilePrefix string `long:"posfile-prefix" default:"maillog" description:"prefix added position file"`
	Version       bool   `short:"v" long:"version" description:"Show version"`
	Verbose       bool   `short:"V" long:"verbose" description:"Show verbose log"`
}

func (opt *Opt) run() (string, error) {
	pp := postfixlog.NewPostfixParser()
	fp := &followparser.Parser{
		WorkDir:  pluginutil.PluginWorkDir(),
		Callback: pp,
		Silent:   !opt.Verbose,
	}
	_, err := fp.Parse(
		fmt.Sprintf("%s-postfixlog", opt.PosFilePrefix),
		opt.LogFile,
	)
	out := pp.Output(time.Now())
	if err != nil {
		return "", err
	}
	return out, nil
}

func main() {
	os.Exit(_main())
}

func _main() int {
	opt := &Opt{}
	psr := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()
	if opt.Version {
		if commit == "" {
			commit = "dev"
		}
		fmt.Printf(
			"%s-%s\n%s/%s, %s, %s\n",
			filepath.Base(os.Args[0]),
			version,
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
			commit)
		return OK
	} else if flags.WroteHelp(err) {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		return OK
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return UNKNOWN
	}

	output, err := opt.run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return CRITICAL
	}
	fmt.Print(output)

	return OK
}
