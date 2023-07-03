#!/bin/bash

read -p "chain_command chain_name base_dir: " cmd chain_name basedir
if [[ -z "$cmd" ]] || [[ -z "$chain_name" ]] || [[ -z "$basedir"  ]];then
    echo "you must input 3 argument"
    exit 0
fi



mkdir $basedir
cd $basedir
mkdir node1
mkdir node2

$cmd init node1 --chain-id ${chain_name} --home ./node1
$cmd keys add ${basedir}-node1-account --home ./node1 > mnemonic_phrase1.txt
$cmd init node2 --chain-id ${chain_name} --home ./node2
$cmd keys add ${basedir}-node2-account --home ./node2 > mnemonic_phrase2.txt
sed -i -e "s|\S*stake\"|\"aarch\"|" node1/config/genesis.json
result1=$(${cmd} keys show ${basedir}-node1-account -a --home ./node1)
result2=$(${cmd} keys show ${basedir}-node2-account -a --home ./node2)
$cmd add-genesis-account $result1 100000000000000aarch --home ./node1
$cmd add-genesis-account $result2 100000000000000aarch --home ./node1
$cmd gentx ${basedir}-node1-account 1000000000aarch --chain-id ${chain_name} --home ./node1
$cmd collect-gentxs --home ./node1
sed -i.bak -e "s|^address *= \"0\.0\.0\.0:9090\"$|address = \"0\.0\.0\.0:9095\"|" node2/config/app.toml
sed -i -e "s|^address *= \"0\.0\.0\.0:9091\"$|address = \"0\.0\.0\.0:9096\"|" node2/config/app.toml
sed -i -e "s|laddr = \"tcp:\/\/127.0.0.1:26657\"|laddr = \"tcp://127.0.0.1:10005\"|" node2/config/config.toml
sed -i -e "s|^pprof_laddr *= \"localhost:6060\"|pprof_laddr = \"localhost:6065\"|" node2/config/config.toml
sed -i -e "s|^laddr *= \"tcp:\/\/0\.0\.0\.0:26656\"|laddr = \"tcp:\/\/0\.0\.0\.0:20005\"|" node2/config/config.toml
cp node1/config/genesis.json node2/config/genesis.json
$cmd start --home ./node1
