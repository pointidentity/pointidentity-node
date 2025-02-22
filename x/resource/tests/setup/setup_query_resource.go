package setup

import "github.com/pointidentity/pointidentity-node/x/resource/types"

func (s *TestSetup) QueryResource(collectionID, resourceID string) (*types.QueryResourceResponse, error) {
	req := &types.QueryResourceRequest{
		CollectionId: collectionID,
		Id:           resourceID,
	}

	return s.ResourceQueryServer.Resource(s.StdCtx, req)
}
