package cli

import (
	integrationhelpers "github.com/pointidentity/pointidentity-node/tests/integration/helpers"
	"github.com/pointidentity/pointidentity-node/x/did/types"
)

func QueryDid(did string, container string) (types.QueryDidDocResponse, error) {
	res, err := Query(container, CliBinaryName, "pointidentity", "did-document", did)
	if err != nil {
		return types.QueryDidDocResponse{}, err
	}

	var resp types.QueryDidDocResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return types.QueryDidDocResponse{}, err
	}

	return resp, nil
}
