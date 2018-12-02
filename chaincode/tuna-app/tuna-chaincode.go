// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario
 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
 */

package main

/* Imports  
* 4 utility libraries for handling bytes, reading and writing JSON, 
formatting, and string manipulation  
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts  
*/ 
import (
	"bytes"
	"encoding/json"
	"fmt"
	// "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Tuna structure, with 4 properties.  
Structure tags are used by encoding/json library
*/
type IPv6 struct {
    AS_number       string   `json:"AS_number"` 
    Assign_by       string   `json:"Assign_by"` 
    Assign_to       string   `json:"Assign_to"`
    IPv6_prefix    string   `json:"IPv6_prefix"`
    Advertisement   string   `json:"Advertisement"`
}

/*
 * The Init method *
 called when the Smart Contract "tuna-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function 
 -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "tuna-chaincode"
 The app also specifies the specific smart contract function to call with args
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryIPv6" {
		return s.queryIPv6(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "allocateIPv6" {
		return s.allocateIPv6(APIstub, args)
	} else if function == "queryAllAllocation" {
		return s.queryAllAllocation(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryIPv6 method *
Used to view the records of one particular tuna
It takes one argument -- the key for the tuna in question
 */
func (s *SmartContract) queryIPv6(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ipAsBytes, _ := APIstub.GetState(args[0])
	if ipAsBytes == nil {
		return shim.Error("Could not query IPv6 prefix")
	}
	return shim.Success(ipAsBytes)
}

/*
 * The initLedger method *
Will add test data (10 tuna catches)to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	ips := []IPv6{
		IPv6{IPv6_prefix: "testAddress1", AS_number: "ASNUMBER1", Assign_by: "admin", Assign_to: "testUser1", Advertisement: "empty"},
		IPv6{IPv6_prefix: "testAddress2", AS_number: "ASNUMBER2", Assign_by: "admin", Assign_to: "testUser2", Advertisement: "empty"},
		IPv6{IPv6_prefix: "testAddress3", AS_number: "ASNUMBER3", Assign_by: "admin", Assign_to: "testUser3", Advertisement: "empty"},
		IPv6{IPv6_prefix: "testAddress4", AS_number: "ASNUMBER4", Assign_by: "admin", Assign_to: "testUser4", Advertisement: "empty"},
		IPv6{IPv6_prefix: "testAddress5", AS_number: "ASNUMBER5", Assign_by: "admin", Assign_to: "testUser5", Advertisement: "empty"},
		IPv6{IPv6_prefix: "testAddress6", AS_number: "ASNUMBER6", Assign_by: "admin", Assign_to: "testUser6", Advertisement: "empty"},
		IPv6{IPv6_prefix: "testAddress7", AS_number: "ASNUMBER7", Assign_by: "admin", Assign_to: "testUser7", Advertisement: "empty"},
		IPv6{IPv6_prefix: "testAddress8", AS_number: "ASNUMBER8", Assign_by: "admin", Assign_to: "testUser8", Advertisement: "empty"},
		IPv6{IPv6_prefix: "testAddress9", AS_number: "ASNUMBER9", Assign_by: "admin", Assign_to: "testUser9", Advertisement: "empty"},
	}

	i := 0
	for i < len(ips) {
		fmt.Println("i is ", i)
		ipAsBytes, _ := json.Marshal(ips[i])
		APIstub.PutState(ips[i].IPv6_prefix, ipAsBytes)
		fmt.Println("Added", ips[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The allocateIPv6 method *
Fisherman like Sarah would use to record each of her tuna catches. 
This method takes in five arguments (attributes to be saved in the ledger). 
 */
func (s *SmartContract) allocateIPv6(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	ipAsBytes, _ := APIstub.GetState(args[0])
	if ipAsBytes != nil {
		return shim.Error("Could not allocate this IPv6 prefix")
	}

	var allocation = IPv6{IPv6_prefix: args[0], AS_number: args[1], Assign_by: args[2], Assign_to: args[3], Advertisement: args[4]}

	ipAsBytes, _ = json.Marshal(allocation)
	err := APIstub.PutState(args[0], ipAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to allocate IPv6 prefix: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllAllocation method *
allows for assessing all the records added to the ledger(all tuna catches)
This method does not take any arguments. Returns JSON string containing results. 
 */
func (s *SmartContract) queryAllAllocation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
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

	fmt.Printf("- queryAllIPv6:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


/*
 * The changeTunaHolder method *
The data in the world state can be updated with who has possession. 
This function takes in 2 arguments, tuna id and new holder name. 
 */
func (s *SmartContract) changePrefixHolder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	prefixAsBytes, _ := APIstub.GetState(args[0])
	if prefixAsBytes == nil {
		return shim.Error("Could not find prefix")
	}
	prefix := IPv6{}

	json.Unmarshal(prefixAsBytes, &prefix)
	// Normally check that the specified argument is a valid holder of prefix
	// we are skipping this check for this example
	prefix.Assign_to = args[1]

	prefixAsBytes, _ = json.Marshal(prefix)
	err := APIstub.PutState(args[0], prefixAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change prefix holder: %s", args[0]))
	}

	return shim.Success(nil)
}


/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}