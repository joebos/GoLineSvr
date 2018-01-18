package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
)

// LineFile is teh data model type (class)
type LineFile struct {
	filePath             string
	fileDir              string
	filePathHash         string
	indexPage            []int
	numLines             int
	numLinesPerIndexPage int
	currentIndexPage     int
	indexCompleted       bool
}

// NewLineFile "constructor"
func NewLineFile(filePath string, numLinesPerIndexPage int) *LineFile {
	lineFile := new(LineFile)
	lineFile.filePath = filePath
	lineFile.numLinesPerIndexPage = numLinesPerIndexPage
	lineFile.fileDir, _ = path.Split(lineFile.filePath)
	lineFile.filePathHash = GetMD5Hash(filePath)
	lineFile.currentIndexPage = -1
	//lineFile.indexPage = make([]int, numLinesPerIndexPage)
	return lineFile
}

// BuildIndex - pre-process the text file and build indexes for O(1) retrieval
// The index is similar to cluster indexes in relational database, whith line numeber as the key,
// and the location (offset from file start) of a line in the text file as the value
//
// Because the text file could have a large number of lines, so the method splits indexes into multiple index files.
// Let's say we store the indexes for m lines in one index file. When a client asking for the nth line, the method
// first find which index page by doing n/m, then load that index page. It thens compute n % m to find the file offset location
// of the line being requested  in the indexes just loaded.
func (obj *LineFile) BuildIndex() {

	log.Printf("Building indexes: %s is started...", obj.filePath)

	file, err := os.Open(obj.filePath)
	obj.checkError(err)

	obj.deleteIndexFiles()

	var lastOffset int
	r := bufio.NewReader(file)
	data, err := obj.readLine(r)
	obj.checkError(err)
	var i int
	for ; err != io.EOF && data != nil && len(data) > 0; i++ {
		obj.checkError(err)
		obj.indexPage = append(obj.indexPage, lastOffset)
		lastOffset += int(len(data)) + 1
		if (i+1)%obj.numLinesPerIndexPage == 0 {
			obj.writeToIndexPage(i / obj.numLinesPerIndexPage)
			obj.indexPage = obj.indexPage[:0]
		}
		data, err = obj.readLine(r)
	}
	if len(obj.indexPage) > 0 {
		obj.writeToIndexPage((i - 1) / obj.numLinesPerIndexPage)
		obj.indexPage = obj.indexPage[:0]
	}
	obj.numLines = i
	obj.indexCompleted = true
	log.Printf("Building indexes is completed: %s", obj.filePath)

	file.Close()
}

// All index files are stored in a subfolder called index under where the text file is located.
// The name of a index file consists of the hash of text file path, and a integer number
// which equals to n/m where n is the line number and m is the number of lines handled by each index page
func (obj *LineFile) getIndexFilePath(indexPageNumber int) string {
	filePrefix := obj.filePathHash
	indexFolder := path.Join(obj.fileDir, "index")
	os.MkdirAll(indexFolder, 0755)
	indexFilePath := path.Join(indexFolder, filePrefix+"_"+strconv.Itoa(indexPageNumber)+".idx")
	return indexFilePath
}

// Delete all index files
func (obj *LineFile) deleteIndexFiles() {
	filePrefix := obj.filePathHash
	indexFolder := path.Join(obj.fileDir, "index")
	os.MkdirAll(indexFolder, 0755)
	DeleteFiles(indexFolder, filePrefix+"_*.idx")
}

// Write to a index file
func (obj *LineFile) writeToIndexPage(indexPageNumber int) {
	indexFilePath := obj.getIndexFilePath(indexPageNumber)
	file, err := os.Create(indexFilePath)
	obj.checkError(err)
	w := bufio.NewWriter(file)
	_, err = w.WriteString(ConvertIntArrayToString(obj.indexPage))
	obj.checkError(err)
	err = w.Flush()
	obj.checkError(err)
	file.Close()
}

// This method reads a line from a reader, regardless of how long the line is
func (obj *LineFile) readLine(r *bufio.Reader) ([]byte, error) {
	var results []byte
	hasMore := true
	bytes := []byte{}
	var err error
	for hasMore {
		bytes, hasMore, err = r.ReadLine()
		if err != nil {
			break
		}
		if len(bytes) > 0 {
			results = append(results, bytes...)
		}
	}
	return results, err
}

// load a particular index file
func (obj *LineFile) loadIndexPage(indexPageNumber int) {
	indexFilePath := obj.getIndexFilePath(indexPageNumber)
	file, err := os.Open(indexFilePath)
	obj.checkError(err)
	r := bufio.NewReader(file)
	bytes, err := obj.readLine(r)
	obj.checkError(err)
	jsonString := ConvertBytesToString(bytes)
	obj.indexPage = ConvertStringTointArray(jsonString)
	obj.currentIndexPage = indexPageNumber
	file.Close()
}

// GetLine is the method for retriving a line. It first finds which index page file,
// then get the file offset of the line, and then open the text file, seek to that offset and read the entire line
func (obj *LineFile) GetLine(lineNo int) (status int, line string) {
	if !obj.indexCompleted {
		obj.BuildIndex()
	}
	lineNumber0 := lineNo - 1
	if lineNumber0 >= 0 && lineNumber0 < obj.numLines {
		indexPageNumber := lineNumber0 / obj.numLinesPerIndexPage
		if obj.currentIndexPage != indexPageNumber {
			obj.loadIndexPage(indexPageNumber)
		}
		lineAddress := obj.indexPage[lineNumber0%obj.numLinesPerIndexPage]
		file, err := os.Open(obj.filePath)
		if err != nil {
			return 500, "Woops! Something went wrong. Please try again"
		}
		defer file.Close()
		r := bufio.NewReader(file)
		file.Seek(int64(lineAddress), 0)
		bytes, err := obj.readLine(r)
		if err == nil {
			line = string(bytes)
			return 200, line
		}
		return 500, "Woops! Something went wrong. Please try again"
	}
	return 404, ""
}

func (obj *LineFile) checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
