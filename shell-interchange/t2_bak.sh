#!/bin/bash

cd testnet

id=$(interchanged status | sed -r "s/(.*\"id\":)(.*)/\2/" | awk -F '"' '{print $2}')

chain_cmd=$(which interchanged)

nohup ${chain_cmd} --home ./node2 start --p2p.seeds ${id}@localhost:26656 &
sleep 10

chain_id=$(cat ./node1/config/genesis.json | grep 'chain_id' | awk -F '"' '{print $4}')

interchanged tx staking create-validator --from node2-account --amount 10000000000000aarch --min-self-delegation 100 --commission-rate 0.01 --commission-max-rate 0.1 --commission-max-change-rate 0.1 --pubkey "$(interchanged tendermint show-validator --home ./node2)" --chain-id ${chain_id} --home ./node2

sleep 10
interchanged --home ./node1 q tendermint-validator-set

