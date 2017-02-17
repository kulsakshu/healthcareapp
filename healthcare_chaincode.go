package main

import (
	"errors"
	"fmt"
    "encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


type HealthCareChaincode struct {
}

type eRewardPoint struct{
	Points string 					
	User   string
}


func (t *HealthCareChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    
    var err error
    var username , points string 
    if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

    //initialize
    username = args[0]
    points = "0"


    ePoints := eRewardPoint{}

    ePoints.Points = points
    ePoints.User = username

    jsonStr := `{"Points : "`+points +`","User : "`+ username +`"} `

    err = stub.PutState(username, []byte(jsonStr))
    if err != nil {
		return nil, err
	}

    return nil, nil
}


func (t *HealthCareChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

    if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}

    if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

    var UserName , jsonResp string 

    UserName = args[0]

    valAsbytes, err := stub.GetState(UserName)	

    if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + UserName + "\"}"
		return nil, errors.New(jsonResp)
	}
     ePoints := eRewardPoint{}
    json.Unmarshal(valAsbytes, &ePoints)


    jsonResp = `{"Username : "`+ ePoints.User+ `", "Points : "` + ePoints.Points +`"}`
    fmt.Printf("Query Response:%s\n", jsonResp)
	return valAsbytes, nil	
}

func (t *HealthCareChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

  var err error
    var username , points string 
    if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

    //initialize
    username = args[0]
    points = args[1]

    fmt.Printf("Points :%s\n", points)
        ePoints := eRewardPoint{}

    ePoints.Points = points
    ePoints.User = username

    jsonStr := `{"Points : "`+points +`","User : "`+ username +`"} `

    err = stub.PutState(username, []byte(jsonStr))
    if err != nil {
		return nil, err
	}

    return nil, nil
}


func main() {
	err := shim.Start(new(HealthCareChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
