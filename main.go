package main

import (
	"C"

	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/abtasty/flagship-go-sdk/v2"
	"github.com/abtasty/flagship-go-sdk/v2/pkg/bucketing"
	"github.com/abtasty/flagship-go-sdk/v2/pkg/client"
	"github.com/abtasty/flagship-go-sdk/v2/pkg/logging"
	"github.com/sirupsen/logrus"
)
import "sort"

var fsClient *client.Client

func main() {
	initFlagship(C.CString("blvo2kijq6pg023l8edg"), C.CString("wwURPfcEB01YVMfTYWfCtaezCkXVLeFZ61FJmXtI"), 60, C.CString("info"))

	ticker := time.Ticker{}

	for {
		select {
		case <-ticker.C:
			log.Println("interval")
		}
	}

}

//export initFlagship
func initFlagship(environmentID *C.char, apiKey *C.char, polling C.int, logLevel *C.char) {
	var err error

	switch C.GoString(logLevel) {
	case "debug":
		logging.SetLevel(logrus.DebugLevel)
	case "info":
		logging.SetLevel(logrus.InfoLevel)
	case "warn":
		logging.SetLevel(logrus.WarnLevel)
	case "error":
		logging.SetLevel(logrus.ErrorLevel)
	default:
		logging.SetLevel(logrus.WarnLevel)
	}
	fsClient, err = flagship.Start(C.GoString(environmentID), C.GoString(apiKey), client.WithBucketing(bucketing.PollingInterval(time.Duration(polling)*time.Second)))
	if err != nil {
		fmt.Printf("err: %s\n", err)
		os.Exit(1)
	}
}

//export getAllFlags
func getAllFlags(visitorID *C.char, contextString *C.char) *C.char {
	context := map[string]interface{}{}
	contextInfos := strings.Split(C.GoString(contextString), ";")
	for _, cKV := range contextInfos {
		cKVInfos := strings.Split(cKV, ":")
		if len(cKVInfos) == 2 {
			context[cKVInfos[0]] = cKVInfos[1]
		}
	}
	fsVisitor, err := fsClient.NewVisitor(C.GoString(visitorID), context)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		os.Exit(1)
	}

	err = fsVisitor.SynchronizeModifications()
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
