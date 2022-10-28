package env

import (
	"errors"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"shipyard/requests"
	"shipyard/requests/uri"
)

func NewRebuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rebuild",
		Short: "Rebuild an environment",
	}

	cmd.AddCommand(newRebuildEnvironmentCmd())

	return cmd
}

func newRebuildEnvironmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Aliases:      []string{"env"},
		Use:          "environment",
		Short:        "Rebuild a running environment",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return rebuildEnvironmentByID(args[0])
			}
			return errors.New("missing environment ID")
		},
	}

	return cmd
}

func rebuildEnvironmentByID(id string) error {
	client, err := requests.NewHTTPClient(os.Stdout)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	org := viper.GetString("org")
	if org != "" {
		params["org"] = org
	}

	body, err := client.Do(http.MethodPost, uri.CreateResourceURI("rebuild", "environment", id, params), nil)
	if err != nil {
		return err
	}

	return client.Write(body)
}
