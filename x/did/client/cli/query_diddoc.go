package cli

import (
	"context"

	"github.com/pointidentity/pointidentity-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetDidDoc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "did-document [id]",
		Short: "Query a DID Document by DID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			did := args[0]
			params := &types.QueryDidDocRequest{
				Id: did,
			}

			resp, err := queryClient.DidDoc(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
