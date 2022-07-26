/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

import (
	"errors"
	"strconv"
)

func Tokenize(src *SourceReference, data []byte, grammar *GrammarTable) ([]Token, error) {
	amount := len(data)
	capacity := amount >> 4
	if capacity < 4 {
		capacity = 4
	}
	tokens := make([]Token, 0, capacity)
	for i := 0; i < amount; {
		c := data[i]
		b := grammar.TokenMap[c]
		kind := b.Kind
		if kind == BYTE_SPACE {
			if c == 10 || c == 13 {
				if c == 13 || i == 0 || data[i-1] != 13 {
					src.Row++
					src.Column = 1
				}
			} else {
				src.Column++
			}
			i++
			continue
		}
		pos := i
		var t *Token
		oper, nextPos := b.findOperator(grammar, data, pos)
		if nextPos > pos {
			t = &Token{DataType: TYPE_OPERATOR, Value: oper}
			i = nextPos
		} else {
			recognizerAmount := len(b.Recognizers)
			var err error
			for j := 0; j < recognizerAmount; j++ {
				t, nextPos, err = b.Recognizers[j].Processor(grammar, data, pos, amount)
				if err != nil {
					return nil, err
				}
				if t != nil {
					i = nextPos
					break
				}
			}
		}
		if t != nil {
			t.Row = src.Row
			t.Column = src.Column
			t.Place = src.Place
			tokens = append(tokens, *t)
		} else {
			return nil, errors.New("Unknown character " + string(c) + "(" + strconv.Itoa(int(c)) + ")" + src.Place + " (" + strconv.Itoa(src.Row) + ":" + strconv.Itoa(src.Column) + ")")
		}
		for ; pos < i; pos++ {
			c := data[pos]
			if c == 13 || c == 10 {
				if c == 13 || pos == 0 || data[pos-1] != 13 {
					src.Row++
					src.Column = 1
				}
			} else {
				src.Column++
			}
		}
	}
	return tokens, nil
}

