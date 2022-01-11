/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

import (
	"log"
	"strings"
)

func CheckCreateGrammarTable(rules *GrammarRuleDefinitions) {
	if rules.Grammar != nil {
		return
	}
	if rules.BaseGrammar == nil {
		rules.BaseGrammar = CommonGrammarBaseDefinition
	}
	rules.Grammar = &GrammarTable{
		IsSingleQuoteTypeChar: (rules.EvaluateOptions & SINGLE_QUOTE_HAS_TYPE_CHAR) != 0,
	}
	blocks := buildStandardTokenTable(rules)
	provideRecognizers(rules, rules.Recognizers, blocks)
	provideRecognizers(rules, standardRecognizers, blocks)
	buildSpecificTokenTable(rules)
}

func provideRecognizers(rules *GrammarRuleDefinitions, recognizers []TypeRecognizer, blocks []*TokenBlock) {
	flags := rules.EvaluateOptions
	for _, v := range recognizers {
		if (v.exceptFlag & flags) == 0 {
			for _, k := range v.StartByteTypes {
				for _, b := range blocks {
					if b.Kind == k {
						if b.Recognizers == nil {
							b.Recognizers = make([]TypeRecognizer, 1)
							b.Recognizers[0] = v
						} else {
							b.Recognizers = append(b.Recognizers, v)
						}
					}
				}
			}
		}
	}
}

func buildStandardTokenTable(rules *GrammarRuleDefinitions) []*TokenBlock {
	var TokenOnlyInsideString = &TokenBlock{
		Kind:     BYTE_ONLY_INSIDE_STRING,
		IsCommon: true,
	}
	var TokenControl = &TokenBlock{
		Kind:     BYTE_CONTROL,
		IsCommon: true,
	}
	var TokenDigit = &TokenBlock{
		Kind:     BYTE_DIGIT,
		IsCommon: true,
	}
	var TokenLetter = &TokenBlock{
		Kind:     BYTE_LETTER,
		IsCommon: true,
	}
	var TokenListClosing = &TokenBlock{
		Kind:     BYTE_LIST_CLOSING,
		IsCommon: true,
	}
	var TokenListSeparator = &TokenBlock{
		Kind:     BYTE_LIST_SEPARATOR,
		IsCommon: true,
	}
	var TokenQuotation = &TokenBlock{
		Kind:     BYTE_QUOTATION,
		IsCommon: true,
	}
	var TokenSpace = &TokenBlock{
		Kind:     BYTE_SPACE,
		IsCommon: true,
	}

	tokenMap := &rules.Grammar.TokenMap
	for i := 0; i < 33; i++ {
		tokenMap[i] = TokenSpace
	}
	for i := 33; i < 48; i++ {
		tokenMap[i] = TokenOnlyInsideString
	}
	for i := 48; i < 58; i++ {
		tokenMap[i] = TokenDigit
	}
	for i := 58; i < 65; i++ {
		tokenMap[i] = TokenOnlyInsideString
	}
	for i := 65; i < 91; i++ {
		tokenMap[i] = TokenLetter
		tokenMap[i+32] = TokenLetter
	}
	for i := 91; i < 97; i++ {
		tokenMap[i] = TokenOnlyInsideString
		tokenMap[i+32] = TokenOnlyInsideString
	}
	nonLatinLettersKind := TokenLetter
	if (rules.EvaluateOptions & NON_LATIN_LETTERS_ONLY_IN_STRINGS) == 0 {
		nonLatinLettersKind = TokenOnlyInsideString
	}
	for i := 128; i < 256; i++ {
		tokenMap[i] = nonLatinLettersKind
	}
	tokenMap['_'] = TokenLetter
	tokenMap['`'] = TokenQuotation
	tokenMap['\''] = TokenQuotation
	tokenMap['"'] = TokenQuotation
	tokenMap[','] = TokenListSeparator
	tokenMap[';'] = TokenListSeparator

	tokenMap['('] = TokenListClosing
	tokenMap[')'] = TokenListClosing
	tokenMap['['] = TokenListClosing
	tokenMap[']'] = TokenListClosing
	tokenMap['{'] = TokenListClosing
	tokenMap['}'] = TokenListClosing
	tokenMap['.'] = TokenControl
	return []*TokenBlock{
		TokenOnlyInsideString,
		TokenControl,
		TokenDigit,
		TokenLetter,
		TokenListSeparator,
		TokenListClosing,
		TokenQuotation,
		TokenSpace,
	}
}

func addTokenControl(tokenMap *TokenMapType, key string, value string) {
	keyWords := strings.Split(key, " ")
	for i := 0; i < len(keyWords); i++ {
		if keyWords[i] == "" {
			keyWords = append(keyWords[:i], keyWords[i+1:]...)
			i--
		}
	}
	maxWords := len(keyWords)
	if maxWords == 0 {
		return
	}
	key = strings.Join(keyWords, " ")
	n := len(key)
	c := key[0]
	b := tokenMap[c]
	if b.IsCommon {
		kind := b.Kind
		recongnizers:=b.Recognizers
		if kind == BYTE_ONLY_INSIDE_STRING {
			kind = BYTE_OPERATOR
			recongnizers = nil
		}
		b = &TokenBlock{
			Kind:      kind,
			Controls:  make([]map[string]string, n),
			MaxWords:  maxWords,
			MaxLength: n,
			Recognizers: recongnizers,
		}
		tokenMap[c] = b
	} else {
		if b.MaxLength < n {
			b.MaxLength = n
		}
		if b.MaxWords < maxWords {
			b.MaxWords = maxWords
		}
		if len(b.Controls) < n {
			b.Controls = append(b.Controls, make([]map[string]string, n-len(b.Controls))...)
		}
	}
	if b.Controls[n-1] == nil {
		b.Controls[n-1] = make(map[string]string)
		b.Controls[n-1][key] = value
	} else {
		val, ok := b.Controls[n-1][key]
		if !ok {
			b.Controls[n-1][key] = value
		} else if val != value {
			if value == key {
				b.Controls[n-1][key] = value
				value = val
				val = b.Controls[n-1][key]
			}
			log.Printf("Warning: there are duplicates in rules: %s and %s for %s", val, value, key)
		}
	}
}

func buildSpecificTokenTable(rules *GrammarRuleDefinitions) {
	tokenMap := &rules.Grammar.TokenMap
	for k, _ := range rules.Visitors {
		if _, ok := rules.BaseGrammar.Operators[k]; !ok {
			log.Printf("Warning: Operator %s is not described in base grammar", k)
		} else {
			addTokenControl(tokenMap, k, k)
		}
	}
	for k, _ := range rules.UnaryPreVisitors {
		addTokenControl(tokenMap, k, k)
	}
	for k, _ := range rules.UnaryPostVisitors {
		addTokenControl(tokenMap, k, k)
	}
	for k, v := range rules.GrammarAliases {
		_, ok := rules.Visitors[v]
		if !ok {
			_, ok = rules.UnaryPreVisitors[v]
			if !ok {
				_, ok = rules.UnaryPostVisitors[v]
			}
		}
		if ok {
			addTokenControl(tokenMap, k, v)
		} else {
			log.Printf("Unknown alias value: %s for key %s", v, k)
		}
	}
}
