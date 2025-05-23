package cli

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/pointidentity/pointidentity-node/tests/integration/helpers"
	"github.com/pointidentity/pointidentity-node/tests/integration/network"

	didtypes "github.com/pointidentity/pointidentity-node/x/did/types"
	resourcetypes "github.com/pointidentity/pointidentity-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

var CLIQueryParams = []string{
	"--chain-id", network.ChainID,
	"--output", OutputFormat,
}

var KeyParams = []string{
	"--output", OutputFormat,
	"--keyring-backend", KeyringBackend,
}

func Query(module, query string, queryArgs ...string) (string, error) {
	args := []string{"query", module, query}

	// Common params
	args = append(args, CLIQueryParams...)

	// Other args
	args = append(args, queryArgs...)

	return Exec(args...)
}

func QueryBalance(address, denom string) (sdk.Coin, error) {
	res, err := Query("bank", "balances", address, "--denom", denom)
	if err != nil {
		return sdk.Coin{}, err
	}

	var resp sdk.Coin
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return sdk.Coin{}, err
	}

	return resp, nil
}

func QueryParams(subspace, key string) (paramproposal.ParamChange, error) {
	res, err := Query("params", "subspace", subspace, key)
	if err != nil {
		return paramproposal.ParamChange{}, err
	}

	var resp paramproposal.ParamChange
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return paramproposal.ParamChange{}, err
	}

	return resp, nil
}

func QueryDidDoc(did string) (didtypes.QueryDidDocResponse, error) {
	res, err := Query("cheqd", "did-document", did)
	if err != nil {
		return didtypes.QueryDidDocResponse{}, err
	}

	var resp didtypes.QueryDidDocResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypes.QueryDidDocResponse{}, err
	}

	return resp, nil
}

func QueryResource(collectionID string, resourceID string) (resourcetypes.QueryResourceResponse, error) {
	res, err := Query("resource", "specific-resource", collectionID, resourceID)
	if err != nil {
		return resourcetypes.QueryResourceResponse{}, err
	}

	var resp resourcetypes.QueryResourceResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryResourceResponse{}, err
	}

	return resp, nil
}

func QueryResourceMetadata(collectionID string, resourceID string) (resourcetypes.QueryResourceMetadataResponse, error) {
	res, err := Query("resource", "metadata", collectionID, resourceID)
	if err != nil {
		return resourcetypes.QueryResourceMetadataResponse{}, err
	}

	var resp resourcetypes.QueryResourceMetadataResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryResourceMetadataResponse{}, err
	}

	return resp, nil
}

func QueryResourceCollection(collectionID string) (resourcetypes.QueryCollectionResourcesResponse, error) {
	res, err := Query("resource", "collection-metadata", collectionID)
	if err != nil {
		return resourcetypes.QueryCollectionResourcesResponse{}, err
	}

	var resp resourcetypes.QueryCollectionResourcesResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryCollectionResourcesResponse{}, err
	}

	return resp, nil
}

func QueryTxn(hash string) (sdk.TxResponse, error) {
	res, err := Query("tx", hash)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var resp sdk.TxResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return resp, nil
}

func QueryProposal(container, id string) (govtypesv1.Proposal, error) {
	fmt.Println("Querying proposal from", container)
	args := append([]string{
		CliBinaryName,
		"query", "gov", "proposal", id,
	}, QueryParamsConst...)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return govtypesv1.Proposal{}, err
	}

	fmt.Println("Proposal", out)

	var resp govtypesv1.Proposal

	err = MakeCodecWithExtendedRegistry().UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return govtypesv1.Proposal{}, err
	}
	return resp, nil
}

func QueryFeemarketParams() (feemarkettypes.Params, error) {
	res, err := Query("feemarket", "params")
	if err != nil {
		return feemarkettypes.Params{}, err
	}

	var resp feemarkettypes.Params
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return feemarkettypes.Params{}, err
	}

	return resp, nil
}

func GetProposalID(rawLog string) (string, error) {
	var logs []sdk.ABCIMessageLog
	err := json.Unmarshal([]byte(rawLog), &logs)
	if err != nil {
		return "", err
	}

	// Iterate over logs and their events
	for _, log := range logs {
		for _, event := range log.Events {
			// Look for the "submit_proposal" event type
			if event.Type == "submit_proposal" {
				for _, attr := range event.Attributes {
					// Look for the "proposal_id" attribute
					if attr.Key == "proposal_id" {
						return attr.Value, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("proposal_id not found")
}

// QueryKeys retrieves the key information and extracts the address
func QueryKeys(name string) (string, error) {
	args := []string{"keys", "show", name}

	args = append(args, KeyParams...)

	output, err := Exec(args...)
	if err != nil {
		return "", err
	}

	var result struct {
		Address string `json:"address"`
	}

	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal JSON")
	}

	return result.Address, nil
}
