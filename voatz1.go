/*
A product for Voatz company of NiMSiM LLC
Implemented from ASF Marbles
Sends votes from voters to candidates
Written by K
*/

// *********** PART I ***********

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"strings"

	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

// Chaincode implementation of electronic voting
type VoatzCC struct {
}

var voteIndexStr = "_voteindex"			//stores list of all known votes in transmission
var openIntentsStr = "_openintents"		//stores all intents to vote

type Voter struct{
}

type Candidate struct{
}

type User struct{
	Voter
	Candidate
	UserID int 'json:"userid"'			//tags each user
}

type Vote struct{
	User
	Tag string 'json:"tag"'			//tags the vote
	Election string 'json:"election"'		//type of election; e.g., presidential elections
	Choice string 'json:"choice"'			//voter's choice (the vote); e.g., Obama
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Main
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func main(){
	err := shim.Start(new(VoatzCC))

	if err != nil {
		fmt.Printf("Error initiating vote transmission", err)
	}
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Init :  Reset all
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** 
func (t *VoatzCC) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var Aval int
	var err error
     
	if len(args) != 1 {
		return nil, errors.New("Error: # of arguments. Expecting: 1")
	}

	//Initialize the CC
	Aval, err = strconv.Atoi(args[0])

	if err != nil {
		return nil, errors.New("Error. Expecting: integer")
	}

	//Record state
	err = stub.PutState("testVar", []byte(strconv.Itoa(Aval)))		//making test variable to test the network

	if err != nil {
		return nil, err
	}

	var empty []string
	jsonB, _ := json.Marshal(empty)
	err = stub.PutState(voteIndexStr, jsonB)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Run
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func (t *VoatzCC) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Voting in progress via " + function)
     
	//Handling functions:
	if function == "init" {		//initialize CC state, or reset
		return t.init(stub, args)
	} else if function == "delete" {		//delete argument from its state
		return t.Delete(stub, args)
	} else if function == "write" {		//write argument to CC state
		return t.Write(stub, args)
	} else if function == "init_vote" {		//initiate new vote
		return t.init_vote(stub, args)
	} else if function == "set_user" {		//change owner of a vote
		return t.set_user(stub, args)
	}

	fmt.Println("Can't find function " + function + " during run")

	return nil, error.New("Error: Invoked unknown function")
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Quest :  Initiate quest
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func (t *VoatzCC) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Quest in progress via " + function)

	//Dealing functions:
	if function == "read" {		//read any variable
		return t.read(stub, args)
	}

	fmt.Println("Can't find " + function + " during quest")

	return nil, errors.New("Error: Unknown function in quest")
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Read ;  var from CC state
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func (t *VoatzCC) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var tag string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Error: # of arguments. Expecting: 1; enter tag of vote to inquire")
	}

	tag = args[0]
	valB, err := stub.GetState(tag)		//get the tag from CC state

	if err != nil {
		return nil, errors.New("Error: Failed to get state for " + tag)
	}

	return valB, nil		//send the tag forward
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Delete ;  from state
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func (t *VoatzCC) Delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1{
		return nil, error.New("Error: # of arguments. Expecting: 1")
	}

	name := args[0]
	err := stub.DelState(tag)		//remove the tag from CC state

	if err != nil {
		return nil, errors.New("Error: Can't delete state")
	}

	//get vote index
	votesB, err := stub.GetState(voteIndexStr)

	if err != nil {
		return nil, errors.New("Error: Can't get vote index")
	}

	var voteIndex []string
	json.Unmarshal(votesB, &voteIndex)		//de-string; Json.parse()

	//retrieve vote from index
	for i, val := range voteIndex {
		fmt.Println(strconv.Itoa(i) + ". Looking at " + val + " for " + tag)

		if val == tag {			//find the vote to remove
			fmt.Println("Found vote")
			voteIndex = append(voteIndex[:i], voteIndex[i+1:]…)		//remove the vote

			for x := range voteIndex {		//debug
				fmt.Println(string(x) + " - " + voteIndex[x])
			}

			break
		}
	}
	jsonB, _ := json.Marshal(voteIndex)
	err = stub.PutState(userIndexStr, jsonB)
	return nil, nil
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Write ;  var to CC state
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func (t *VoatzCC) Write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var tag, val string
	var err error
	fmt.Println("Write in progress")

	if len(args) != 2 {
		return nil, errors.New("Error: # of arguments. Expecting: 2; var tag and value to set")
	}

	tag = args[0]		//update tag
	val = args[1]
	err = stub.PutState(tag, []byte(val))			//write the tag into CC state

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Initiate Vote :  Create new Vote and store to CC state
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func (t *VoatzCC) init_vote(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	// 0 1 2 3
	// userID, tag, election, choice
	// e.g., "4850", "FromJaniceBlack", "presidential", "BarrackObama" for vote sent from voter
	// e.g., "9374", "FromJaniceBlackToBarrackObama", "presidential", "BarrackObama" for vote received by candidate

	if len(args) != 4{
		return nil, errors.New("Error: # of arguments. Expecting: 4")
	}

	fmt.Println("Initiating vote..")
	
	if len(args[0]) <= 0 {
		return nil, errors.New("Error: 1st argument. Expecting: User ID, as a string of numbers")
	}

	userID, err := strconv.Atoi(args[0])

	if err != nil {
		return nil, errors.New("Error: 1st argument. Expecting: User ID, as a numeric string")
	}

	if len(args[1]) <= 0 {
		return nil, errors.New("Error: 2nd argument. Expecting: Vote tag, as a string of letters")
	}

	tag := strings.ToLower(args[1])
	
	if len(args[2]) <= 0 {
		return nil, errors.New("Error: 3rd argument. Expecting: Election, as a string of letters")
	}

	election := strings.ToLower(args[2])

	if len(args[3]) <= 0 {
		return nil, errors.New("Error: 4th argument. Expecting: Choice, as a string of letters")
	}

	choice := strings.ToLower(args[3])

	//Manually building json string for a vote:
	str := '{"userID": "' + strconv.Itoa(userID) + '", "tag": "' + tag + '", "election": "' + election + '", "choice": "' + choice + '"}'
	err = stub.PutState(args[1], []byte(str))		//store vote with tag as key
	
	if err != nil {
		return nil, err
	}

	//get the vote index
	votesB, err := stub.GetState(voteIndexStr)

	if err != nil {
		return nil, errors.New(Error: Can't get vote index)
	}

	var voteIndex []string
	json.Unmarshal(votesB, &voteIndex)		//de-string; Json.parse()

	//append
	voteIndex = append(voteIndex, args[1])
	fmt.Println("Vote index is ", voteIndex)		//add vote tag to index list
	jsonB, _ := json.Marshal(voteIndex)
	err = stub.PutState(voteIndexStr, jsonB)		//store vote tag

	fmt.Println(Completing vote initialization..)

	return nil, nil
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Set User Permissions on Vote; identical for voter & for candidate
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func (t *VoatzCC) set_user(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	// 0  1
	// vote tag, userID
	// "FromJaniceBlack", "4850"
	if len(args) < 2 {
		return nil, error.New(Error: # of arguments. Expecting: 2)
	}

	fmt.Println("Initiating new user..")
	fmt.Println(args[0] + " : " + args[1])
	votesB, err := stub.GetState(args[0])

	if err != nil {
		return nil, errors.New("Error: Setting user")
	}
	
	line := Vote{}
	json.Unmarshal(votesB, &line)		//de-string; Json.parse()
	line.User = args[1]			//change user

	jsonB, _ := json.Marshal(line)
	err = stub.PutState(args[0]. jsonB)	//rewrite user with tag as key

	if err != nil {
		return nil, err
	}

	fmt.Println("Completing user initialization..")

	return nil, nil
}









