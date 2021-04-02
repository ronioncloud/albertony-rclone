package about

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rclone/rclone/cmd"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config/flags"
	"github.com/rclone/rclone/fs/operations"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
)

func init() {
	cmd.Root.AddCommand(commandDefinition)
	cmdFlags := commandDefinition.Flags()
	flags.BoolVarP(cmdFlags, &jsonOutput, "json", "", false, "Format output as JSON")
}

// printValue formats uv to be output
func printValue(what string, uv *int64, humanReadable bool, isSize bool) {
	what += ":"
	if uv == nil {
		return
	}
	if isSize {
		fmt.Printf("%-9s%s\n", what, operations.SizeString(*uv, humanReadable))
	} else {
		fmt.Printf("%-9s%s\n", what, operations.CountString(*uv, humanReadable))
	}
}

var commandDefinition = &cobra.Command{
	Use:   "about remote:",
	Short: `Get quota information from the remote.`,
	Long: `
` + "`rclone about`" + ` prints quota information about a remote to standard
output. The output is typically used, free, quota and trash contents.

E.g. Typical output from ` + "`rclone about remote:`" + ` is:

    Total:   18253611008
    Used:    7993453766
    Free:    1411001220
    Trashed: 104857602
    Other:   8849156022

Where the fields are:

  * Total: Total size available.
  * Used: Total size used.
  * Free: Total space available to this user.
  * Trashed: Total space used by trash.
  * Other: Total amount in other storage (e.g. Gmail, Google Photos).
  * Objects: Total number of objects in the storage.

All sizes are in number of bytes.

Applying global flag ` + "`--human-readable`" + ` to the command prints, e.g.

    Total:   17Gi
    Used:    7.444Gi
    Free:    1.315Gi
    Trashed: 100.000Mi
    Other:   8.241Gi

A ` + "`--json`" + ` flag generates conveniently computer readable output, e.g.

    {
        "total": 18253611008,
        "used": 7993453766,
        "trashed": 104857602,
        "other": 8849156022,
        "free": 1411001220
    }

Not all backends print all fields. Information is not included if it is not
provided by a backend. Where the value is unlimited it is omitted.

Some backends does not support the ` + "`rclone about`" + ` command at all,
see complete list in [documentation](https://rclone.org/overview/#optional-features).
`,
	Run: func(command *cobra.Command, args []string) {
		cmd.CheckArgs(1, 1, command, args)
		f := cmd.NewFsSrc(args)
		cmd.Run(false, false, command, func() error {
			doAbout := f.Features().About
			if doAbout == nil {
				return errors.Errorf("%v doesn't support about", f)
			}
			ctx := context.Background()
			ci := fs.GetConfig(context.Background())
			u, err := doAbout(ctx)
			if err != nil {
				return errors.Wrap(err, "About call failed")
			}
			if u == nil {
				return errors.New("nil usage returned")
			}
			if jsonOutput {
				out := json.NewEncoder(os.Stdout)
				out.SetIndent("", "\t")
				return out.Encode(u)
			}

			printValue("Total", u.Total, ci.HumanReadable, true)
			printValue("Used", u.Used, ci.HumanReadable, true)
			printValue("Free", u.Free, ci.HumanReadable, true)
			printValue("Trashed", u.Trashed, ci.HumanReadable, true)
			printValue("Other", u.Other, ci.HumanReadable, true)
			printValue("Objects", u.Objects, ci.HumanReadable, false)
			return nil
		})
	},
}
