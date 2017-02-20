package main

import (
	"errors"
	"fmt"
    "strconv"
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
    if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

    //initialize
    username = args[0]
    points = args[1]


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

    if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

    //initialize
    //username = args[0]
   // points = args[1]

 
    //ePoints := eRewardPoint{}

    //ePoints.Points = points
   // ePoints.User = username

    //jsonStr := `{"Points : "`+points +`","User : "`+ username +`"} `

   // err = stub.PutState(username, []byte(jsonStr))
   // if err != nil {
//		return nil, err
//	}

    if function == "assign" {
		// Assign ownership
		return t.assign(stub, args)
	} else if function == "redeem" {
		// Transfer ownership
		return t.redeem(stub, args)
	}

    return nil, nil
}


func (t *HealthCareChaincode) assign(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var err error
    var fromUser , toUser, points string
    var jsonResp string

    var userPoints , assignPoint int 

    fromUser = args[0] 
    toUser = args[1]   
    points = args[2]

    fmt.Printf("fromUser = %s, toUser = %s , points = %s\n", fromUser, toUser,points)
    valAsbytes, err := stub.GetState(toUser)	

    if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + toUser + "\"}"
		return nil, errors.New(jsonResp)
	}

    ePoints := eRewardPoint{}
    json.Unmarshal(valAsbytes, &ePoints)

    return  nil, errors.New("Points :  "+ string(valAsbytes))

    assignPoint , err = strconv.Atoi(points)
    if err != nil {
		return nil, errors.New("Invalid points, expecting a integer value for assign points")
	}

    userPoints , err   = strconv.Atoi(ePoints.Points)


    if err != nil {
		return nil, errors.New("Invalid points, expecting a integer value for user points ")
	}
    

    userPoints = userPoints + assignPoint

    ePoints.Points = strconv.Itoa(userPoints)

    jsonStr := `{"Points : "`+ePoints.Points +`","User : "`+ toUser +`"} `

    err = stub.PutState(toUser, []byte(jsonStr))

    if err != nil {
		return nil, err
	}

    val, err := stub.GetState(toUser)	
    if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + toUser + "\"}"
		return nil, errors.New(jsonResp)
	}
    return val, nil 
}

func (t *HealthCareChaincode) redeem(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
   var err error
    var fromUser , toUser, points string
    var jsonResp string

    var userPoints , assignPoint int 

    fromUser = args[0] 
    toUser = args[1]   
    points = args[2]

  fmt.Printf("fromUser = %s, toUser = %s , points = %s\n", fromUser, toUser,points)
    valAsbytes, err := stub.GetState(toUser)	

    if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + toUser + "\"}"
		return nil, errors.New(jsonResp)
	}

    ePoints := eRewardPoint{}
    json.Unmarshal(valAsbytes, &ePoints)

    userPoints , err   = strconv.Atoi(string(ePoints.Points))
    if err != nil {
		return nil, errors.New("Invalid points, expecting a integer value for user points")
	}
    assignPoint , err = strconv.Atoi(string(points))
    if err != nil {
		return nil, errors.New("Invalid points, expecting a integer value for assign points")
	}

    userPoints = userPoints - assignPoint

     ePoints.Points = strconv.Itoa(userPoints)

    jsonStr := `{"Points : "`+ePoints.Points +`","User : "`+ toUser +`"} `

    err = stub.PutState(toUser, []byte(jsonStr))

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
