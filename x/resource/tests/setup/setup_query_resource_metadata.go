package setup

import "github.com/pointidentity/pointidentity-node/x/resource/types"

func (s *TestSetup) QueryResourceMetadata(collectionID, resourceID string) (*types.QueryResourceMetadataResponse, error) {
	req := &types.QueryResourceMetadataRequest{
		CollectionId: collectionID,
		Id:           resourceID,
	}

	return s.ResourceQueryServer.ResourceMetadata(s.StdCtx, req)
}
