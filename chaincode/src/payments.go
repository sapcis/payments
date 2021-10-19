package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	// "github.com/hyperledger/fabric-chaincode-go/shim"
	// "github.com/hyperledger/fabric-protos-go/peer"
)

func Success(rc int32, doc string, payload []byte) peer.Response {
	return peer.Response{
		Status:  rc,
		Message: doc,
		Payload: payload,
	}
}

func Error(rc int32, doc string) peer.Response {
	logger.Errorf("Error %d = %s", rc, doc)
	return peer.Response{
		Status:  rc,
		Message: doc,
	}
}

func Validate(funcName string, args []string, desc ...interface{}) peer.Response {

	logger.Debugf("Function: %s(%s)", funcName, strings.TrimSpace(strings.Join(args, ",")))

	nrArgs := len(desc) / 3

	if len(args) != nrArgs {
		return Error(http.StatusBadRequest, "Parameter Mismatch")
	}

	for i := 0; i < nrArgs; i++ {
		switch desc[i*3] {

		case "%json":
			var jsonData map[string]interface{}
			if jsonErr := json.Unmarshal([]byte(args[i]), &jsonData); jsonErr != nil {
				return Error(http.StatusBadRequest, "JSON Payload Not Valid")
			}
			fallthrough

		case "%s":
			var minLen = desc[i*3+1].(int)
			var maxLen = desc[i*3+2].(int)
			if len(args[i]) < minLen || len(args[i]) > maxLen {
				return Error(http.StatusBadRequest, "Parameter Length Error")
			}
		}
	}

	return Success(0, "OK", nil)
}

var logger = shim.NewLogger("chaincode")

type PaymentsChaincode struct {
}

func main() {
	if err := shim.Start(new(PaymentsChaincode)); err != nil {
		fmt.Printf("Main: Error starting chaincode: %s", err)
	}
	logger.SetLevel(shim.LogDebug)
}

func (cc *PaymentsChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return Success(http.StatusNoContent, "OK", nil)
}

func (cc *PaymentsChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	function, args := stub.GetFunctionAndParameters()

	switch strings.ToLower(function) {
	case "create":
		return cc.create(stub, args)
	case "read":
		return cc.read(stub, args)
	case "update":
		return cc.update(stub, args)
	case "delete":
		return cc.delete(stub, args)
	case "list":
		return cc.list(stub, args)
	default:
		logger.Warningf("Invoke('%s') invalid!", function)
		return Error(http.StatusNotImplemented, "Invalid method! Valid methods are 'create|read|update|delete|list'!")
	}
}

func (cc *PaymentsChaincode) create(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if rc := Validate("create", args, "%s", 1, 255, "%json", 2, 4096); rc.Status > 0 {
		return rc
	}

	if msg, err := stub.GetState(args[0]); err != nil || msg != nil {
		return Error(http.StatusLocked, "Already Exists")
	}

	if err := stub.PutState(args[0], []byte(args[1])); err == nil {
		if err := stub.SetEvent("create", []byte(args[1])); err != nil {
			return shim.Error(err.Error())
		}
		return Success(http.StatusCreated, "Created", nil)
	} else {
		return Error(http.StatusInternalServerError, err.Error())
	}
}

func (cc *PaymentsChaincode) read(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if rc := Validate("read", args, "%s", 1, 255); rc.Status > 0 {
		return rc
	}

	if msg, err := stub.GetState(args[0]); err == nil && msg != nil {
		return Success(http.StatusOK, "OK", msg)
	} else {
		return Error(http.StatusNotFound, "Not Found")
	}
}

func (cc *PaymentsChaincode) update(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if rc := Validate("update", args, "%s", 1, 255, "%json", 2, 4096); rc.Status > 0 {
		return rc
	}

	if msg, err := stub.GetState(args[0]); err != nil || msg == nil {
		return Error(http.StatusNotFound, "Not Found")
	}

	if err := stub.PutState(args[0], []byte(args[1])); err == nil {
		if err := stub.SetEvent("update", []byte(args[1])); err != nil {
			return shim.Error(err.Error())
		}
		return Success(http.StatusNoContent, "Updated", nil)
	} else {
		return Error(http.StatusInternalServerError, err.Error())
	}
}

func (cc *PaymentsChaincode) delete(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if rc := Validate("delete", args, "%s", 1, 255); rc.Status > 0 {
		return rc
	}

	if msg, err := stub.GetState(args[0]); err != nil || msg == nil {
		return Error(http.StatusNotFound, "Not Found")
	}

	if err := stub.DelState(args[0]); err == nil {
		if err := stub.SetEvent("delete", []byte(args[0])); err != nil {
			return shim.Error(err.Error())
		}
		return Success(http.StatusNoContent, "Deleted", nil)
	} else {
		return Error(http.StatusInternalServerError, err.Error())
	}
}

func (cc *PaymentsChaincode) list(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 0 {
		return Error(http.StatusBadRequest, "Parameter Mismatch")
	}

	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("{ \"ids\": [")

	for resultsIterator.HasNext() {
		it, _ := resultsIterator.Next()
		buffer.WriteString("\"")
		buffer.WriteString(it.Key)
		buffer.WriteString("\",")
	}

	buffer.Truncate(buffer.Len() - 1)
	buffer.WriteString("]}")
	return Success(http.StatusOK, "OK", buffer.Bytes())
}
