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
	"github.com/urfave/cli/v2"
)

const (
	flagCfg         = "config"
	flagCfgAlias    = "c"
	flagDryRun      = "dry-run"
	flagDryRunAlias = "d"
)

// Run is the app starter
func Run(args []string) error {
	app := cli.NewApp()
	app.Name = "logsync"
	app.Version = "1.0.1"
	app.Usage = "Event Forwarding"
	app.Copyright = "Swiss Learning Hub AG"
	app.Commands = []*cli.Command{
		newCmdRun(),
	}
	return app.Run(args)
}

type command struct {
	args   []cliArg
	argMap map[string]string
}

func (cmd *command) bootstrap(beforeFunc cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {

		var err error

		if err = cmd.bootstrapArgs(c); err != nil {
			return err
		}

		if beforeFunc != nil {
			return beforeFunc(c)
		}

		return nil
	}
}

func (cmd *command) bootstrapArgs(c *cli.Context) error {

	cmd.setCmdPrinter()

	if c.Args().Len() < len(cmd.args) {
		_ = cli.ShowSubcommandHelp(c)
		return fmt.Errorf("expected arguments: %d; given arguments: %d", len(cmd.args), c.Args().Len())
	}

	cmd.argMap = map[string]string{}
	for _, arg := range cmd.args {
		v, err := cmd.getArgFromContext(c, arg.Name)
		if err != nil {
			return err
		}
		cmd.argMap[arg.Name] = v
	}

	return nil
}
