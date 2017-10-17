// ioops project ioops.go
package ioops

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

//creates a directory with the given name.
//if directory already exists it ignores.
//returns error if any error occures while creation
func CreateDirectory(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return os.Mkdir(dir, os.ModePerm)
	}
	return nil
}

//creates a file with the given name.
//returns error if any error occures while creation
func CreateFile(fn string) (err error) {
	file, err := os.Create(fn)
	if err != nil {
		return nil
	}
	_, err = file.WriteString("0")
	if err != nil {
		return err
	}
	return nil
}

//moves a file from source directory to the destination directory.
//file move is easy option than actually copying and deleting from the destination later..
func MoveFile(src, dest string) error {
	err := os.Rename(src, dest)
	if err != nil {
		return err
	}
	return nil
}

//copies a file to a dest file name from the io.reader stream.
func CopyToFile(r io.Reader, dest string) (err error) {
	//defer file.Close()
	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)

	if err != nil {
		return err
	}
	return nil
}

//read bytes and return as reader
func Read(buf []byte) io.Reader {
	return bytes.NewReader(buf)
}

//extract the file extension from the file name.
func GetFileExt(fn string) string {
	return fn[strings.LastIndex(fn, "."):]
}

func FileCountIncrement(filename string) (count int, err error) {
	FileMutex := &sync.Mutex{}
	FileMutex.Lock()
	defer FileMutex.Unlock()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return count, err
	}
	count, err = strconv.Atoi(string(data))
	if err != nil {
		return count, err
	}
	count = count + 1
	ioutil.WriteFile(filename, []byte(strconv.Itoa(count)), 0644)

	return count, err
}

//Get the count of countfile for each movie id
func GetFileCount(filename string) (count int, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return count, err
	}
	if string(data) == "" {
		return count, errors.New("wrong count")
	}
	count, err = strconv.Atoi(string(data))
	return count, err
}

//splits the path as root folder, movieid and inside folder and filename
func SplitPath(path string) (root, movieid, folder, file string) {
	s := strings.Split(path, "/")
	if len(s) == 4 {
		root = s[0]
		movieid = s[1]
		folder = s[2]
		file = s[3]
		return root, movieid, folder, file
	}
	return "", "", "", ""
}
