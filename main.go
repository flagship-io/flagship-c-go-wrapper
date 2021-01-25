package main

/*
#include <stdlib.h>
typedef struct FSValue {
    char *name;   // name for the key context
    int varType; // 1 = string, 2 = bool, 3 = float, 4 = int
    void *data;   // the value
} FSValue;
// UserContext
typedef struct FSContext{
    int numberOfAttribute;
    struct FSValue *userContextList;  // list of context
} FSContext;
// Modification
typedef struct FSModifications{
    int numberOfFlags;
    struct FSValue *flagList;
} FSModifications;
*/
import "C"

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
	"unsafe"

	"github.com/abtasty/flagship-go-sdk/v2"
	"github.com/abtasty/flagship-go-sdk/v2/pkg/bucketing"
	"github.com/abtasty/flagship-go-sdk/v2/pkg/client"
	"github.com/abtasty/flagship-go-sdk/v2/pkg/logging"
	"github.com/sirupsen/logrus"
)

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

	log.Println("------------ initFlagship from GO wrapper library ------------ ")

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

/*  Added */

//export getVisitorContext
func getVisitorContext() *C.char {
	log.Println("------------ getVisitorContext ------------ ")

	// typedef struct userContext {
	//     char *name;   // name for the key context
	//     int var_type; // 1 = string, 2 = bool, 3 = float, 4 = int
	//     void *data;   // the value
	// } userContext;
	//

	//attribs := C.struct_userContexts{C.CString("the current context from GO package "),1,nil}

	return C.CString("the current context from GO package ")
}

//export initCFlagship
func initCFlagship(environmentID *C.char, apiKey *C.char, polling C.int, logLevel *C.char, context *C.struct_FSContext) {
	var err error

	log.Println("------------ initCFlagship from GO wrapper library ------------ ")

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

//export getAllCFlags
func getAllCFlags(visitorID *C.char, contextString *C.char) C.struct_FSModifications {

	log.Println("------------ getAllCFlags ------------ ")

	/// For the moment the struct_FSContext is not connected

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

	/// Work in progress

	//// list
	n := C.ulong(len(keys))
	list := (*C.FSValue)(C.malloc(C.sizeof_FSValue * n))

	var arrayValueStruct (*C.FSValue) = list
	lenght := C.ulong(len(keys))

	slicePointer := (*[1 << 30](C.FSValue))(unsafe.Pointer(arrayValueStruct))[:lenght:lenght]

	for i, k := range keys {

		typeData := 0
		switch v := flags[k].Value.(type) {
		case int:
			//	fmt.Println("&&&&&&&&&&&&&& int: &&&&&&&&&&&&&&&&&&&&", v)
			typeData = 4
		case float64:
			//fmt.Println("&&&&&&&&&&&&&&&& float64: &&&&&&&&&&&&&&", v)
			typeData = 4
		case string:
			//fmt.Println("&&&&&&&&&&& string &&&&&&&&&&&&&&:\n", v)
			typeData = 1
		case bool:
			//fmt.Println("&&&&&&&&&&& Boolean &&&&&&&&&&&&&&:\n", v)
			typeData = 3
		default:
			fmt.Println("unknown value for the key: ", v)
		}

		slicePointer[i].varType = C.int(typeData)

		slicePointer[i].name = C.CString(k)
		valueString := fmt.Sprintf("%v", flags[k].Value)
		slicePointer[i].data = unsafe.Pointer(C.CString(valueString))
	}

	return C.struct_FSModifications{C.int(len(keys)), &slicePointer[0]}
}
