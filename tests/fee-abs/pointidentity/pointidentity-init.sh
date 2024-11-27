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

CONFIG_TOML="$HOME/.node/config/config.toml"

# Config
sed -i $SED_EXT 's|minimum-gas-prices = ""|minimum-gas-prices = "50upoint"|g' "$HOME/.node/config/app.toml"
sed -i $SED_EXT 's|addr_book_strict = true|addr_book_strict = false|g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/timeout_propose = "3s"/timeout_propose = "500ms"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/timeout_prevote = "1s"/timeout_prevote = "500ms"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/timeout_precommit = "1s"/timeout_precommit = "500ms"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/timeout_commit = "5s"/timeout_commit = "2s"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/log_level = "info"/log_level = "debug"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/"voting_period": "172800s"/"voting_period": "12s"/' "$HOME/.node/config/genesis.json"

# shellcheck disable=SC2086
sed -i $SED_EXT 's|laddr = "tcp://127.0.0.1:26647"|laddr = "tcp://0.0.0.0:26647"|g' "$HOME/.node/config/config.toml"
sed -i $SED_EXT 's|address = "localhost:9090"|address = "0.0.0.0:9090"|g' "$HOME/.node/config/app.toml"
sed -i $SED_EXT 's|log_level = "error"|log_level = "info"|g' "$HOME/.node/config/config.toml"

# Genesis
GENESIS="$HOME/.node/config/genesis.json"
sed -i $SED_EXT 's/"stake"/"upoint"/' "$GENESIS"

pointidentity-noded genesis add-genesis-account pointidentity-user 1000000000000000000upoint --keyring-backend test
pointidentity-noded genesis gentx pointidentity-user 10000000000000000npoint--chain-id $CHAIN_ID --pubkey "$NODE_0_VAL_PUBKEY" --keyring-backend test

pointidentity-noded genesis collect-gentxs
pointidentity-noded genesis validate-genesis
