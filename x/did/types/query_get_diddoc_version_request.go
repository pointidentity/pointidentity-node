package types

import (
	"github.com/pointidentity/pointidentity-node/x/did/utils"
)

func (query *QueryDidDocVersionRequest) Normalize() {
	query.Id = utils.NormalizeDID(query.Id)
	query.Version = utils.NormalizeUUID(query.Version)
}
