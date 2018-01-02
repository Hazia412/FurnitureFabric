/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the furniture structure, with 4 properties.  Structure tags are used by encoding/json library
type Furniture struct {
	Type   string `json:"type"`
	Size  string `json:"size"`
	Material string `json:"material"`
	Owner  string `json:"owner"`
}

/*
 * The Init method is called when the Smart Contract "furniture" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "furniture"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryFurniture" {
		return s.queryFurniture(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createFurniture" {
		return s.createFurniture(APIstub, args)
	} else if function == "queryAllFurnitures" {
		return s.queryAllFurnitures(APIstub)
	} else if function == "changeFurnitureOwner" {
		return s.changeFurnitureOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryFurniture(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	furnitureAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(furnitureAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	furnitures := []Furniture{
		Furniture{Type:"chair", Size:"medium" ,Material:"plastic", Owner:"HaziaF"},
		Furniture{Type:"bed" ,Size:"king", Material:"wood" ,Owner:"Prajyot"},
		Furniture{Type:"table" ,Size:"large-round", Material:"glass and wood", Owner:"Vibhuti"},
		Furniture{Type:"couch", Size:"L-shaped", Material:"leather", Owner:"Prabhav"},
		Furniture{Type:"cupboard", Size:"large" ,Material:"ply-wood", Owner:"Madhushree"},
		Furniture{Type:"crib", Size:"medium" ,Material:"metal" ,Owner:"Lysanne"},
	}

	i := 0
	for i < len(furnitures) {
		fmt.Println("i is ", i)
		furnitureAsBytes, _ := json.Marshal(furnitures[i])
		APIstub.PutState("FURNITURE"+strconv.Itoa(i), furnitureAsBytes)
		fmt.Println("Added", furnitures[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createFurniture(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var furniture = Furniture{Type: args[1], Size: args[2], Material: args[3], Owner: args[4]}

	furnitureAsBytes, _ := json.Marshal(furniture)
	APIstub.PutState(args[0], furnitureAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllFurnitures(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "FURNITURE0"
	endKey := "FURNITURE999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString("abc")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllFurnitures:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeFurnitureOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	furnitureAsBytes, _ := APIstub.GetState(args[0])
	furniture := Furniture{}

	json.Unmarshal(furnitureAsBytes, &furniture)
	furniture.Owner = args[1]

	furnitureAsBytes, _ = json.Marshal(furniture)
	APIstub.PutState(args[0], furnitureAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
