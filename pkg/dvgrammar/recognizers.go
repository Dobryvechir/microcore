/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

import (
	"errors"
)

func RecognizerByteLetter(grammar *GrammarTable, data []byte, pos int, n int) (t *Token, nextPos int, err error) {
	i := pos
	for ; i < n; i++ {
		c := data[i]
		b := grammar.TokenMap[c].Kind
		if b != BYTE_LETTER && b != BYTE_DIGIT {
			break
		}
	}
	t = &Token{DataType: TYPE_DATA, Value: string(data[pos:i])}
	nextPos = i
	return
}

func RecognizerByteDigit(grammar *GrammarTable, data []byte, pos int, n int) (t *Token, nextPos int, err error) {
	nextPos = readNumber(data, pos, n) + 1
	t = &Token{DataType: TYPE_NUMBER, Value: string(data[pos:nextPos])}
	return
}

func RecognizerByteListControl(grammar *GrammarTable, data []byte, pos int, n int) (t *Token, nextPos int, err error) {
	t = &Token{DataType: TYPE_CONTROL, Value: string(data[pos : pos+1])}
	nextPos = pos + 1
	return
}

func RecognizerByteQuotation(grammar *GrammarTable, data []byte, pos int, n int) (t *Token, nextPos int, err error) {
	c := data[pos]
	i := pos
	switch c {
	case '`':
		g := 0
		for i < n && data[i] == '`' && g < 3 {
			g++
			i++
		}
		for i < n {
			ge := data[i]
			if ge == '`' {
				gc := 0
				for i < n && data[i] == '`' && gc < g {
					gc++
					i++
				}
				if gc == g {
					i--
					break
				}
			}
		}
		if i == n {
			return nil, 0, errors.New("Unclosed string with `")
		}
		i++
		t = &Token{DataType: TYPE_STRING, Value: string(data[pos+g : i-g])}
	case '"', '\'':
		i++
		for i < n && data[i] != '"' {
			ge := data[i]
			if ge == 10 || ge == 13 {
				ending := string([]byte{c})
				return nil, 0, errors.New("Unclosed string with " + ending)
			}
			if ge == '\\' {
				i++
			}
			i++
		}
		if i < n {
			dataType := TYPE_STRING
			if c == '\'' && grammar.IsSingleQuoteTypeChar {
				dataType = TYPE_CHAR
			}
			t = &Token{DataType: dataType, Value: GetEscapedString(data[pos+1 : i])}
			i++
		} else {
			ending := string([]byte{c})
			return nil, 0, errors.New("Unclosed string with " + ending)
		}
	}
	nextPos = i
	return
}
func RecognizerByteControl(grammar *GrammarTable, data []byte, pos int, n int) (t *Token, nextPos int, err error) {
	if data[pos] == '.' && pos+1 < n && data[pos+1] <= '9' && data[pos+1] >= '0' {
		nextPos = readNumber(data, pos, n) + 1
		t = &Token{DataType: TYPE_NUMBER, Value: string(data[pos:nextPos])}
	} else {
		t = &Token{DataType: TYPE_CONTROL, Value: string(data[pos : pos+1])}
		nextPos = pos + 1
	}
	return
}

var standardRecognizers = []TypeRecognizer{
	{
		StartByteTypes: []int{BYTE_LETTER},
		Processor:      RecognizerByteLetter,
		exceptFlag:     EXCLUDE_RECOGNIZER_BYTE_LETTER,
	},
	{
		StartByteTypes: []int{BYTE_DIGIT},
		Processor:      RecognizerByteDigit,
		exceptFlag:     EXCLUDE_RECOGNIZER_BYTE_DIGIT,
	},
	{
		StartByteTypes: []int{BYTE_QUOTATION},
		Processor:      RecognizerByteQuotation,
		exceptFlag:     EXCLUDE_RECOGNIZER_BYTE_QUOTATION,
	},
	{
		StartByteTypes: []int{BYTE_LIST_CLOSING, BYTE_LIST_SEPARATOR},
		Processor:      RecognizerByteListControl,
		exceptFlag:     EXCLUDE_RECOGNIZER_BYTE_LIST_CONTROL,
	},
	{
		StartByteTypes: []int{BYTE_CONTROL},
		Processor:      RecognizerByteControl,
		exceptFlag:     EXCLUDE_RECOGNIZER_BYTE_CONTROL,
	},
}
