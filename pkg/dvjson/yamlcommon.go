/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import "github.com/Dobryvechir/microcore/pkg/dvevaluation"

func PrintToJsonAtLevel(dvEntry *dvevaluation.DvVariable, res []byte, level int, indent int, noIndentAtFirst bool) []byte {
	if dvEntry == nil {
		return res
	}
	n := indent * level
	nextN := n + indent
	nextLevel := level + 1
	var indentBuf, nextIndentBuf []byte
	if indent > 0 {
		indentBuf = make([]byte, n)
		for i := 0; i < n; i++ {
			indentBuf[i] = ' '
		}

		if !noIndentAtFirst {
			res = append(res, indentBuf...)
		}
	}
	switch dvEntry.Kind {
	case dvevaluation.FIELD_STRING:
		res = appendJsonEscapedString(res, dvEntry.Value)
	case dvevaluation.FIELD_ARRAY:
		arrayAmount := len(dvEntry.Fields)
		res = append(res, '[')
		if arrayAmount > 0 {
			if indent > 0 {
				res = append(res, byte(10))
				nextIndentBuf = make([]byte, nextN)
				for i := 0; i < nextN; i++ {
					nextIndentBuf[i] = ' '
				}
			}

			for i := 0; i < arrayAmount; i++ {
				if i > 0 {
					res = append(res, ',', byte(10))
				}
				if indent > 0 {
					res = append(res, nextIndentBuf...)
				}
				res = PrintToJsonAtLevel(dvEntry.Fields[i], res, nextLevel, indent, true)
			}
			if indent > 0 {
				res = append(res, byte(10))
				res = append(res, indentBuf...)
			}
		}
		res = append(res, ']')

	case dvevaluation.FIELD_OBJECT:
		mapAmount := len(dvEntry.Fields)
		res = append(res, '{')
		if mapAmount > 0 {
			if indent > 0 {
				res = append(res, byte(10))
				nextIndentBuf = make([]byte, nextN)
				for i := 0; i < nextN; i++ {
					nextIndentBuf[i] = ' '
				}
			}
			isNext := false
			entryList := dvEntry.Fields
			for _, v := range entryList {
				if isNext {
					res = append(res, ',', byte(10))
				} else {
					isNext = true
				}
				if indent > 0 {
					res = append(res, nextIndentBuf...)
				}
				res = appendJsonEscapedString(res, v.Name)
				res = append(res, ':', ' ')
				res = PrintToJsonAtLevel(v, res, nextLevel, indent, true)
			}
			if indent > 0 {
				res = append(res, byte(10))
				res = append(res, indentBuf...)
			}
		}
		res = append(res, '}')

	default:
		res = append(res, dvEntry.Value...)
	}
	return res
}

func appendJsonEscapedString(res []byte, add []byte) []byte {
	n := len(add)
	res = append(res, '"')
	for i := 0; i < n; i++ {
		b := add[i]
		if b == '"' || b == '\\' {
			res = append(res, '\\')
		}
		res = append(res, b)
	}
	res = append(res, '"')
	return res
}

func isReservedWord(data []byte) bool {
	word := string(data)
	return word == "null" || word == "false" || word == "true"
}

func appendYamlEscapedString(res []byte, add []byte) []byte {
	n := len(add)
	isNumber := true
	isSimple := true
	for i := 0; i < n; i++ {
		b := add[i]
		if !(b >= '0' && b <= '9') {
			if b != '.' && b != '+' && b != 'e' && b != '-' {
				isNumber = false
			}
			if !(b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b > 127 || (i > 0 && i < n-1 && (b == '.' || b == '/' || b == '-'))) {
				isSimple = false
				if !isNumber {
					break
				}
			}
		}
	}
	if isNumber {
		res = append(res, '\'')
		res = append(res, add...)
		res = append(res, '\'')
		return res
	}
	if isSimple && isReservedWord(add) {
		isSimple = false
	}
	if isSimple {
		res = append(res, add...)
		return res
	}
	res = append(res, '"')
	for i := 0; i < n; i++ {
		b := add[i]
		if b == '"' || b == '\\' {
			res = append(res, '\\')
		}
		res = append(res, b)
	}
	res = append(res, '"')
	return res
}

func PrintToYamlAtLevel(dvEntry *dvevaluation.DvVariable, res []byte, level int, indent int, noIndentAtFirst bool) []byte {
	if dvEntry == nil {
		return res
	}
	n := indent * level
	nextLevel := level + 1
	var indentBuf []byte
	indentBuf = make([]byte, n+1)
	indentBuf[0] = byte(10)
	for i := 1; i <= n; i++ {
		indentBuf[i] = ' '
	}
	switch dvEntry.Kind {
	case dvevaluation.FIELD_STRING:
		if !noIndentAtFirst {
			res = append(res, indentBuf...)
		}
		res = appendYamlEscapedString(res, dvEntry.Value)
	case dvevaluation.FIELD_ARRAY:
		arrayAmount := len(dvEntry.Fields)
		if arrayAmount > 0 {
			for i := 0; i < arrayAmount; i++ {
				res = append(res, indentBuf...)
				res = append(res, '-', ' ')
				res = PrintToYamlAtLevel(dvEntry.Fields[i], res, nextLevel, indent, true)
			}
		} else {
			res = append(res, '[', ']')
		}
	case dvevaluation.FIELD_OBJECT:
		mapAmount := len(dvEntry.Fields)
		if mapAmount > 0 {
			entryList := dvEntry.Fields
			for _, v := range entryList {
				res = append(res, indentBuf...)
				res = appendYamlEscapedString(res, v.Name)
				res = append(res, ':', ' ')
				res = PrintToYamlAtLevel(v, res, nextLevel, indent, true)
			}
		} else {
			res = append(res, '{', '}')
		}
	default:
		if !noIndentAtFirst {
			res = append(res, indentBuf...)
		}
		res = append(res, dvEntry.Value...)
	}
	return res
}

func PrintToJson(dvEntry *dvevaluation.DvVariable, indent int) []byte {
	res := make([]byte, 0, 64000)
	return PrintToJsonAtLevel(dvEntry, res, 0, indent, true)
}

func PrintToYaml(dvEntry *dvevaluation.DvVariable, indent int) []byte {
	if indent < 1 {
		indent = 2
	}
	res := make([]byte, 0, 64000)
	res = PrintToYamlAtLevel(dvEntry, res, 0, indent, true)
	n := len(res)
	if n > 0 && res[0] == 10 {
		res = res[1:]
		n--
	}
	if n > 0 && res[n-1] != 10 {
		res = append(res, byte(10))
	}
	return res
}

func IsCurrentFormatJson(data []byte) bool {
	first := getFirstNonSpaceByte(data)
	last := getLastNonSpaceByte(data)
	return first == '[' && last == ']' || first == '{' && last == '}'
}

func getFirstNonSpaceByte(data []byte) byte {
	l := len(data)
	for i := 0; i < l; i++ {
		if data[i] > 32 {
			return data[i]
		}
	}
	return 0
}

func getLastNonSpaceByte(data []byte) byte {
	l := len(data)
	for i := l - 1; i >= 0; i-- {
		if data[i] > 32 {
			return data[i]
		}
	}
	return 0
}
