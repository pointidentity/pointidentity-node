package setup

import (
	"github.com/pointidentity/pointidentity-node/x/resource"
	"github.com/pointidentity/pointidentity-node/x/resource/types"
)

func (s *TestSetup) StorePortWithGenesis() {
	genesis := types.DefaultGenesis()
	resource.InitGenesis(s.SdkCtx, s.ResourceKeeper, genesis)
}
