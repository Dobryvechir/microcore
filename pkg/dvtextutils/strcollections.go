/***********************************************************************
MicroCore
Copyright 2020 -2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvtextutils

func IsStringContainedInArray(s string, collection []string) bool {
	n:=len(collection)
	for i:=0;i<n;i++ {
		if collection[i]==s {
			return true
		}
	}
	return false
}

func IsStringContainedInStringToStringMap(s string,collection map[string]string) bool {
	_,res:=collection[s]
	return res
}

func IsStringContainedInStringToAnyMap(s string,collection map[string]interface{}) bool {
	_,res:=collection[s]
	return res
}

