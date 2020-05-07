/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvcom

import (
	"bufio"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type ReadLinesPool struct {
	lines    []string
	err      error
	fileName string
}

func readLinesFromFile(fileName string) ReadLinesPool {
	pool := ReadLinesPool{fileName: fileName, err: nil, lines: make([]string, 0, 100)}
	file, err := os.Open(fileName)
	if err != nil {
		pool.err = err
		return pool
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pool.lines = append(pool.lines, scanner.Text())
	}

	pool.err = scanner.Err()
	return pool
}

func writeLinesToFile(pool *ReadLinesPool) {
	eol := EOL_BYTES
	eolLen := len(eol)
	count := len(pool.lines) * eolLen
	for _, s := range pool.lines {
		count += len(s)
	}
	buf := make([]byte, 0, count)
	for _, s := range pool.lines {
		buf = append(buf, []byte(s)...)
		buf = append(buf, eol...)
	}
	pool.err = ioutil.WriteFile(pool.fileName, buf, 0644)
	if pool.err != nil {
		name := strings.ToLower(pool.err.Error())
		p := strings.Index(name, "access")
		p1 := strings.Index(name[p+1:], "denied")
		if p >= 0 && p1 > 0 {
			tmpName := dvlog.GetTemporaryFileName()
			pool.err = ioutil.WriteFile(tmpName, buf, 0644)
			if LogHosts && dvlog.CurrentLogLevel >= dvlog.LogDetail {
				log.Printf("Writing file %d bytes to %s error:%v", len(buf), pool.fileName, pool.err)
			}
			if pool.err == nil {
				AddAdministrativeTask("copy", []string{tmpName, pool.fileName})
				if LogHosts && dvlog.CurrentLogLevel >= dvlog.LogWarning {
					log.Printf("Administrative task will run to write to %s", pool.fileName)
				}
			}
		}
	}
}
