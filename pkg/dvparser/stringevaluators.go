package dvparser

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"log"
)

func ByteArrayEvaluatorAsInterface(data []byte, extraParams *dvevaluation.DvObject, source string) (interface{}, error) {
	res := dvevaluation.ParseForDvObject(data, extraParams, 1, 1, source)
	return res.FinalResult, res.Err
}

func StringEvaluatorAsString(data string, extraParams *dvevaluation.DvObject) (string, error) {
	res, err := ByteArrayEvaluatorAsInterface([]byte(data), extraParams, data)
	if err != nil {
		return "", err
	}
	return dvevaluation.AnyToString(res), nil

}

func StringEvaluatorAsStringWithErrorLoggedAsWarning(data string, extraParams *dvevaluation.DvObject, defaultIfError string) string {
	res, err := StringEvaluatorAsString(data, extraParams)
	if err != nil {
		log.Printf("Error in %s: %s", data, err.Error())
		res = defaultIfError
	}
	return res
}

func StringEvaluatorAsBoolean(data string, extraParams *dvevaluation.DvObject) (bool, error) {
	res, err := ByteArrayEvaluatorAsInterface([]byte(data), extraParams, data)
	if err != nil {
		return false, err
	}
	return dvevaluation.AnyToBoolean(res), nil
}

func StringEvaluatorAsBooleanWithErrorLoggedAsWarning(data string, extraParams *dvevaluation.DvObject, defaultIfError bool) bool {
	res, err := StringEvaluatorAsBoolean(data, extraParams)
	if err != nil {
		log.Printf("Error in %s: %s", data, err.Error())
		res = defaultIfError
	}
	return res
}

func CleanLogByteArray(data []byte) string {
	n := len(data)
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		c := data[i]
		if c < ' ' {
			res[i] = ' '
		} else {
			res[i] = c
		}
	}
	return string(res)
}

func CleanLogString(data string) string {
	return CleanLogByteArray([]byte(data))
}
