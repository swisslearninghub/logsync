// Copyright 2023 Swiss Learning Hub AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"fmt"
	"github.com/swisslearninghub/logsync/api"
	"github.com/swisslearninghub/logsync/cefsyslog"
	"github.com/swisslearninghub/logsync/config"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// CmdRun ...
type CmdRun struct {
	command
	cfg     *config.Config
	api     *api.HubAPI
	logfile *os.File
	ceflog  *cefsyslog.Writer
}

// newRealmPolicySetDefault ...
func newCmdRun() *cli.Command {

	cmd := &CmdRun{
		command: command{
			args: []cliArg{
				// {"argname", "Argument description."},
			},
		},
	}

	return &cli.Command{
		Name:        "run",
		Description: "Run with given configuration",
		Before:      cmd.bootstrap(cmd.before),
		Action:      cmd.action,
		ArgsUsage:   cmd.genArgsUsage(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      flagCfg,
				Usage:     "use `FILE` as config",
				Aliases:   []string{flagCfgAlias},
				TakesFile: true,
				Value:     "logsync.json",
			},
			&cli.BoolFlag{
				Name:    flagDryRun,
				Usage:   "do not report to syslog server",
				Aliases: []string{flagDryRunAlias},
			},
		},
	}
}

// before bootstraps run
func (cmd *CmdRun) before(c *cli.Context) error {

	var err error

	if err = cmd.setConfig(c); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	if err = cmd.setLogging(); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	if err = cmd.setAPI(); err != nil {
		log.Println(err.Error())
		cmd.close()
		return cli.Exit(err.Error(), 1)
	}

	if err = cmd.setSyslog(); err != nil {
		log.Println(err.Error())
		cmd.close()
		return cli.Exit(err.Error(), 1)
	}

	return nil
}

// action executes command
func (cmd *CmdRun) action(c *cli.Context) error {

	defer cmd.close()

	log.Println("Starting " + c.App.Name)

	log.Printf("[Option] dryRun: %v\n", c.Bool(flagDryRun))

	values := cmd.values()
	for k, v := range values {
		log.Printf("[Query] %s: %v\n", k, v)
	}

	events, err := cmd.api.QueryClientEvents(values)
	if err != nil {
		log.Println(err.Error())
		log.Println("Exiting")
		return cli.Exit(err.Error(), 1)
	}

	log.Printf("Retrieved %d event(s)\n", len(events))

	if len(events) == 0 {
		log.Println("Exiting")
		return nil
	}

	var reported int

	for _, ev := range events {
		reported += cmd.report(ev, c.Bool(flagDryRun))
	}

	log.Printf("Reported %d event(s)\n", reported)
	log.Println("Exiting")

	return nil
}

// report handles configuret report checks and returns *CEF if needed. nil otherwise
func (cmd *CmdRun) report(er api.EventRepresentation, dryRun bool) int {
	reported := 0
	for _, detection := range cmd.cfg.Detections {
		if detection.Report(&er) {
			cef := detection.CEF()
			cmd.updFromEvent(cef, &er)
			log.Printf("[%d] %s", er.Time, cef.String())
			if !dryRun {
				if err := cmd.ceflog.Log(time.UnixMilli(er.Time), detection.LogLevel, cef.String()); err != nil {
					log.Printf("[%d] %s\n", er.Time, err.Error())
					continue
				}
			}
			reported++
		}
	}
	return reported
}

// updFromEvent enriches cef with event attributes and details
func (cmd *CmdRun) updFromEvent(cef *cefsyslog.CEF, er *api.EventRepresentation) {
	cef.Extension[cefsyslog.ExtSourceUserName] = er.GetDetail("username", "unknown")
	cef.Extension[cefsyslog.ExtReceiptTime] = fmt.Sprintf("%d", er.Time)
	if er.UserID != nil {
		cef.Extension[cefsyslog.ExtSourceUserID] = *er.UserID
	}
	if er.IPAddress != nil {
		cef.Extension[cefsyslog.ExtSourceAddress] = *er.IPAddress
	}
}

// values returns url.Values for api request
func (cmd *CmdRun) values() url.Values {
	const timeDay = time.Hour * 24
	dateTo := time.Now()
	dateFrom := dateTo.Add(-(timeDay * time.Duration(cmd.cfg.Filter.Days)))
	values := url.Values{
		api.QueryParamFrom: []string{dateFrom.Format(api.EventDateLayout)},
		api.QueryParamTo:   []string{dateTo.Format(api.EventDateLayout)},
		api.QueryParamMax:  []string{api.EventMax},
	}
	if cmd.cfg.Filter.Max > 0 {
		values[api.QueryParamMax] = []string{fmt.Sprintf("%d", cmd.cfg.Filter.Max)}
	}
	if len(cmd.cfg.Filter.Type) > 0 {
		values[api.QueryParamType] = cmd.cfg.Filter.Type
	}
	return values
}

// setConfig loads configuration
func (cmd *CmdRun) setConfig(c *cli.Context) error {

	var app string
	var cwd string
	var err error

	if cwd, err = os.Getwd(); err != nil {
		return err
	}

	if app, err = filepath.Abs(os.Args[0]); err != nil {
		return err
	}

	cmd.cfg, err = config.NewFromFiles(
		c.String(flagCfg),
		cwd+"/logsync.json",
		app+"/logsync.json",
	)

	return err
}

// setLogging initializes logging
func (cmd *CmdRun) setLogging() error {
	var err error
	log.SetOutput(os.Stdout)
	const perm = 0600
	if cmd.cfg.Logfile != "" {
		if cmd.logfile, err = os.OpenFile(cmd.cfg.Logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm); err != nil {
			return err
		}
		log.SetOutput(io.MultiWriter(os.Stdout, cmd.logfile))
	}
	return nil
}

// setLogging initializes http app
func (cmd *CmdRun) setAPI() error {
	var err error
	cmd.api, err = api.NewAPI(
		cmd.cfg.OAuth2.ClientID,
		cmd.cfg.OAuth2.Secret,
		cmd.cfg.OAuth2.TokenURL,
		cmd.cfg.OAuth2.ContextURL,
	)
	return err
}

// setLogging initializes syslog client
func (cmd *CmdRun) setSyslog() error {
	var err error
	cmd.ceflog, err = cefsyslog.SyslogWriterDial(
		cmd.cfg.Syslog.Proto,
		cmd.cfg.Syslog.Address,
		cmd.cfg.Syslog.Facility,
		cmd.cfg.Syslog.Tag,
	)
	return err
}

// close takes care about open resources
func (cmd *CmdRun) close() {
	cmd.api = nil
	if cmd.ceflog != nil {
		_ = cmd.ceflog.Close()
		cmd.ceflog = nil
	}
	if cmd.logfile != nil {
		_ = cmd.logfile.Close()
		cmd.logfile = nil
	}
}
