package types

import (
	"github.com/pointidentity/pointidentity-node/x/did/utils"
)

func (query *QueryResourceRequest) Normalize() {
	query.CollectionId = utils.NormalizeID(query.CollectionId)
	query.Id = utils.NormalizeUUID(query.Id)
}
