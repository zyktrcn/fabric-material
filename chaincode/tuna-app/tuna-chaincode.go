
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)


type SmartContract struct {
}

type IPv6 struct {
    AS_number       string   `json:"AS_number"`
    Assign_by       string   `json:"Assign_by"`
    Assign_to       string   `json:"Assign_to"`
    IPv6_prefix    string   `json:"IPv6_prefix"`
    Advertisement   string   `json:"Advertisement"`
    Timestamp   string   `json:"Timestamp"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()

	if function == "queryIPv6" {
		return s.queryIPv6(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "allocateIPv6" {
		return s.allocateIPv6(APIstub, args)
	} else if function == "queryAllAllocation" {
		return s.queryAllAllocation(APIstub, args)
	} else if function == "changePrefixHolder" {
		return s.changePrefixHolder(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

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

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	ips := []IPv6{
		IPv6{IPv6_prefix: "testAddress1", AS_number: "ASNUMBER1", Assign_by: "admin", Assign_to: "testUser1", Advertisement: "empty", Timestamp: "1543803219165"},
		IPv6{IPv6_prefix: "testAddress2", AS_number: "ASNUMBER2", Assign_by: "admin", Assign_to: "testUser2", Advertisement: "empty", Timestamp: "1543803219165"},
		IPv6{IPv6_prefix: "testAddress3", AS_number: "ASNUMBER3", Assign_by: "admin", Assign_to: "testUser3", Advertisement: "empty", Timestamp: "1543803219165"},
		IPv6{IPv6_prefix: "testAddress4", AS_number: "ASNUMBER4", Assign_by: "admin", Assign_to: "testUser4", Advertisement: "empty", Timestamp: "1543803219165"},
		IPv6{IPv6_prefix: "testAddress5", AS_number: "ASNUMBER5", Assign_by: "admin", Assign_to: "testUser5", Advertisement: "empty", Timestamp: "1543803219165"},
		IPv6{IPv6_prefix: "testAddress6", AS_number: "ASNUMBER6", Assign_by: "admin", Assign_to: "testUser6", Advertisement: "empty", Timestamp: "1543803219165"},
		IPv6{IPv6_prefix: "testAddress7", AS_number: "ASNUMBER7", Assign_by: "admin", Assign_to: "testUser7", Advertisement: "empty", Timestamp: "1543803219165"},
		IPv6{IPv6_prefix: "testAddress8", AS_number: "ASNUMBER8", Assign_by: "admin", Assign_to: "testUser8", Advertisement: "empty", Timestamp: "1543803219165"},
		IPv6{IPv6_prefix: "testAddress9", AS_number: "ASNUMBER9", Assign_by: "admin", Assign_to: "testUser9", Advertisement: "empty", Timestamp: "1543803219165"},
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

func (s *SmartContract) allocateIPv6(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	ipAsBytes, _ := APIstub.GetState(args[0])
	if ipAsBytes != nil {
		return shim.Error("Could not allocate this IPv6 prefix")
	}

	var allocation = IPv6{IPv6_prefix: args[0], AS_number: args[1], Assign_by: args[2], Assign_to: args[3], Advertisement: args[4], Timestamp: args[5]}

	ipAsBytes, _ = json.Marshal(allocation)
	err := APIstub.PutState(args[0], ipAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to allocate IPv6 prefix: %s", args[0]))
	}

	return shim.Success(nil)
}


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


	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")

		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllIPv6:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


func (s *SmartContract) changePrefixHolder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	prefixAsBytes, _ := APIstub.GetState(args[0])
	if prefixAsBytes == nil {
		return shim.Error("Could not find prefix")
	}
	prefix := IPv6{}

	json.Unmarshal(prefixAsBytes, &prefix)

	prefix.Assign_to = args[1]
	prefix.Timestamp = args[2]

	prefixAsBytes, _ = json.Marshal(prefix)
	err := APIstub.PutState(args[0], prefixAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change prefix holder: %s", args[0]))
	}

	return shim.Success(nil)
}



func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
