#!/bin/bash


sed -i.bak -e "s|^address *= \"0\.0\.0\.0:9092\"$|address = \"0\.0\.0\.0:9095\"|" ./c1/node2/config/app.toml
sed -i -e "s|^address *= \"0\.0\.0\.0:9091\"$|address = \"0\.0\.0\.0:9096\"|" ./c1/node2/config/app.toml
sed -i -e "s|laddr = \"tcp:\/\/127.0.0.1:10002\"|laddr = \"tcp://127.0.0.1:10005\"|" ./c1/node2/config/config.toml
sed -i -e "s|^pprof_laddr *= \"localhost:6062\"|pprof_laddr = \"localhost:6065\"|" ./c1/node2/config/config.toml
sed -i -e "s|^laddr *= \"tcp:\/\/0\.0\.0\.0:20002\"|laddr = \"tcp:\/\/0\.0\.0\.0:20005\"|" ./c1/node2/config/config.toml

chain_cmd=$(which interchanged)

id=$(interchanged  tendermint show-node-id --home ./c1/node1)
echo $id
${chain_cmd} --home ./c1/node2 start --p2p.seeds ${id}@localhost:20002
