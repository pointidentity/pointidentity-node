package cli

import (
	integrationhelpers "github.com/pointidentity/pointidentity-node/tests/integration/helpers"
	"github.com/pointidentity/pointidentity-node/x/resource/types"
)

func QueryResource(collectionID string, resourceID string, container string) (types.QueryResourceResponse, error) {
	res, err := Query(container, CliBinaryName, "resource", "specific-resource", collectionID, resourceID)
	if err != nil {
		return types.QueryResourceResponse{}, err
	}

	var resp types.QueryResourceResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return types.QueryResourceResponse{}, err
	}

	return resp, nil
}
