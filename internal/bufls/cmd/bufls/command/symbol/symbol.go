// Copyright 2022 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package symbol

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/bufbuild/buf-language-server/internal/bufls"
	"github.com/bufbuild/buf-language-server/internal/bufls/buflscli"
	"github.com/bufbuild/buf/private/pkg/app/appcmd"
	"github.com/bufbuild/buf/private/pkg/app/appflag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewCommand returns a new Command.
func NewCommand(
	name string,
	builder appflag.Builder,
) *appcmd.Command {
	return &appcmd.Command{
		Use:   name + " <file>",
		Short: "Provides information about the symbol at the given file.",
		Args:  cobra.MaximumNArgs(1),
		Run: builder.NewRunFunc(
			func(ctx context.Context, container appflag.Container) error {
				engine, err := buflscli.NewEngine(
					ctx,
					container,
					true,
				)
				if err != nil {
					return err
				}
				fp := container.Arg(0)
				// validate file path
				if _, err := os.Stat(
					fp,
				); err != nil {
					return err
				}
				sbs, err := engine.Symbols(ctx, bufls.FilePath(fp))
				if err != nil {
					return err
				}
				buf := bytes.NewBuffer(nil)
				for _, sb := range sbs {
					fmt.Fprintf(
						buf,
						"%s:%s\n",
						sb.Name,
						sb.Kind.String(),
					)
				}
				if _, err := container.Stdout().Write(
					buf.Bytes(),
				); err != nil {
					return err
				}
				return nil
			},
		),
		BindFlags: func(fs *pflag.FlagSet) {},
	}
}
