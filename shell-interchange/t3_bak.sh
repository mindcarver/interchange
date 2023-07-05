#!/bin/bash


chain_cmd=$(which interchanged)

id=$(interchanged status | sed -r "s/(.*\"id\":)(.*)/\2/" | awk -F '"' '{print $2}')
echo $id
${chain_cmd} --home ./node2 start --p2p.seeds ${id}@localhost:26656
