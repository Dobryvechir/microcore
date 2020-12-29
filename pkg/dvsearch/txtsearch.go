// package dvsearch provides functions to search in text
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvsearch

/*********************************************************************
literal \ must be escaped as \\
literal / must be escaped as \057 (=0x2f)
literal * must be escaped as \052 (=0x2a)
Simple search/replace: Hope Hoffnung
Escaped search/replace: \040hope(\042Love\042)\n\r
Simple wildcard: hope(\042*\042)*;
Simple regular expressions in search: hope(\042/[A-Z]* /\042)/\s/;
Simple replacement places in replace: Eternal life with $1 is from $2
********************************************************************/

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type SearchOptionRunBuffer struct {
	pos     int
	size    int
	restPos int
	restLen int
}

type SearchOption struct {
	Pattern         string
	Search          []byte
	SearchLen       int
	SearchRest      [][]byte
	SearchPatterns  []string
	RestPosNumber   int
	ReplaceSkip     int
	ReplaceLimit    int
	Replace         []byte
	ReplaceTotalLen int
	ReplaceRest     [][]byte
	ReplaceOrder    []int
	Verbose         bool
	Check           bool
	TrimSpace       int
	ReadOnly        bool
	runBuffer       []SearchOptionRunBuffer
	saveBuffer      []byte
	baseName        []byte
	wholeName       []byte
	currentRestPos  int
	restInfo        []int
}

const (
	TrimSpaceLeft  = 1
	TrimSpaceRight = 2
)

func addRestInfo(option *SearchOption, pos int, size int) {
	n := option.currentRestPos + 2
	m := len(option.restInfo)
	if n > m {
		buf := make([]int, m+1024)
		for i := 0; i < m; i++ {
			buf[i] = option.restInfo[i]
		}
		option.restInfo = buf
	}
	option.restInfo[n-2] = pos
	option.restInfo[n-1] = size
	option.currentRestPos = n
}

func calculatePositionAsLineColumn(data []byte, pos int) string {
	row := 1
	col := 1
	if pos > len(data) {
		pos = len(data)
	}
	for i := 0; i < pos; i++ {
		c := data[i]
		if c == 13 || (c == 10 && (i == 0 || data[i-1] != 13)) {
			row++
			col = 1
		} else if c != 10 {
			col++
		}
	}
	return strconv.Itoa(row) + ":" + strconv.Itoa(col)
}

func logChanges(origBuf []byte, origStart int, origEnd int, changeBuf []byte, changeStart int, changeEnd int, name []byte) {
	message := "In " + string(name) + " " + calculatePositionAsLineColumn(origBuf, origStart) + " - " + calculatePositionAsLineColumn(origBuf, origEnd) + " into " + calculatePositionAsLineColumn(changeBuf, changeStart) + " - " + calculatePositionAsLineColumn(changeBuf, changeEnd) + "\n"
	preStart := origStart - 15
	if preStart < 0 {
		preStart = 0
	}
	postEnd := origEnd + 15
	if postEnd > len(origBuf) {
		postEnd = len(origBuf)
	}
	message += string(origBuf[preStart:origStart]) + "[[[" + string(origBuf[origStart:origEnd]) + "]]]" + string(origBuf[origEnd:postEnd])
	message += "\n----[[" + string(changeBuf[changeStart:changeEnd]) + "]]----\n"
	log.Println(message)
}

func extendInBuf(size int, data []byte, option *SearchOption, count int) []byte {
	buf := option.saveBuffer
	if len(buf) < size {
		buf = make([]byte, size+1024)
		option.saveBuffer = buf
	}
	dstPos := 0
	srcPos := 0
	n := len(data)
	r := option.Replace
	nr := len(r)
	for j := option.ReplaceSkip; j < count; j++ {
		srcLimit := option.runBuffer[j].pos
		for ; srcPos < srcLimit; srcPos++ {
			buf[dstPos] = data[srcPos]
			dstPos++
		}
		srcPos += option.runBuffer[j].size
		oldDstPos := dstPos
		for i := 0; i < nr; i++ {
			buf[dstPos] = r[i]
			dstPos++
		}
		restPos := option.runBuffer[j].restPos
		restLen := option.runBuffer[j].restLen
		m := len(option.ReplaceRest)
		for k := 0; k < m; k++ {
			ord := option.ReplaceOrder[k]
			var bf []byte = nil
			if ord == 0 {
				bf = option.baseName
			} else if ord >= 1 && ord <= restLen {
				ord = (ord-1)*2 + restPos
				bf = data[option.restInfo[ord] : option.restInfo[ord]+option.restInfo[ord+1]]
			}
			ord = len(bf)
			for l := 0; l < ord; l++ {
				buf[dstPos] = bf[l]
				dstPos++
			}
			bf = option.ReplaceRest[k]
			ord = len(bf)
			for l := 0; l < ord; l++ {
				buf[dstPos] = bf[l]
				dstPos++
			}
		}
		if option.Check || option.Verbose {
			logChanges(data, srcPos-option.runBuffer[j].size, srcPos, buf, oldDstPos, dstPos, option.wholeName)
		}
	}
	for ; srcPos < n; srcPos++ {
		buf[dstPos] = data[srcPos]
		dstPos++
	}
	return buf[:dstPos]
}

func countInsertionSize(option *SearchOption, count int) int {
	restPos := option.runBuffer[count].restPos
	restLen := option.runBuffer[count].restLen
	n := len(option.ReplaceRest)
	size := 0
	for i := 0; i < n; i++ {
		ord := option.ReplaceOrder[i]
		if ord == 0 {
			size += len(option.baseName)
		} else if ord >= 1 && ord <= restLen {
			ord = (ord-1)*2 + restPos
			size += option.restInfo[ord+1]
		}
	}
	return size
}

func SearchProcessFile(name string, option *SearchOption) int {
	option.wholeName = []byte(name)
	data, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Printf("! %s %v", name, err)
		return 0
	}
	n := len(data)
	size := n
	count := 0
	pos := 0
	option.currentRestPos = 0
searchProcessFileLoopMain:
	for pos < n {
		blockPos := bytes.Index(data[pos:], option.Search) + pos
		if blockPos < pos {
			break
		}
		if len(option.runBuffer) < count+1 {
			alterPuffer := option.runBuffer
			n := len(alterPuffer)
			buffer := make([]SearchOptionRunBuffer, n+1024)
			for i := 0; i < n; i++ {
				buffer[i].pos = alterPuffer[i].pos
				buffer[i].size = alterPuffer[i].size
				buffer[i].restPos = alterPuffer[i].restPos
				buffer[i].restLen = alterPuffer[i].restLen
			}
			option.runBuffer = buffer
		}
		blockSize := option.SearchLen
		restPos := option.currentRestPos
		pos = blockPos + blockSize
		m := len(option.SearchRest)
		if m > 0 {
			for i := 0; i < m; i++ {
				bf := option.SearchRest[i]
				bflen := len(bf)
				if bflen == 0 {
					addRestInfo(option, pos, 0)
				} else {
					blockEnd := bytes.Index(data[pos:], option.SearchRest[i]) + pos
					if blockEnd < pos {
						break searchProcessFileLoopMain
					}
					addRestInfo(option, pos, blockEnd-pos)
					pos = blockEnd + bflen
				}
			}
			blockSize = pos - blockPos
		}
		//TODO process search for complicated cases
		/*********************************************************************************************
		                      CONSIDER ALL ALGORITHM IN CLOSE CONNECTION WITH ALREADY CREATED STRUCTURE
				      Algorythm:
				      1. Write a function
				        func processSequentialRegularExpressions(regs []string???, data []byte, posStart int, posEnd int, fixedStart bool, fixedEnd bool) (ok bool, finalStart int, finalEnd int,res [][]byte) {
		                        }
		                        which looks for sequential regular expressions taking place
		                      2. 1) make places for known values
		                         2) inside make calls processSequentialRegularExpressions()
		                         3) make another positions if previous fails, until all places are passed
		                         4) make further search only after the position where the
		                         5) update common size for good
		                      3. Introduce options as follows:
		                         **XXX - at least XXX characters before the next value
				         ***XXX to search for HTML-like tag and its whole enclosure (number indicates 0 - current, 1 - first inside, 2 - second...-1 - external by 1 level, -2 external by 2 levels)
		                         **** subchild search
		                         *****X find whole quote and X is a kind info
		                      3. Add an option to search for whole line
		                      Transfer the code to the dvserver and make it work for replacement functionality

				**********************************************************************************************/
		if (option.TrimSpace & TrimSpaceLeft) != 0 {
			for blockPos > 0 && data[blockPos-1] <= ' ' {
				blockPos--
				blockSize++
			}
		}
		pos = blockPos + blockSize
		if (option.TrimSpace & TrimSpaceRight) != 0 {
			for pos < n && data[pos] <= ' ' {
				pos++
			}
			blockSize = pos - blockPos
		}
		option.runBuffer[count].pos = blockPos
		option.runBuffer[count].size = blockSize
		option.runBuffer[count].restPos = restPos
		option.runBuffer[count].restLen = (option.currentRestPos - restPos) >> 1
		size += option.ReplaceTotalLen - blockSize + countInsertionSize(option, count)
		count++
		if option.ReplaceLimit > 0 && count >= option.ReplaceSkip+option.ReplaceLimit {
			break
		}
	}
	if count > 0 {
		data = extendInBuf(size, data, option, count)
		if option.Check {
			name += ".check"
		}
		if !option.ReadOnly {
			err = ioutil.WriteFile(name, data, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return count
}

func SearchProcessDir(dir string, option *SearchOption) (int, int) {
	if dir == "" {
		dir = "./"
	} else {
		c := dir[len(dir)-1]
		if c != '/' && c != '\\' {
			dir += "/"
		}
	}
	d, err := os.Open(dir)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	defer d.Close()
	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	if option.Verbose {
		fmt.Printf("Dir %s has %d entries\n", dir, len(files))
	}
	all := 0
	found := 0
	for _, file := range files {
		nm := file.Name()
		if nm == "." || nm == ".." {
			continue
		}
		switch mode := file.Mode(); {
		case mode.IsRegular():
			matched, err := filepath.Match(option.Pattern, nm)
			if err != nil {
				fmt.Println(err)
				return 0, 0
			} else if matched {
				name := dir + nm
				all++
				option.baseName = []byte(nm)
				hasSearch := SearchProcessFile(name, option)
				if hasSearch > 0 {
					found++
					if option.Verbose {
						fmt.Printf("%d %s\n", hasSearch, name)
					}
				} else if option.Verbose {
					fmt.Println("- " + nm)
				}
			}
		case mode.IsDir():
			newAll, newFound := SearchProcessDir(dir+nm, option)
			all += newAll
			found += newFound
		}
	}
	return all, found
}

func convertSearchReplacePattern(pattern string) []byte {
	n := len(pattern)
	res := []byte(pattern)
	resPos := 0
	for i := 0; i < n; i++ {
		c := res[i]
		if c == '\\' && i+1 < n {
			i++
			c = res[i]
			if c == 'n' {
				res[resPos] = byte(13)
			} else if c == 'r' {
				res[resPos] = byte(10)
			} else if c == 't' {
				res[resPos] = byte(9)
			} else if c >= '0' && c <= '9' && i+2 < n && res[i+1] >= '0' && res[i+1] <= '9' && res[i+2] >= '0' && res[i+2] <= '9' {
				res[resPos] = ((c - '0') << 6) | ((res[i+1] - '0') << 3) | (res[i+2] - '0')
				i += 2
			} else {
				res[resPos] = c
			}
			resPos++
		} else {
			res[resPos] = c
			resPos++
		}
	}
	return res[:resPos]
}

func analyzeSearchPattern(search string) (first []byte, rest [][]byte, patterns []string) {
	p := strings.Index(search, "/")
	if p < 0 {
		s := strings.Split(search, "*")
		n := len(s)
		first = convertSearchReplacePattern(s[0])
		n--
		if n > 0 {
			patterns = make([]string, n)
			rest = make([][]byte, n)
			for i := 0; i < n; i++ {
				rest[i] = convertSearchReplacePattern(s[i+1])
			}
		}
		return
	}
	p1 := strings.Index(search, "*")
	if p1 >= 0 && p1 < p {
		p = p1
	}
	first = convertSearchReplacePattern(search[:p])
	patterns = make([]string, 0, 7)
	rest = make([][]byte, 0, 7)
	search = search[p:]
	for len(search) > 0 {
		switch search[0] {
		case '*':
			for len(search) > 0 && search[0] == '*' {
				search = search[1:]
			}
			patterns = append(patterns, "")
		case '/':
			p := strings.Index(search[1:], "/")
			if p < 0 {
				fmt.Println("Incorrect structure in " + search)
				p = len(search)
			}
			patterns = append(patterns, search[1:p])
			if p+1 >= len(search) {
				search = ""
			} else {
				search = search[p+1:]
			}
		}
		p = strings.Index(search, "/")
		p1 = strings.Index(search, "*")
		if p1 >= 0 && (p1 < p || p < 0) {
			p = p1
		}
		if p < 0 {
			p = len(search)
		}
		rest = append(rest, convertSearchReplacePattern(search[:p]))
		search = search[p:]
	}
	return
}

func findFirstReplacement(replace string) (pos int, nextPos int, number int) {
	pos = 0
	n := len(replace)
	for pos < n {
		i := strings.Index(replace[pos:], "$") + pos
		if i < pos || i+1 >= n {
			break
		}
		pos = i + 1
		c := replace[pos]
		if c >= '0' && c <= '9' {
			nextPos = pos + 1
			number = int(c - '0')
			pos--
			for nextPos < n && replace[nextPos] >= '0' && replace[nextPos] <= '9' {
				number = number*10 + int(replace[nextPos]-'0')
				nextPos++
			}
			return
		}
	}
	return n, -1, -1
}

func analyzeReplacePattern(replace string) (first []byte, rest [][]byte, order []int) {
	pos, nextPos, number := findFirstReplacement(replace)
	first = convertSearchReplacePattern(replace[:pos])
	if number < 0 {
		return
	}
	rest = make([][]byte, 0, 7)
	order = make([]int, 0, 7)
	for number >= 0 {
		order = append(order, number)
		replace = replace[nextPos:]
		pos, nextPos, number = findFirstReplacement(replace)
		rest = append(rest, convertSearchReplacePattern(replace[:pos]))
	}
	return
}

func readNumberInOptionsAfter(options string, letter string, defValue int) int {
	lpos := strings.Index(options, letter)
	n := len(options)
	if lpos >= 0 && lpos+1 < n && options[lpos+1] >= '0' && options[lpos+1] <= '9' {
		defValue = int(options[lpos+1] - '0')
		lpos += 2
		for lpos < n && options[lpos] >= '0' && options[lpos] <= '9' {
			defValue = defValue*10 + int(options[lpos]-'0')
		}
	}
	return defValue
}

func GenerateSearchOptions(search string, replace string, options string, pattern string) *SearchOption {
	spaceTrim := 0
	spos := strings.Index(options, "t")
	if spos >= 0 {
		next := byte('B')
		if spos+1 < len(options) {
			next = options[spos+1]
		}
		if next == 'L' {
			spaceTrim = TrimSpaceLeft
		} else if next == 'R' {
			spaceTrim = TrimSpaceRight
		} else {
			spaceTrim = TrimSpaceLeft | TrimSpaceRight
		}
	}
	skip := readNumberInOptionsAfter(options, "s", 0)
	limit := readNumberInOptionsAfter(options, "l", 0)
	searchFirst, searchRest, searchPatterns := analyzeSearchPattern(search)
	rests := len(searchRest)
	replaceFirst, replaceRest, replaceOrder := analyzeReplacePattern(replace)
	replaceTotalLen := len(replaceFirst)
	n := len(replaceRest)
	for i := 0; i < n; i++ {
		replaceTotalLen += len(replaceRest[i])
	}
	searchOptions := &SearchOption{
		Pattern:         pattern,
		Search:          searchFirst,
		SearchLen:       len(searchFirst),
		SearchRest:      searchRest,
		SearchPatterns:  searchPatterns,
		RestPosNumber:   rests,
		Replace:         replaceFirst,
		ReplaceTotalLen: replaceTotalLen,
		ReplaceRest:     replaceRest,
		ReplaceOrder:    replaceOrder,
		ReplaceSkip:     skip,
		ReplaceLimit:    limit,
		Verbose:         strings.Index(options, "v") >= 0,
		Check:           strings.Index(options, "c") >= 0,
		ReadOnly:        strings.Index(options, "r") >= 0,
		TrimSpace:       spaceTrim,
		saveBuffer:      make([]byte, 65536),
	}
	return searchOptions
}
