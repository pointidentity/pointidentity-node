package setup

import "github.com/pointidentity/pointidentity-node/x/resource/types"

func (s *TestSetup) CollectionResources(collectionID string) (*types.QueryCollectionResourcesResponse, error) {
	req := &types.QueryCollectionResourcesRequest{
		CollectionId: collectionID,
	}

	return s.ResourceQueryServer.CollectionResources(s.StdCtx, req)
}
