package types

import (
	"github.com/pointidentity/pointidentity-node/x/did/utils"
)

func (query *QueryDidDocRequest) Normalize() {
	query.Id = utils.NormalizeDID(query.Id)
}
