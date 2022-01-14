/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdir

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CopyWholeFolder(src string, dest string) error {
	src = GetFolderNameWithoutLastSlash(src)
	dest = GetFolderNameWithoutLastSlash(dest)
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyWholeFolder(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := CopySingleFile(sourcePath, destPath); err != nil {
				return err
			}
		}

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, entry.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func CopySingleFile(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	defer out.Close()
	if err != nil {
		return err
	}

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CreateOrCleanDir(dir string, perm os.FileMode) error {
	if Exists(dir) {
		os.RemoveAll(dir)
	}
	return CreateIfNotExists(dir, perm)
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

func DeleteFileIfExists(fileName string) bool {
	dat, err := os.Stat(fileName)
	if err == nil {
		if dat.IsDir() {
			err = os.RemoveAll(fileName)
		} else {
			err = os.Remove(fileName)
		}
		if err != nil {
			return false
		}
	}
	return true
}

func DeleteFilesIfExist(fileNames []string) int {
	k := 0
	n := len(fileNames)
	for i := 0; i < n; i++ {
		if DeleteFileIfExists(fileNames[i]) {
			k++
		}
	}
	return k
}

func MakeALlDirs(fileNames []string) int {
	k := 0
	n := len(fileNames)
	for i := 0; i < n; i++ {
		err := os.MkdirAll(fileNames[i], 0664)
		if err != nil {
			k++
		}
	}
	return k
}

func MakeLastPathIfNotAbsolute(fileName string) string {
	pos := strings.LastIndex(fileName, "/")
	if fileName == "" || fileName[0] == '/' || pos <= 0 || pos == len(fileName)-1 {
		return fileName
	}
	return fileName[pos+1:]
}

func EnsureFolderForUsualFileToBeSaved(fileName string) {
	pos := strings.LastIndex(fileName, "/")
	pos1 := strings.LastIndex(fileName, "\\")
	if pos1 > pos {
		pos = pos1
	}
	if pos <= 0 {
		return
	}
	os.MkdirAll(fileName[:pos+1], 0755)
}
