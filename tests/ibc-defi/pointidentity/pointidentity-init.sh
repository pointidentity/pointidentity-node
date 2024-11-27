#!/bin/bash
# shellcheck disable=SC2086

set -euox pipefail

# sed in MacOS requires extra argument
if [[ "$OSTYPE" == "darwin"* ]]; then
  SED_EXT='.orig'
else 
  SED_EXT=''
fi

CHAIN_ID="pointidentity"

# Node
pointidentity-noded init node0 --chain-id "$CHAIN_ID"
NODE_0_VAL_PUBKEY=$(pointidentity-noded tendermint show-validator)

# User
pointidentity-noded keys add pointidentity-user --keyring-backend test

# Config
sed -i $SED_EXT 's|minimum-gas-prices = ""|minimum-gas-prices = "50upoint"|g' "$HOME/.node/config/app.toml"

# shellcheck disable=SC2086
sed -i $SED_EXT 's|laddr = "tcp://127.0.0.1:26647"|laddr = "tcp://0.0.0.0:26647"|g' "$HOME/.node/config/config.toml"
sed -i $SED_EXT 's|address = "localhost:9090"|address = "0.0.0.0:9090"|g' "$HOME/.node/config/app.toml"

# Genesis
GENESIS="$HOME/.node/config/genesis.json"
sed -i $SED_EXT 's/"stake"/"upoint"/' "$GENESIS"

pointidentity-noded genesis add-genesis-account pointidentity-user 1000000000000000000upoint --keyring-backend test
pointidentity-noded genesis gentx pointidentity-user 10000000000000000upoint --chain-id $CHAIN_ID --pubkey "$NODE_0_VAL_PUBKEY" --keyring-backend test

pointidentity-noded genesis collect-gentxs
pointidentity-noded genesis validate-genesis
