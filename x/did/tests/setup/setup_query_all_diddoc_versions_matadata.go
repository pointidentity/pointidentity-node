package setup

import "github.com/pointidentity/pointidentity-node/x/did/types"

func (s *TestSetup) QueryAllDidDocVersionsMetadata(did string) (*types.QueryAllDidDocVersionsMetadataResponse, error) {
	req := &types.QueryAllDidDocVersionsMetadataRequest{
		Id: did,
	}

	return s.QueryServer.AllDidDocVersionsMetadata(s.StdCtx, req)
}
