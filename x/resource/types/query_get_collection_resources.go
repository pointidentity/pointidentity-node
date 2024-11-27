package types

import (
	"github.com/pointidentity/pointidentity-node/x/did/utils"
)

func (query *QueryCollectionResourcesRequest) Normalize() {
	query.CollectionId = utils.NormalizeID(query.CollectionId)
}
