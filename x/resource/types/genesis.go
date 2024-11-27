package types

import (
	"fmt"
)

const (
	DefaultCreateResourceImageFee   = 0e9                   // 0 POINT or 00000000000 upoint
	DefaultCreateResourceJSONFee    = 0e9                   // 0 POINT or 0000000000 upoint
	DefaultCreateResourceDefaultFee = 0e9                    // 0 POINT or 0000000000 upoint
	DefaultBurnFactor               = "0.000000000000000000" // 0.0 or 0%
)

// DefaultGenesis returns the default `resource` genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Resources: []*ResourceWithMetadata{},
		FeeParams: DefaultFeeParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.ValidateNoDuplicates(); err != nil {
		return err
	}

	if err := gs.FeeParams.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

func (gs GenesisState) ValidateNoDuplicates() error {
	// Group resources by collection
	resourcesByCollection := make(map[string][]*ResourceWithMetadata)

	for _, resource := range gs.Resources {
		existing := resourcesByCollection[resource.Metadata.CollectionId]
		resourcesByCollection[resource.Metadata.CollectionId] = append(existing, resource)
	}

	// Check that there are no collisions within each collection
	for _, resources := range resourcesByCollection {
		resourceIDMap := make(map[string]bool)

		for _, resource := range resources {
			if _, ok := resourceIDMap[resource.Metadata.Id]; ok {
				return fmt.Errorf("duplicated id for resource within the same collection. collection: %s, id: %s", resource.Metadata.CollectionId, resource.Metadata.Id)
			}

			resourceIDMap[resource.Metadata.Id] = true
		}
	}

	return nil
}
