package main

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GetRandom(max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	result := r1.Intn(max)
	return result
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func DeleteFiles(dir string, filePattern string) (err error) {
	filePathPattern := path.Join(dir, filePattern)
	files, err := filepath.Glob(filePathPattern)
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func ConvertBytesToString(input []byte) string {
	return string(input)
	var results []string
	for i := 0; i < len(input); i++ {
		results = append(results, strconv.Itoa(int(input[i])))
	}
	return strings.Join(results, "")
}

func ConvertIntArrayToString(input []int) string {
	var results []string
	for i := 0; i < len(input); i++ {
		results = append(results, strconv.Itoa(input[i]))
	}
	result := strings.Join(results, ",")
	return result
}

func ConvertStringTointArray(input string) []int {
	stringList := strings.Split(input, ",")
	var results []int
	for i := 0; i < len(stringList); i++ {
		lineNumber, _ := strconv.Atoi(stringList[i])
		results = append(results, lineNumber)
	}
	return results
}
