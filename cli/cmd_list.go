package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all OWASP cheat sheets",
		RunE: func(cmd *cobra.Command, _ []string) error {
			limit := a.effectiveLimit(0)
			sheets, err := a.client.List(cmd.Context(), limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(sheets, len(sheets))
		},
	}
	return cmd
}
