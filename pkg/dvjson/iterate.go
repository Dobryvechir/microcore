/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

type IterateProcessor func (string,interface{},int,interface{}) (interface{},bool)

func IterateOnAnyType(val interface{},processor IterateProcessor,initial interface{}) interface{} {
	res:=initial
	index:=0
	toBreak:=false
	switch val.(type) {
	case map[string]string:
		for k,v:=range val.(map[string]string) {
			res,toBreak = processor(k, v, index, res)
			if toBreak {
				break
			}
			index++
		}
	case map[string]interface{}:
		for k,v:=range val.(map[string]interface{}) {
			res,toBreak = processor(k, v, index, res)
			if toBreak {
				break
			}
			index++
		}
	case *DvFieldInfo:
		fieldInfo:=val.(*DvFieldInfo)
		fields:=fieldInfo.Fields
		n:=len(fields)
		for index=0;index<n;index++ {
			f:=fields[index]
			res,toBreak = processor(string(f.Name), f, index, res)
			if toBreak {
				break
			}
		}
	}
	return res
}
