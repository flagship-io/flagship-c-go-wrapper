package main

import (
	"C"

	"fmt"
	"os"
	"strings"
	"time"

	"github.com/flagship-io/flagship-go-sdk/v2"
	"github.com/flagship-io/flagship-go-sdk/v2/pkg/bucketing"
	"github.com/flagship-io/flagship-go-sdk/v2/pkg/client"
	"github.com/flagship-io/flagship-go-sdk/v2/pkg/logging"
	"github.com/sirupsen/logrus"
)
import (
	"sort"

	"github.com/flagship-io/flagship-go-sdk/v2/pkg/model"
)

var fsClient *client.Client

type FakeTrackingAPIClient struct{}

func (*FakeTrackingAPIClient) SendHit(visitorID string, anonymousID *string, hit model.HitInterface) error {
	return nil
}
func (*FakeTrackingAPIClient) ActivateCampaign(request model.ActivationHit) error { return nil }
func (*FakeTrackingAPIClient) SendEvent(request model.Event) error {
	return nil
}

func main() {
	initFlagship(C.CString("blvo2kijq6pg023l8edg"), C.CString("wwURPfcEB01YVMfTYWfCtaezCkXVLeFZ61FJmXtI"), 60, C.CString("info"), 0)

	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			getAllFlags(C.CString("visitor_id"), C.CString(""))
		}
	}
}

//export initFlagship
func initFlagship(environmentID *C.char, apiKey *C.char, polling C.int, logLevel *C.char, trackingEnabled C.int) {
	level, err := logrus.ParseLevel(C.GoString(logLevel))
	if err != nil {
		fmt.Printf("err: %s\n", err)
		os.Exit(1)
	}

	logging.SetLevel(level)

	// set bucketing options with custom polling interval
	options := []client.OptionBuilder{
		client.WithBucketing(bucketing.PollingInterval(time.Duration(polling) * time.Second)),
	}

	// if tracking is disabled, set custom "fake" trackging api client
	if trackingEnabled == 0 {
		options = append(options, client.WithTrackingAPIClient(&FakeTrackingAPIClient{}))
	}

	fsClient, err = flagship.Start(C.GoString(environmentID), C.GoString(apiKey), options...)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		os.Exit(1)
	}
}

func createVisitor(visitorID *C.char, contextString *C.char) *client.Visitor {
	fsVisitor, err := fsClient.NewVisitor(C.GoString(visitorID), extractContext(contextString))
	if err != nil {
		fmt.Printf("err: %s\n", err)
		os.Exit(1)
	}
	return fsVisitor
}

func extractContext(contextString *C.char) map[string]interface{} {
	context := map[string]interface{}{}
	contextInfos := strings.Split(C.GoString(contextString), ";")
	for _, cKV := range contextInfos {
		cKVInfos := strings.Split(cKV, ":")
		if len(cKVInfos) == 2 {
			context[cKVInfos[0]] = cKVInfos[1]
		}
	}
	return context
}

//export getFlagBool
func getFlagBool(visitorID *C.char, contextString *C.char, key *C.char, defaultValue C.int, activate C.int) C.int {
	fsVisitor := createVisitor(visitorID, contextString)
	err := fsVisitor.SynchronizeModifications()
	if err != nil {
		fmt.Printf("err when synchronizing visitor: %s\n", err)
		return defaultValue
	}

	flag, err := fsVisitor.GetModificationBool(C.GoString(key), defaultValue == 1, activate == 1)
	if err != nil {
		return defaultValue
	}

	var ret C.int = 0
	if flag {
		ret = 1
	}
	return ret
}

//export getFlagNumber
func getFlagNumber(visitorID *C.char, contextString *C.char, key *C.char, defaultValue C.double, activate C.int) C.double {
	fsVisitor := createVisitor(visitorID, contextString)
	err := fsVisitor.SynchronizeModifications()
	if err != nil {
		fmt.Printf("err when synchronizing visitor: %s\n", err)
		return defaultValue
	}

	flag, err := fsVisitor.GetModificationNumber(C.GoString(key), float64(defaultValue), activate == 1)
	if err != nil {
		return defaultValue
	}
	return C.double(flag)
}

//export getFlagString
func getFlagString(visitorID *C.char, contextString *C.char, key *C.char, defaultValue *C.char, activate C.int) *C.char {
	fsVisitor := createVisitor(visitorID, contextString)
	err := fsVisitor.SynchronizeModifications()
	if err != nil {
		fmt.Printf("err when synchronizing visitor: %s\n", err)
		return defaultValue
	}

	flag, err := fsVisitor.GetModificationString(C.GoString(key), C.GoString(defaultValue), activate == 1)
	if err != nil {
		return defaultValue
	}
	return C.CString(flag)
}

//export getAllFlags
func getAllFlags(visitorID *C.char, contextString *C.char) *C.char {
	fsVisitor := createVisitor(visitorID, contextString)

	err := fsVisitor.SynchronizeModifications()
	if err != nil {
		fmt.Printf("err: %s\n", err)
		os.Exit(1)
	}

	flags := fsVisitor.GetAllModifications()

	flagsString := ""
	keys := []string{}
	for k := range flags {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		flagsString += fmt.Sprintf("%s:%v;", k, flags[k].Value)
	}

	return C.CString(flagsString)
}
