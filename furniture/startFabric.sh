#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

starttime=$(date +%s)

# launch network; create channel and join peer to channel
cd ../network
./start.sh

# Now launch the CLI container in order to install, instantiate chaincode
# and prime the ledger with our 10 cars
docker-compose -f ./docker-compose.yml up -d cli

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.furniture.com/users/Admin@org1.furniture.com/msp" cli peer chaincode install -n furniture -v 1.0 -p github.com/furniture
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.furniture.com/users/Admin@org1.furniture.com/msp" cli peer chaincode instantiate -o orderer.furniture.com:7050 -C furniturechannel -n furniture -v 1.0 -c '{"Args":[""]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
sleep 10
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.furniture.com/users/Admin@org1.furniture.com/msp" cli peer chaincode invoke -o orderer.furniture.com:7050 -C furniturechannel -n furniture -c '{"function":"initLedger","Args":[""]}'

printf "\nTotal setup execution time : $(($(date +%s) - starttime)) secs ...\n\n\n"
printf "Start by installing required packages run 'npm install'\n"
printf "Then run 'node enrollAdmin.js', then 'node registerUser'\n\n"
printf "The 'node queryAllFurnitures.js' to query all furnitures\n\n"
printf "The 'node queryFurniture.js' to query a specific furniture\n\n"
printf "The 'node createNewFurniture.js' to create new furniture\n\n"
printf "The 'node changeOwner.js' to change owner of furniture\n\n"