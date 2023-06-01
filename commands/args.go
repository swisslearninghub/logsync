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
	"io"
	"strings"
)

var commandHelpTemplate = `NAME:
    {{.HelpName}} - {{.Usage}}

USAGE:
    {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

CATEGORY:
    {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
    {{.Description | nindent 3 | trim}}{{end}}{{if .MyArgs}}

ARGUMENTS:
    {{range .MyArgs}}{{.}}
    {{end}}{{end}}{{if .VisibleFlags}}
OPTIONS:
    {{range .VisibleFlags}}{{.}}
    {{end}}{{end}}
`

type helpCommand struct {
	cli.Command
	MyArgs []cliArg
}

type cliArg struct {
	Name  string
	Usage string
}

// String returns representation for custom CommandHelpText
func (arg *cliArg) String() string {
	return fmt.Sprintf("<%s>\t\t%s", arg.Name, arg.Usage)
}

func (cmd *command) genArgsUsage() string {
	if len(cmd.args) == 0 {
		return ""
	}
	var args []string
	for _, arg := range cmd.args {
		args = append(args, arg.Name)
	}
	return fmt.Sprintf("<" + strings.Join(args, "> <") + ">")
}

func (cmd *command) getArgFromContext(c *cli.Context, name string) (string, error) {
	for x := 0; x < len(cmd.args); x++ {
		if cmd.args[x].Name != name {
			continue
		}
		if c.Args().Len() < x+1 {
			return "", fmt.Errorf("argument not found: %s", name)
		}
		v := strings.TrimSpace(c.Args().Get(x))
		if v == "" {
			return "", fmt.Errorf("argument empty: %s", name)
		}
		return v, nil
	}
	return "", fmt.Errorf("argument not found: %s", name)
}

// func (cmd *command) getArg(name string) string {
//	if v, ok := cmd.argMap[name]; ok {
//		return v
//	}
//	return ""
//}

func (cmd *command) setCmdPrinter() {
	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		c, ok := data.(*cli.Command)
		if ok {
			cli.HelpPrinterCustom(w, commandHelpTemplate, &helpCommand{Command: *c, MyArgs: cmd.args}, nil)
			return
		}
		cli.HelpPrinterCustom(w, templ, data, nil)
	}
}
