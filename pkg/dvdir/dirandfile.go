/***********************************************************************
MicroCore
Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdir

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyWholeFolder(src string, dest string) error {
	src = GetFolderNameWithoutLastSlash(src)
	dest = GetFolderNameWithoutLastSlash(dest)
	entries, err := os.ReadDir(src)
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
		info, e := entry.Info()
		if e != nil {
			continue
		}
		isSymlink := (info.Mode() & os.ModeSymlink) != 0
		if !isSymlink {
			if err := os.Chmod(destPath, info.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func CopySingleFile(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer in.Close()

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

func MoveFile(sourcePath, destPath string, strictRemove bool) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove

	err = os.Remove(sourcePath)
	if err != nil && strictRemove {
		return fmt.Errorf("couldn't remove source file: %v", err)
	}
	return nil
}

func RenameOrMoveFile(sourcePath, destPath string, strictRemove bool) error {
	err := os.Rename(sourcePath, destPath)
	if err == nil {
		return nil
	}
	err = MoveFile(sourcePath, destPath, strictRemove)
	return err
}
