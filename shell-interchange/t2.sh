#!/bin/bash


id=$(interchanged status -n tcp://localhost:10002 | sed -r "s/(.*\"id\":)(.*)/\2/" | awk -F '"' '{print $2}')

chain_cmd=$(which interchanged)

nohup ${chain_cmd} --home ./c1/node2 start --p2p.seeds ${id}@localhost:20002 &
sleep 5


