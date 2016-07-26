/*
A product for Voatz company of NiMSiM LLC
Implemented from ASF Marbles
Sends votes from voters to candidates
Written by K
*/

// *********** PART II ***********

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"time"
	"strings"
	
	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

// Chaincode implementation of electronic voting
type VoatzCC struct {
}

var voteIndexStr = "_voteindex"			//stores list of all known votes in transmission
var openIntentsStr = "_openintents"			//stores all intents to vote

type Voter struct{
}

type Candidate struct{
}

type User struct{
	Voter
	Candidate
	UserID int `json:"userid"`			//tags each user
}

type Vote struct{
	User
	Tag string `json:"tag"`			//tags the vote
	Election string `json:"election"`		//type of election; e.g., presidential elections
	Choice string `json:"choice"`			//voter's choice (the vote); e.g., Obama
}

type Description struct {
	Choice string `json:"choice"`			// argument to be exchanged; send voter's choice to candidate
}

type AnOpenIntent struct{
	User
	Timestamp int64 `json:"timestamp"`
	Want Description  `json:"want"`
	Willing []Description `json:"willing"`
}

type AllIntents struct {
	OpenIntents []AnOpenIntent `json:"open_trades"`
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Main
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func main(){
	err := shim.Start(new(VoatzCC))

	if err != nil {
		fmt.Printf("Error initiating vote transmission: %s", err)
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

	var intents AllIntents
	jsonB, _ = json.Marshal(intents)
	err = stub.PutState(openIntentsStr, jsonB)

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
	return t.Invoke(stub, function, args)
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Invoke :  Entry point for invocations
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func (t *VoatzCC) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Invocation in progress via " + function)

  	//Handling functions:
	if function == "init" {				//initialize CC state, or reset
		return t.init(stub, args)
	} else if function == "delete" {			//delete argument from its state
		line, err := t.Delete(stub, args)
		cleanIntents(stub)
		return line, err
	} else if function == "write" {				//write argument to CC state
		return t.Write(stub, args)
	} else if function == "init_vote" {			//initiate new vote
		return t.init_vote(stub, args)
	} else if function == "set_user" {			//change owner of a vote
		line, err := t.set_user(stub, args)
		cleanIntents(stub)
		return line, err
	} else if function == "vote_intent" {
		return t.vote_intent(stub, args)
	} else if function == "transmit" {		//change owner of a vote; send from voter to candidate
		line, err := t.transmit(stub, args)
		cleanIntents(stub)
		return line, err
	} else if function == "remove_intent" {
		return t.remove_intent(stub, args)
	}

	fmt.Println("Can't find function " + function + " while invoking")

	return nil, errors.New("Error: Invoked unknown function")
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
		return nil, errors.New("Error: # of arguments. Expecting: 1")
	}

	tag := args[0]
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
			voteIndex = append(voteIndex[:i], voteIndex[i+1:]...)		//remove the vote

			for x := range voteIndex {		//debug
				fmt.Println(string(x) + " - " + voteIndex[x])
			}

			break
		}
	}
	jsonB, _ := json.Marshal(voteIndex)
	err = stub.PutState(voteIndexStr, jsonB)
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

	//check if vote exists already
	tagB, err := stub.GetState(tag)
	
	if err != nil {
		return nil, errors.New("Error: Can't get vote tag")
	}

	line := Vote {}
	json.Unmarshal(tagB, &line)

	if line.Tag == tag {
		fmt.Println("Vote already active: " + tag)
		fmt.Println(line);
	
		return nil, errors.New("Vote already exists")
	}

	//Manually building json string for a vote:
	str := `{"userID": "` + strconv.Itoa(userID) + `", "tag": "` + tag + `", "election": "` + election + `", "choice": "` + choice + `"}`
	err = stub.PutState(args[1], []byte(str))		//store vote with tag as key
	
	if err != nil {
		return nil, err
	}

	//get the vote index
	votesB, err := stub.GetState(voteIndexStr)

	if err != nil {
		return nil, errors.New("Error: Can't get vote index")
	}

	var voteIndex []string
	json.Unmarshal(votesB, &voteIndex)		//de-string; Json.parse()

	//append
	voteIndex = append(voteIndex, args[1])
	fmt.Println("Vote index is ", voteIndex)		//add vote tag to index list
	jsonB, _ := json.Marshal(voteIndex)
	err = stub.PutState(voteIndexStr, jsonB)		//store vote tag

	fmt.Println("Completing vote initialization..")

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
		return nil, errors.New("Error: # of arguments. Expecting: 2")
	}

	fmt.Println("Initiating new user..")
	fmt.Println(args[0] + " : " + args[1])
	votesB, err := stub.GetState(args[0])

	if err != nil {
		return nil, errors.New("Error: Setting user")
	}
	
	line := Vote{}
	json.Unmarshal(votesB, &line)		//de-string; Json.parse()
	line.UserID, err = strconv.Atoi(args[1])			//change user

	if err != nil {
	   return nil, errors.New("Error. Expecting: String of numbers")
	   }

	jsonB, _ := json.Marshal(line)
	err = stub.PutState(args[0], jsonB)	//rewrite user with tag as key

	if err != nil {
		return nil, err
	}

	fmt.Println("Completing user initialization..")

	return nil, nil
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Open Intent :  Create a vote intent ;  to be sent from voter to candidate
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func (t *VoatzCC) vote_intent(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error
	var will_choice string
	var intent_sent Description

	// Alternative {
	// 0 1 2 3
	// "4850", "BernieSanders", "HillaryClinton"], "DonaldTrump" 
	// [userID (of voter), 1st vote choice of voter, 2nd vote choice of voter], new vote choice of voter
	// }

	//0 1 2 3
	// ["9374", "BarrackObama", "BarrackObama"], "BarrackObama"
	// [userID (of candidate), choice; of vote present at candidate, choice; of vote willed by user (e.g., BO collects all vote intents for himself)], choice; new vote sent to candidate

	if len(args) < 3 {
		return nil, errors.New("Error: # of arguments. Expecting: 3")
	}
	
	userID, err := strconv.Atoi(args[0])

	if err != nil {
		return nil, errors.New("Error: 1st argument. Expecting: User ID; as numeric string")
	}

	open := AnOpenIntent{}
	open.UserID = userID
	open.Timestamp = makeTimestamp()
	open.Want.Choice = args[1]
	
	fmt.Println("Submitting vote intent to candidate..")
	jsonB, _ := json.Marshal(open)
	err = stub.PutState("_debug1", jsonB)

	for i := 1 ;  i < len(args) ;  i++ {
		will_choice = strings.ToLower(args[i])

		// if err != nil {
		// 	warning := "Choice; " + args[i] + " is not an alphabetical string"
		// 	fmt.Println(warning)
		// 	return nil, errors.New("Error: " + warning)
		// }
	
		intent_sent = Description{}
		intent_sent.Choice = will_choice
		fmt.Println("Sending intent: " + will_choice)
		jsonB, _ = json.Marshal(intent_sent)
		err = stub.PutState("_debug2", jsonB)

		open.Willing = append(open.Willing, intent_sent)
		fmt.Println("Appending will to open..")
		i++;
	}

	// Obtaining open intent struct:
	intentsB, err := stub.GetState(openIntentsStr)
	
	if err != nil {
		return nil, errors.New("Error: Can't get open intents")
	}

	var intents AllIntents
	json.Unmarshal(intentsB, &intents)

	intents.OpenIntents = append(intents.OpenIntents, open)
	fmt.Println("Appending open intent to the existent..")
	jsonB, _ = json.Marshal(intents)
	err = stub.PutState(openIntentsStr, jsonB)

	if err != nil {
		return nil, err
	}

	fmt.Println("Ending open intent..")

	return nil, nil
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Transmit :  Close an open intent and move ownership ;  transmit vote from voter to client
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func (t *VoatzCC) transmit(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error
	var userIDIntentStr string

	// 0 1 2 3 4
	//[data.id(timestamp how ??), data.closer.voter(userID), data.closer.tag, data.opener.candidate(userID), data.opener.choice]
	// This should transfer vote from voter to candidate since it already classifies opener vs. closer ??
	if len(args) < 5 {
		return nil, errors.New("Error: # of arguments. Expecting: 5")
	}

	fmt.Println("Transmitting vote..")
	timestamp, err := strconv.ParseInt(args[0], 10, 55)  //not sure about these numbers; 10, 55=(9*len(args))+10 ??

	if err != nil {
		return nil, errors.New("Error: 1st argument. Expecting: Numeric string")
	}

	// voterID, err := strconv.Atoi(args[1])

	// if err != nil {
	//	return nil, errors.New("Error: 2nd argument. Expecting: Numeric string")
	// }

	// candidateID, err := strconv.Atoi(args[3])

	// if err != nil {
	//	return nil, errors.New("Error: 4th argument. Expecting: Numeric string")
	// }

	if len(args[4]) <= 0 {
	   return nil, errors.New("Error: 5th argument. Expecting: Choice, as a string of letters")
	   }

	choice := strings.ToLower(args[3])

	// Obtaining open intent struct:
	intentsB, err := stub.GetState(openIntentsStr)

	if err != nil {
		return nil, errors.New("Error: Can't get open intents")
	}

	var intents AllIntents
	json.Unmarshal(intentsB, &intents)

	for i := range intents.OpenIntents {
		fmt.Println("Searching for " + strconv.FormatInt(timestamp, 10) + " at " + strconv.FormatInt(intents.OpenIntents[i].Timestamp, 10))

		if intents.OpenIntents[i].Timestamp == timestamp {
			fmt.Println("Intent located")
			votesB, err := stub.GetState(args[2])

			if err != nil {
				return nil, errors.New("Error: Locating intent")
			}

			closersVote := Vote{}
			json.Unmarshal(votesB, &closersVote)

			userIDIntentStr	= strconv.Itoa(intents.OpenIntents[i].UserID) // UserID belongs to voter

			// Checking for vote requirements
			if closersVote.Choice != intents.OpenIntents[i].Want.Choice {
				line := "Vote not as intended"
				fmt.Println(line)
				return nil, errors.New("Error: " + line)
			}

			vote, e := findIntendedVote(stub, userIDIntentStr, choice)

			if (e == nil) {
				fmt.Println("Proceeding..")
				t.set_user(stub, []string{args[2], userIDIntentStr})		// Voter instead of User ??
				t.set_user(stub, []string{vote.Tag, args[1]})
				intents.OpenIntents = append(intents.OpenIntents[:i], intents.OpenIntents[i+1:]...)	// remove vote
				jsonB, _ := json.Marshal(intents)
				err = stub.PutState(openIntentsStr, jsonB)

				if err != nil {
					return nil, err
				}
			}
		}
	}
	fmt.Println("Ending vote transmission..")
	
	return nil, nil
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Find the Intended Vote :  Search for vote this candidate owns
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func findIntendedVote(stub *shim.ChaincodeStub, user string, choice string)(v Vote, err error){
	var fail Vote
	var userIDStr string

	fmt.Println("Initiating search for intended vote..")
	fmt.Println("Searching for " + user + ", " + choice)

	// Obtaining vote index:
	votesB, err := stub.GetState(voteIndexStr)

	if err != nil {
		return fail, errors.New("Error: Can't get vote index")
	}

	var voteIndex []string
	json.Unmarshal(votesB, &voteIndex)

	for i := range voteIndex {
		votesB, err := stub.GetState(voteIndex[i])

		if err != nil {
			return fail, errors.New("Error: Can't get vote")
		}

		line := Vote{}
		json.Unmarshal(votesB, &line)
		   
		userIDStr = strconv.Itoa(line.UserID)

		if strings.ToLower(userIDStr) == strings.ToLower(user) && strings.ToLower(line.Choice) == strings.ToLower(choice) {
			fmt.Println("Found a vote:" + line.Tag)
			fmt.Println("Ending search for intended vote..")

			return line, nil
		}
	}

	fmt.Println("Error: Ending search for intended vote..")

	return fail, errors.New("Error: Couldn't find intended vote")
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Make Timestamp ;  in ms
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Remove Open Intent ; remove intent ;  finalize vote transmission
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func (t *VoatzCC) remove_intent(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	// 0
	// [data.id]
	if len(args) < 1 {
		return nil, errors.New("Error: # of arguments. Expecting: 1")
	}

	fmt.Println("Removing intent..")
	timestamp, err := strconv.ParseInt(args[0], 10, 55)

	if err != nil {
		return nil, errors.New("Error: 1st argument. Expecting: Numeric string")
	}

	// Obtaining open intent struct:
	intentsB, err := stub.GetState(openIntentsStr)
	
	if err != nil {
		return nil, errors.New("Error: Can't get open intents")
	}

	var intents AllIntents
	json.Unmarshal(intentsB, &intents)

	for i := range intents.OpenIntents {
		if intents.OpenIntents[i].Timestamp == timestamp {
			fmt.Println("Located intent")
			intents.OpenIntents = append(intents.OpenIntents[:i], intents.OpenIntents[i+1:]...)
			jsonB, _ := json.Marshal(intents)
			err = stub.PutState(openIntentsStr, jsonB)
			
			if err != nil {
				return nil, err
			}
	
			break
		}
	}

	fmt.Println("Ending intent removal..")
	
	return nil, nil
}

// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
// Clean Intents ;  remove disapproved or closed intents
// *********** *********** *********** *********** *********** *********** *********** *********** *********** *********** ***********
func cleanIntents(stub *shim.ChaincodeStub)(err error){
	var valid = false
	var userIDStr string

	fmt.Println("Cleaning intents..")

	// Obtaining open intent struct:
	intentsB, err := stub.GetState(openIntentsStr)

	if err != nil {
		return errors.New("Error: Can't get open intents")
	}

	var intents AllIntents
	json.Unmarshal(intentsB, &intents)

	fmt.Println("Trades: " + strconv.Itoa(len(intents.OpenIntents)))

	for i := 0 ;  i < len(intents.OpenIntents) ; {
		fmt.Println(strconv.Itoa(i) + ". Searching trade: " + strconv.FormatInt(intents.OpenIntents[i].Timestamp, 10))
		fmt.Println("Options: " + strconv.Itoa(len(intents.OpenIntents[i].Willing)))

		userIDStr = strconv.Itoa(intents.OpenIntents[i].UserID)

		for x := 0 ;  x < len(intents.OpenIntents[i].Willing) ; {
			fmt.Println("Next Option: " + strconv.Itoa(i) + " : " + strconv.Itoa(x))
			_, err := findIntendedVote(stub, userIDStr, intents.OpenIntents[i].Willing[x].Choice)

			if (err != nil) {
				fmt.Println("Error. Removing option..")
				valid = true
				intents.OpenIntents[i].Willing = append(intents.OpenIntents[i].Willing[:x], intents.OpenIntents[i].Willing[x+1:]...)
				x--
			} else {
				fmt.Println("Option is valid")
			}

			x++
			fmt.Println()

			if x >= len(intents.OpenIntents[i].Willing) {
				break
			}
		}

		if len(intents.OpenIntents[i].Willing) == 0 {
			fmt.Println("Ran out of options. Removing intent..")
			valid = true
			intents.OpenIntents = append(intents.OpenIntents[:i], intents.OpenIntents[i+1:]...)
			i--
		}

		i++
		fmt.Println("i: " + strconv.Itoa(i))

		if i >= len(intents.OpenIntents) {
			break
		}
	}

	if (valid) {
		fmt.Println("Saving changes on open intents..")
		jsonB, _ := json.Marshal(intents)
		err = stub.PutState(openIntentsStr, jsonB)

		if err != nil {
			return err
			}
	} else {
		fmt.Println("All open intents are valid")
	}

	fmt.Println("Ending intent clean-up..")
		
	return nil
}