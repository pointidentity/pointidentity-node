package cli

import (
	"context"

	"github.com/pointidentity/pointidentity-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetDidDocVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "did-version [id] [version-id]",
		Short: "Query specific version of a DID Document",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			did := args[0]
			versionID := args[1]
			params := &types.QueryDidDocVersionRequest{
				Id:      did,
				Version: versionID,
			}

			resp, err := queryClient.DidDocVersion(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
