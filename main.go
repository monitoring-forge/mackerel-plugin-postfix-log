package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mackerelio/golib/pluginutil"
	"github.com/monitoring-forge/flagrun"
	"github.com/monitoring-forge/followparser"
	"github.com/monitoring-forge/mackerel-plugin-postfix-log/postfixlog"
)

var version string

type Opt struct {
	LogFile       string `long:"logfile" default:"/var/log/maillog" description:"path to postfix/maillog logfile" required:"true"`
	PosFilePrefix string `long:"posfile-prefix" default:"maillog" description:"prefix added position file"`
	Version       bool   `short:"v" long:"version" description:"Show version"`
	Verbose       bool   `short:"V" long:"verbose" description:"Show verbose log"`
}

func (opt *Opt) Run(_ []string) (any, int) {
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
		return err, flagrun.CRITICAL
	}
	return out, flagrun.OK
}

func main() {
	opt := &Opt{}
	os.Exit(flagrun.Go(opt, flagrun.Version(version)))
}
