#!/bin/bash

# Inits node configuration and runs the node.
# e -> exit immediately, u -> treat unset variables as errors and immediately, o -> sets the exit code to the rightmost command 
set -euo pipefail

# within the container, $HOME=/home/pointidentity
POINTIDENTITY_ROOT_DIR="$HOME/.node"

# Init node config directory
if [ ! -d "${POINTIDENTITY_ROOT_DIR}/config" ]
then
    echo "Node config not found. Initializing."
    pointidentity-noded init "moniker-placeholder"
else
    echo "Node config exists. Skipping initialization."
fi

# Check if a genesis file has been passed in config
if [ -f "/genesis" ]
then
    echo "Genesis file passed. Adding/replacing current genesis file."
    cp /genesis "${POINTIDENTITY_ROOT_DIR}/config/genesis.json"
else
    echo "No genesis file passed. Skipping and retaining existing genesis."
fi

# Check if a seeds file has been passed in config
if [ -f "/seeds" ]
then
    echo "Seeds file passed. Overriding current seeds."
    cp /seeds "${POINTIDENTITY_ROOT_DIR}/config/seeds.txt"
    POINTIDENTITY_NODED_P2P_SEEDS="$(cat "${POINTIDENTITY_ROOT_DIR}/config/seeds.txt")"
    export POINTIDENTITY_NODED_P2P_SEEDS
else
    echo "No seeds file passed. Skipping and retaining existing seeds."
fi

# Check if a node_key file has been passed in config
if [ -f "/node_key" ]
then
    echo "Node key file passed. Overriding current key."
    cp /node_key "${POINTIDENTITY_ROOT_DIR}/config/node_key.json"
else
    echo "No node key file passed. Skipping and retaining existing node key."
fi

# Check if a priv_validator_key file has been passed in config
if [ -f "/priv_validator_key" ]
then
    echo "Private validator key file passed. Replacing current key."
    cp /priv_validator_key "${POINTIDENTITY_ROOT_DIR}/config/priv_validator_key.json"
else
    echo "No private validator key file passed. Skipping and retaining existing key."
fi

# Run node
pointidentity-noded start
