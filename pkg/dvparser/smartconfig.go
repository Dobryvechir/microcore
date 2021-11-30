/***********************************************************************
MicroCore
Copyright 2020 -2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvparser

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
)

const (
	CONFIG_RESULT_NOT_CONTAINING_PREMAP         = 1 << iota
	CONFIG_REPLACEMENT_NOT_MANDATORY            = 1 << iota
	CONFIG_IS_NOT_ESCAPED                       = 1 << iota
	CONFIG_IS_NOT_VARIABLES                     = 1 << iota
	CONFIG_DISABLE_INITIAL_MAP_FROM_ENVIRONMENT = 1 << iota
	CONFIG_PRESERVE_SPACE                       = 1 << iota
)

var DvParserLog bool = false
var NumberOfBracketsInConfigParsing = 3

type internalParseInfo struct {
	position int
	row      int
	column   int
}

type ConfigInfo struct {
	NumberOfBrackets int
	ParamMap		*dvevaluation.DvObject
	InputMap         map[string]string
	Options          int
	OutputMap        map[string]string
	OutputLines      [][]byte
	Err              error
	FilePaths        []string
	insider          string
}

const (
	IFELSE_ELSE_NOTRUN = 0
	IFELSE_ELSE_RUN    = 1
	IFELSE_ELSE_NOMORE = 2
	IFELSE_ELSE_ERROR  = 3
	IFELSE_LOW_MASK    = 7
	IFELSE_CONSUMED    = 16
	IFELSE_NONEWIF     = 8
)

/*************************************************************************************
fileName is the relative or full name, if it starts with X:... or / it is considered the absolute
  if it starts with other letters, it is found with regard to all FilePaths if any. If they are not given,
  it is looked for in the current folder.
level = -1 if the file not existing problem is ignored, 0, 1, ... to specify from what FilePaths to look for files
FilePaths should not contain final slashes
***************************************************************************************/

func linearSmartConfigParse_readFile(fileName string, configInfo *ConfigInfo, level int) ([]byte, string) {
	startLevel := level
	if startLevel < 0 {
		startLevel = 0
	}
	endLevel := len(configInfo.FilePaths)
	fileName = strings.TrimSpace(fileName)
	if fileName != "" {
		if !(len(fileName) > 1 && fileName[1] == ':' || fileName[0] == '/' || fileName[0] == '\\') {
			if fileName[0] == '.' && len(fileName) > 2 && (fileName[1] == '/' || fileName[1] == '\\' || fileName[1] == '.' && (fileName[2] == '/' || fileName[2] == '\\')) && configInfo.insider != "" {
				f := filepath.Dir(configInfo.insider) + "/" + fileName
				if _, err := os.Stat(f); err == nil {
					data, err1 := ioutil.ReadFile(f)
					if err1 == nil {
						fileName, err1 = filepath.Abs(f)
					}
					configInfo.Err = err1
					return data, fileName
				}
			}
			for startLevel < endLevel {
				f := configInfo.FilePaths[startLevel] + "/" + fileName
				if _, err := os.Stat(f); err == nil {
					data, err1 := ioutil.ReadFile(f)
					if err1 == nil {
						fileName, err1 = filepath.Abs(f)
					}
					configInfo.Err = err1
					return data, fileName
				}
				startLevel++
			}

		}
		if _, err := os.Stat(fileName); err == nil {
			data, err1 := ioutil.ReadFile(fileName)
			if err1 == nil {
				fileName, err1 = filepath.Abs(fileName)
			}
			configInfo.Err = err1
			return data, fileName
		}
	}
	if level >= 0 {
		configInfo.Err = errors.New("File " + fileName + " not found in any lookupfolders")
	}
	return nil, fileName
}

func getWhereInfo(pos int, place string, internalParsingInfo []internalParseInfo) string {
	l := len(internalParsingInfo)
	i := 0
	for i+1 < l && internalParsingInfo[i+1].position <= pos {
		i++
	}
	row := internalParsingInfo[i].row
	column := internalParsingInfo[i].column + pos - internalParsingInfo[i].position
	return " (" + place + " [" + strconv.Itoa(row) + ":" + strconv.Itoa(column) + "])"
}

/********************************************************************************
This function parses the following preconditions:
#if, #else, #elif #endif, #ifdef #ifndef      conditional operators
#define #undef                                define a new variable or undefine
#error                                        generate error
#include "filename" or #include <filename>    includes other files
{{{expression}}} the expression is immediately inserted in the text upon loading, based on the
current availabe variables.  The number of braces for this purpose is specified in configInfo.NumberOfBrackets
Comments can be started with # at the beginning of line or with / *** ... *** / to comment a whole block
*********************************************************************************/
func linearSmartConfigParse_internal(data []byte, configInfo *ConfigInfo, place string) {
	l := len(data)
	row := 1
	col := 1
	replacementMandatory := (configInfo.Options & CONFIG_REPLACEMENT_NOT_MANDATORY) == 0
	configIsEscaped := (configInfo.Options & CONFIG_IS_NOT_ESCAPED) == 0
	replaceNumber := configInfo.NumberOfBrackets
	isDictionary := (configInfo.Options & CONFIG_IS_NOT_VARIABLES) == 0
	preserveSpace := (configInfo.Options&CONFIG_PRESERVE_SPACE) != 0 && !isDictionary
	internalParsingInfo := make([]internalParseInfo, 1, 10)
	ifElseInfo := make([]int, 0, 10)
	currentLinePos := 0
	for i := 0; i < l; {
		currentPos := i
		for ; i < l && data[i] <= 32; i++ {
			if data[i] == 13 || data[i] == 10 {
				if data[i] == 13 || i == 0 || data[i-1] != 13 {
					row++
					col = 1
					if preserveSpace && currentPos < i {
						configInfo.OutputLines = append(configInfo.OutputLines, data[currentPos:i])
					}
				}
				currentPos = i + 1
			} else {
				col++
			}
		}
		if i == l {
			break
		}
		pos := i
		sequence := 0
		internalParsingInfo = internalParsingInfo[:1]
		internalParsingInfo[0].position = 0
		internalParsingInfo[0].row = row
		internalParsingInfo[0].column = col
		currentLinePos = 0
		var currentLine []byte = nil
		currentLineEnd := 0
		if col > 1 && preserveSpace {
			pos -= col - 1
		}
		colBasis := i
		for ; i < l && data[i] != 13 && data[i] != 10; i++ {
			if currentLinePos < currentLineEnd {
				currentLine[currentLinePos] = data[i]
				currentLinePos++
			}
			if data[i] == '{' {
				sequence = 1
				for i++; i < l && data[i] == '{'; i++ {
					sequence++
					if currentLinePos < currentLineEnd {
						currentLine[currentLinePos] = data[i]
						currentLinePos++
					}
				}
				if sequence >= replaceNumber {
					k := i - sequence
					j := k - pos
					if currentLinePos < j {
						currentLine = append(currentLine, data[k-j+currentLinePos:k]...)
						currentLineEnd = len(currentLine)
					}
					currentLinePos = j
					k = i
					j = sequence
					sequence = 0
					for ; i < l; i++ {
						if data[i] == '}' {
							sequence = 0
							for sequence = 0; i < l && data[i] == '}'; i++ {
								sequence++
							}
							if sequence == j {
								break
							}
						}
					}
					if j != sequence {
						configInfo.Err = errors.New("Unclosed expression at (" + strconv.Itoa(row) + "," + strconv.Itoa(col+currentLinePos) + ")")
						return
					}
					expr := dvevaluation.Parse(data[k:i-sequence], configInfo.InputMap, configInfo.OutputMap, row, col+colBasis-k, place)
					if expr.Err != nil {
						if replacementMandatory {
							configInfo.Err = expr.Err
							return
						} else {
							if DvParserLog {
								log.Print(expr.Err.Error())
							}
							expr.FinalResult = ""
						}
					}
					finalResult := []byte(dvevaluation.AnyToString(expr.FinalResult))
					newBlockLen := len(finalResult)
					k = currentLinePos + newBlockLen
					if k > currentLineEnd {
						currentLine = append(currentLine[:currentLinePos], finalResult...)
						currentLinePos = k
						currentLineEnd = len(currentLine)
					} else {
						for j = 0; currentLinePos < k; currentLinePos++ {
							currentLine[currentLinePos] = finalResult[j]
							j++
						}
					}
					pos = i - currentLinePos
				}
				i--
			}
		}
		if currentLine == nil {
			currentLine = data[pos:i]
		} else {
			if i-pos > currentLineEnd {
				currentLine = append(currentLine, data[pos+currentLineEnd:i]...)
			} else {
				currentLine = currentLine[:currentLinePos]
			}
		}
		subBlock := currentLine
		subBlockLen := len(subBlock)
		for subBlockLen > 0 && subBlock[subBlockLen-1] <= 32 {
			subBlockLen--
		}
		subStart := 0
		if !preserveSpace {
			for subStart < subBlockLen && subBlock[subStart] <= 32 {
				subStart++
			}
		}
		if subStart == subBlockLen {
			if preserveSpace {
				configInfo.OutputLines = append(configInfo.OutputLines, currentLine)
			}
			continue
		}
		if subBlock[subStart] == '#' {
			isComment := true
			subStart++
			b1 := byte(' ')
			if subStart < subBlockLen {
				b1 = subBlock[subStart]
				subStart++
			}
			b2 := byte(' ')
			if subStart < subBlockLen {
				b2 = subBlock[subStart]
				subStart++
			}
			b3 := byte(' ')
			if subStart < subBlockLen {
				b3 = subBlock[subStart]
				subStart++
			}
			bIfCondition := -1
			if b1 == 'i' {
				if b2 == 'f' {
					if b3 <= 32 { //if
						res := dvevaluation.EvalAsBoolean(subBlock[subStart:subBlockLen], configInfo.InputMap, configInfo.OutputMap, internalParsingInfo[0].row, internalParsingInfo[0].column+subStart, place)
						if res.Err != nil {
							configInfo.Err = res.Err
							return
						}
						if res.FinalResult {
							bIfCondition = 1
						} else {
							bIfCondition = 0
						}
					} else if b3 == 'd' { //ifdef
						if subStart+1 < subBlockLen && subBlock[subStart] == 'e' && subBlock[subStart+1] == 'f' && (subStart+2 >= subBlockLen || subBlock[subStart+2] <= 32) {
							subStart += 2
							var err error
							scope := dvevaluation.NewDvObjectFrom2Maps(configInfo.OutputMap, configInfo.InputMap)
							bIfCondition, err = dvevaluation.IsDefined(subBlock[subStart:subBlockLen], scope, internalParsingInfo[0].row, internalParsingInfo[0].column+subStart, place, 0)
							if err != nil {
								configInfo.Err = err
								return
							}
						}
					} else if b3 == 'n' { //ifndef
						if subStart+2 < subBlockLen && subBlock[subStart] == 'd' && subBlock[subStart+1] == 'e' && subBlock[subStart+2] == 'f' && (subStart+3 >= subBlockLen || subBlock[subStart+3] <= 32) {
							subStart += 3
							scope := dvevaluation.NewDvObjectFrom2Maps(configInfo.OutputMap, configInfo.InputMap)
							v, err := dvevaluation.IsDefined(subBlock[subStart:subBlockLen], scope, internalParsingInfo[0].row, internalParsingInfo[0].column+subStart, place, dvevaluation.EVALUATE_OPTION_UNDEFINED)
							if err != nil {
								configInfo.Err = err
							}
							bIfCondition = v
						}
					}
				} else if b2 == 'n' { //include
					if b3 == 'c' && subStart+4 < subBlockLen && subBlock[subStart] == 'l' && subBlock[subStart+1] == 'u' && subBlock[subStart+2] == 'd' && subBlock[subStart+3] == 'e' && subBlock[subStart+4] <= 32 {
						isComment = false
						subStart += 4
						for subStart < subBlockLen && subBlock[subStart] <= 32 {
							subStart++
						}
						if subStart < subBlockLen {
							kind := -1
							b1 = subBlock[subStart]
							b2 = subBlock[subBlockLen-1]
							if b1 == b2 && b1 == '"' {
								kind = 0
								subStart++
								subBlockLen--
							} else if b1 == '<' && b2 == '>' {
								subStart++
								subBlockLen--
								kind = 1
							}
							if subStart < subBlockLen {
								LinearSmartConfigFromFile(string(subBlock[subStart:subBlockLen]), configInfo, kind)
							} else {
								if DvParserLog {
									log.Print("include instruction is empty" + getWhereInfo(subStart, place, internalParsingInfo))
								}
							}
						} else {
							if DvParserLog {
								log.Print("include instruction is empty" + getWhereInfo(subStart, place, internalParsingInfo))
							}
						}
					}
				}
			} else if b1 == 'e' {
				if b2 == 'l' {
					n := len(ifElseInfo)
					if n == 0 {
						bIfCondition = -(IFELSE_ELSE_ERROR + IFELSE_NONEWIF)
					} else {
						bIfCondition = -(ifElseInfo[n-1] + IFELSE_NONEWIF)
					}
					if b3 == 's' { //else
						if subStart < subBlockLen && subBlock[subStart] == 'e' && (subStart+1 >= subBlockLen || subBlock[subStart+1] <= 32) {
							bIfCondition = -bIfCondition
							if bIfCondition <= IFELSE_NONEWIF+IFELSE_ELSE_RUN {
								ifElseInfo[n-1] = IFELSE_ELSE_NOMORE
							}
						}
					} else if b3 == 'i' { //elif...
						if subStart < subBlockLen && subBlock[subStart] == 'f' {
							if subStart+1 >= subBlockLen || subBlock[subStart+1] <= 32 {
								//elif
								bIfCondition = -bIfCondition
								if bIfCondition == IFELSE_NONEWIF+IFELSE_ELSE_RUN {

									res := dvevaluation.EvalAsBoolean(subBlock[subStart+1:subBlockLen], configInfo.InputMap, configInfo.OutputMap, internalParsingInfo[0].row, internalParsingInfo[0].column, place)
									if res.Err != nil {
										configInfo.Err = res.Err
										return
									}
									if !res.FinalResult {
										bIfCondition = IFELSE_NONEWIF + IFELSE_ELSE_NOTRUN
									}
								}
							}
							if subStart+3 < subBlockLen && subBlock[subStart+1] == 'd' && subBlock[subStart+2] == 'e' && subBlock[subStart+3] == 'f' && (subStart+4 >= subBlockLen || subBlock[subStart+4] <= 32) {
								//elifdef
								bIfCondition = -bIfCondition
								if bIfCondition == IFELSE_NONEWIF+IFELSE_ELSE_RUN {
									scope := dvevaluation.NewDvObjectFrom2Maps(configInfo.OutputMap, configInfo.InputMap)
									v, err := dvevaluation.IsDefined(subBlock[subStart+4:subBlockLen], scope, row, col, place, 0)
									if err != nil {
										configInfo.Err = err
										return
									}
									if v == 0 {
										bIfCondition = IFELSE_NONEWIF + IFELSE_ELSE_NOTRUN
									}
								}
							}
							if subStart+4 < subBlockLen && subBlock[subStart+1] == 'n' && subBlock[subStart+2] == 'd' && subBlock[subStart+3] == 'e' && subBlock[subStart+4] == 'f' && (subStart+5 >= subBlockLen || subBlock[subStart+5] <= 32) {
								//elifndef
								bIfCondition = -bIfCondition
								if bIfCondition == IFELSE_NONEWIF+IFELSE_ELSE_RUN {
									scope := dvevaluation.NewDvObjectFrom2Maps(configInfo.OutputMap, configInfo.InputMap)
									v, err := dvevaluation.IsDefined(subBlock[subStart+5:subBlockLen], scope, row, col, place, dvevaluation.EVALUATE_OPTION_UNDEFINED)
									if err != nil {
										configInfo.Err = err
										return
									}
									if v == 0 {
										bIfCondition = IFELSE_NONEWIF + IFELSE_ELSE_NOTRUN
									}
								}
							}
							if bIfCondition == IFELSE_NONEWIF+IFELSE_ELSE_RUN {
								ifElseInfo[n-1] = IFELSE_ELSE_NOTRUN
							}
						}
					}
					if bIfCondition > IFELSE_NONEWIF+IFELSE_ELSE_RUN {
						if bIfCondition == IFELSE_NONEWIF+IFELSE_ELSE_ERROR {
							configInfo.Err = errors.New("else/elif/elifdef/elifndef directive is outside of if block" + getWhereInfo(subStart, place, internalParsingInfo))
						} else {
							configInfo.Err = errors.New("endif expected but else/elif/elifdef/elifndef found" + getWhereInfo(subStart, place, internalParsingInfo))
						}
						return
					}
				} else if b2 == 'n' { //endif
					if b3 == 'd' && subStart+1 < subBlockLen && subBlock[subStart] == 'i' && subBlock[subStart+1] == 'f' && (subStart+2 >= subBlockLen || subBlock[subStart+2] <= 32) {
						isComment = false
						n := len(ifElseInfo)
						if n == 0 {
							configInfo.Err = errors.New("Unexpected endif directive" + getWhereInfo(subStart, place, internalParsingInfo))
							return
						}
						ifElseInfo = ifElseInfo[:n-1]
					}
				} else if b2 == 'r' { //error
					if b3 == 'r' && subStart+1 < subBlockLen && subBlock[subStart] == 'o' && subBlock[subStart+1] == 'r' && (subStart+2 >= subBlockLen || subBlock[subStart+2] <= 32) {
						isComment = false
						subStart += 3
						for subStart < subBlockLen && subBlock[subStart] <= 32 {
							subStart++
						}
						errorMessage := "Error"
						if subStart < subBlockLen {
							errorMessage = string(subBlock[subStart:subBlockLen])
						}
						configInfo.Err = errors.New(errorMessage + " " + getWhereInfo(subStart, place, internalParsingInfo))
						return
					}
				}
			} else if b1 == 'd' { //define
				if b2 == 'e' && b3 == 'f' && subStart+4 < subBlockLen && subBlock[subStart] == 'i' && subBlock[subStart+1] == 'n' && subBlock[subStart+2] == 'e' && subBlock[subStart+3] <= 32 {
					isComment = false
					subStart += 4
					for subStart < subBlockLen && subBlock[subStart] <= 32 {
						subStart++
					}
					pos1 := subStart
					for pos1 < subBlockLen && subBlock[pos1] > 32 {
						pos1++
					}
					pos2 := pos1
					for pos2 < subBlockLen && subBlock[pos2] <= 32 {
						pos2++
					}
					if subStart < pos1 {
						keyDefined := string(subBlock[subStart:pos1])
						valueDefined := string(subBlock[pos2:subBlockLen])
						configInfo.OutputMap[keyDefined] = valueDefined
					}
				}
			} else if b1 == 'u' { //undef
				if b2 == 'n' && b3 == 'd' && subStart+3 < subBlockLen && subBlock[subStart] == 'e' && subBlock[subStart+1] == 'f' && subBlock[subStart+2] <= 32 {
					isComment = false
					subStart += 3
					for subStart < subBlockLen && subBlock[subStart] < 32 {
						subStart++
					}
					if subStart < subBlockLen {
						keyUndefined := string(subBlock[subStart:subBlockLen])
						if _, keyOk := configInfo.OutputMap[keyUndefined]; keyOk {
							delete(configInfo.OutputMap, keyUndefined)
						}
					}
				}
			}
			if bIfCondition >= 0 {
				bIfRunCondition := bIfCondition & IFELSE_LOW_MASK
				if (bIfCondition & IFELSE_NONEWIF) == 0 {
					ifElseInfo = append(ifElseInfo, 1-bIfRunCondition)
				}
				if bIfRunCondition == 0 {
					//consumption
					for i < l {
						for ; i < l && data[i] <= 32; i++ {
							if data[i] == 13 || data[i] == 10 {
								if data[i] == 13 || i > 0 && data[i-1] != 13 {
									row++
									col = 1
								}
							} else {
								col++
							}
						}
						if i == l {
							break
						}
						if data[i] == '#' && i+5 < l {
							b1 := data[i+1]
							b2 := data[i+2]
							if b1 == 'i' && b2 == 'f' {
								if data[i+3] <= 32 || data[i+3] == 'd' && i+6 < l && data[i+4] == 'e' && data[i+5] == 'f' && data[i+6] <= 32 || data[i+3] == 'n' && i+7 < l && data[i+4] == 'd' && data[i+5] == 'e' && data[i+6] == 'f' && data[i+7] <= 32 {
									ifElseInfo = append(ifElseInfo, IFELSE_CONSUMED)
								}
							} else if b1 == 'e' {
								if b2 == 'l' && ifElseInfo[len(ifElseInfo)-1] == IFELSE_ELSE_RUN {
									if data[i+3] == 's' && data[i+4] == 'e' && data[i+5] <= 32 {
										break
									}
									if data[i+3] == 'i' && data[i+4] == 'f' {
										if data[i+5] <= 32 || data[i+5] == 'd' && i+8 < l && data[i+6] == 'e' && data[i+7] == 'f' && data[i+8] <= 32 || data[i+5] == 'n' && i+9 < l && data[i+6] == 'd' && data[i+7] == 'e' && data[i+8] == 'f' && data[i+9] <= 32 {
											break
										}
									}
								} else if b2 == 'n' && data[i+3] == 'd' && data[i+4] == 'i' && data[i+5] == 'f' && (i+6 == l || data[i+6] <= 32) {
									ifElseInfo = ifElseInfo[:len(ifElseInfo)-1]
									if len(ifElseInfo) == 0 || ifElseInfo[len(ifElseInfo)-1] != IFELSE_CONSUMED {
										for ; i < l && data[i] != 13 && data[i] != 10; i++ {
										}
										break
									}
								}
							}
						}
						for ; i < l && data[i] != 13 && data[i] != 10; i++ {
						}
					}
				}
			} else if isComment && preserveSpace {
				configInfo.OutputLines = append(configInfo.OutputLines, currentLine)
			}
		} else {
			if isDictionary {
				subBlock = subBlock[subStart:subBlockLen]
				subBlockLen -= subStart
				if configIsEscaped {
					pos = 0
					pnt := 0
					bEnd := subBlock[subBlockLen-1]
					if subBlockLen > 2 && subBlock[subBlockLen-2] == '\\' && bEnd != '`' {
						bEnd = 0
					}
					for pos = 0; pos < subBlockLen; pos++ {
						b := subBlock[pos]
						if b == '\\' {
							if pos < subBlockLen-1 {
								pos++
								subBlock[pnt] = subBlock[pos]
								pnt++
							}
						} else if b == '=' {
							break
						} else {
							subBlock[pnt] = b
							pnt++
						}
					}
					if pos == 0 || pos >= subBlockLen {
						if DvParserLog {
							log.Print("key=value expected for " + string(subBlock) + getWhereInfo(subStart, place, internalParsingInfo))
						}
					} else {
						for pnt > 0 && subBlock[pnt-1] <= 32 {
							pnt--
						}
						key := string(subBlock[:pnt])
						pos++
						for pos < subBlockLen && subBlock[pos] <= 32 {
							pos++
						}
						pnt = 0
						if pos < subBlockLen-1 && subBlock[pos] == bEnd && bEnd == '`' {
							val := string(subBlock[pos+1 : subBlockLen-1])
							configInfo.OutputMap[key] = val
						} else {
							if pos < subBlockLen-1 && subBlock[pos] == bEnd && bEnd == '"' {
								pos++
								subBlockLen--
							}
							pnt = 0
							for pnt = 0; pos < subBlockLen; pos++ {
								b := subBlock[pos]
								if b == '\\' {
									if pos < subBlockLen-1 {
										pos++
										subBlock[pnt] = subBlock[pos]
										pnt++
									}
								} else {
									subBlock[pnt] = b
									pnt++
								}
							}
							val := string(subBlock[:pnt])
							configInfo.OutputMap[key] = val
						}
					}
				} else {
					pnt := bytes.IndexByte(subBlock, '=')
					if pnt <= 0 {
						if DvParserLog {
							log.Print("key=value expected for " + string(subBlock) + getWhereInfo(subStart, place, internalParsingInfo))
						}
					} else {
						pntEnd := pnt
						for pntEnd > 0 && subBlock[pntEnd-1] <= 32 {
							pntEnd--
						}
						key := string(subBlock[:pntEnd])
						pnt++
						bEnd := subBlock[subBlockLen-1]
						for pnt < subBlockLen && subBlock[pnt] <= 32 {
							pnt++
						}
						if pnt < subBlockLen-1 && subBlock[pnt] == bEnd && (bEnd == '"' || bEnd == '`') {
							pnt++
							subBlockLen--
						}
						val := string(subBlock[pnt:subBlockLen])
						configInfo.OutputMap[key] = val
					}
				}
			} else {
				if preserveSpace {
					configInfo.OutputLines = append(configInfo.OutputLines, currentLine)
				} else {
					configInfo.OutputLines = append(configInfo.OutputLines, subBlock[subStart:subBlockLen])
				}
			}
		}
	}
	if len(ifElseInfo) != 0 {
		if DvParserLog {
			log.Print("If-else-endif directives are not closed by " + strconv.Itoa(len(ifElseInfo)) + " levels in " + place)
		}
		configInfo.Err = errors.New("If-Else-Endif directives are not closed in " + place)
	}
}

func LinearSmartConfigParse(data []byte, configInfo *ConfigInfo, sourceName string) {
	if configInfo.insider == "" {
		if configInfo.NumberOfBrackets <= 0 {
			configInfo.NumberOfBrackets = 3
		}
		if (configInfo.Options&CONFIG_DISABLE_INITIAL_MAP_FROM_ENVIRONMENT) == 0 && configInfo.InputMap == nil {
			configInfo.InputMap = LinearSmartConfigFromEnvironment()
		}
		if (configInfo.Options & CONFIG_RESULT_NOT_CONTAINING_PREMAP) == 0 {
			configInfo.OutputMap = configInfo.InputMap
		} else {
			configInfo.OutputMap = make(map[string]string)
		}
		if (configInfo.Options & CONFIG_IS_NOT_VARIABLES) != 0 {
			configInfo.OutputLines = make([][]byte, 0, 1024)
		}
	}
	if len(data) > 0 {
		configInfo.insider = sourceName
		linearSmartConfigParse_internal(data, configInfo, sourceName)
	}
}

/**************************************************************************
level = -1 if even missing the file is not a problem
      >= 0 start looking the file from paths configInfo.FilePaths[level] and up
***************************************************************************/
func LinearSmartConfigFromFile(fileName string, configInfo *ConfigInfo, level int) {
	insider := configInfo.insider
	data, absFileName := linearSmartConfigParse_readFile(fileName, configInfo, level)
	if configInfo.Err == nil {
		LinearSmartConfigParse(data, configInfo, absFileName)
	}
	configInfo.insider = insider
}

func LinearSmartConfigFromEnvironment() map[string]string {
	mappa := make(map[string]string)
	for k, v := range globalPremap {
		mappa[k] = v
	}
	return mappa
}

func ReadPropertiesFileWithEnvironmentVariables(currentDir string, propertiesName string, setFilePaths func() error) error {
	DvParserLog = true
	if currentDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			dir = "."
		}
		currentDir = dir
	}
	dvlog.SetCurrentNamespace(currentDir)
	GlobalProperties = LinearSmartConfigFromEnvironment()
	initializeRegisteredFunctions()
	GeneralFilePaths = make([]string, 1, 4)
	GeneralFilePaths[0] = currentDir
	err := setFilePaths()
	if err != nil {
		return err
	}
	filename := FindInGeneralPaths(propertiesName)
	if filename == "" {
		return nil
	}
	configInfo := ConfigInfo{InputMap: GlobalProperties, FilePaths: GeneralFilePaths}
	LinearSmartConfigFromFile(filename, &configInfo, -1)
	if configInfo.Err != nil {
		return errors.New("Error: cannot read properties " + filename + ": " + configInfo.Err.Error())
	}
	return setFilePaths()
}

func defaultSetFilePath() error {
	return nil
}

func ConvertStringByGlobalProperties(data string, sourceName string) (string, error) {
	return ConvertByteArrayByGlobalProperties([]byte(data), sourceName)
}

func ConvertByteArrayByGlobalProperties(data []byte, sourceName string) (res string, err error) {
	var buf []byte
	buf, err = ConvertByteArrayBySpecificProperties(data, sourceName, GlobalProperties, NumberOfBracketsInConfigParsing, CONFIG_PRESERVE_SPACE)
	return string(buf), err
}

func ConvertByteArrayByGlobalPropertiesRuntime(data []byte, sourceName string) (res string, err error) {
	var buf []byte
	buf, err = ConvertByteArrayBySpecificProperties(data, sourceName, GlobalProperties, 2, 0)
	return string(buf), err
}

func ConvertByteArrayBySpecificPropertiesInLines(data []byte, sourceName string, properties map[string]string, numberOfBrackets int, configOptions int) ([][]byte, error) {

	configInfo := &ConfigInfo{
		NumberOfBrackets: numberOfBrackets,
		InputMap:         properties,
		Options:          CONFIG_IS_NOT_VARIABLES | configOptions,
		OutputMap:        make(map[string]string),
		OutputLines:      nil,
		Err:              nil,
		FilePaths:        []string{"."},
	}
	LinearSmartConfigParse(data, configInfo, sourceName)
	if configInfo.Err != nil {
		return nil, configInfo.Err
	}
	return configInfo.OutputLines, nil
}

func ConvertByteArrayBySpecificProperties(data []byte, sourceName string, properties map[string]string, numberOfBrackets int, configOptions int) ([]byte, error) {
	res, err := ConvertByteArrayBySpecificPropertiesInLines(data, sourceName, properties, numberOfBrackets, configOptions)
	return bytes.Join(res, []byte{byte(10)}), err
}

func ConvertByteArrayByGlobalPropertiesInByteLines(data []byte, sourceName string) (res [][]byte, err error) {
	return ConvertByteArrayBySpecificPropertiesInLines(data, sourceName, GlobalProperties, NumberOfBracketsInConfigParsing, 0)
}

func ConvertByteArrayByGlobalPropertiesInStringLines(data []byte, sourceName string) (res []string, err error) {
	resInBytes, err := ConvertByteArrayByGlobalPropertiesInByteLines(data, sourceName)
	if err != nil {
		return
	}
	n := len(resInBytes)
	res = make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = string(resInBytes[i])
	}
	return
}

func SetNumberOfBracketsInConfigParsing(n int) {
	NumberOfBracketsInConfigParsing = n
}

func JoinStringsWithClean(data []string) []byte {
	res := make([]byte, 0, 1024)
	pos := 0
	n := len(data)
	for i := 0; i < n; i++ {
		if i > 0 {
			res = append(res, byte(10))
			pos++
		}
		res = append(res, []byte(data[i])...)
		lastPos := len(res)
		for ; pos < lastPos; pos++ {
			if res[pos] == 10 || res[pos] == 13 {
				res[pos] = ' '
			}
		}
	}
	return res
}

func ConvertStringArrayByGlobalProperties(data []string, sourceName string) (res []string, err error) {
	if sourceName == "" {
		sourceName = "."
	}
	res = make([]string, len(data))
	n := len(data)
	for i := 0; i < n; i++ {
		if len(data[i]) > 0 {
			converted, err1 := ConvertByteArrayByGlobalProperties([]byte(data[i]), sourceName)
			if err1 != nil {
				return nil, err1
			}
			res[i] = converted
		}
	}
	return
}

func ConvertFileByGlobalProperties(fileName string) ([]byte, error) {
	return SmartReadTemplate(fileName, 3, 10)
}
