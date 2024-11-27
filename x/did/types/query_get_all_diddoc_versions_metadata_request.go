package types

import (
	"github.com/pointidentity/pointidentity-node/x/did/utils"
)

func (query *QueryAllDidDocVersionsMetadataRequest) Normalize() {
	query.Id = utils.NormalizeDID(query.Id)
}
