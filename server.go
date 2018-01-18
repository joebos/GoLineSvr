package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// LineServer encapsulates the server object.
type LineServer struct {
	settings       *Settings
	lineFile       *LineFile
	controlChannel chan string
}

// Start is the public method for starting the line server
func (obj *LineServer) Start(texFilePath string) {

	// This channel is for the goroutine handing client request to send message back to shutdown the server
	// whenever a user issues "shutdown" command
	obj.controlChannel = make(chan string)

	//lineFile is teh data model object which handles text file pre-processing (building indexes) and retrieve a line by line number
	obj.lineFile = NewLineFile(texFilePath, obj.settings.numLinesPerIndexPage)
	obj.lineFile.BuildIndex()

	hostPort := obj.settings.host + ":" + strconv.Itoa(obj.settings.port)
	tcpAddress, err := net.ResolveTCPAddr("tcp4", hostPort)
	obj.checkError(err)

	// Start to listening a TCP port
	listener, err := net.ListenTCP("tcp", tcpAddress)
	obj.checkError(err)

	log.Printf("The Line Server is about to start....")

	// Start checcking shutdown command on a seperate goroutine
	go obj.handleControlMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// New thread to handle client connection
		go obj.handleClientConnection(conn)
	}
}

// Exit the program whenever "shutdown" is received from the control channel
func (obj *LineServer) handleControlMessage() {
	for {
		msg := <-obj.controlChannel
		if msg == "SHUTDOWN" {
			os.Exit(0)
		}
	}
}

// handleClientConnection - handles all client commands.
// When a connect issues quit, shutdown or disconnect, this method will return
func (obj *LineServer) handleClientConnection(conn net.Conn) {

	log.Printf("Received a client connection from %s...", conn.RemoteAddr().String())

	// close connection on exit
	defer conn.Close()

	var buf [1024]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		if n == 0 {
			return
		}

		log.Printf("Received %s from %s...", string(buf[:]), conn.RemoteAddr().String())
		isValid, cmd, params := obj.validateCcommand(buf, conn)
		if isValid {
			switch cmd {
			case "SHUTDOWN":
				obj.controlChannel <- "SHUTDOWN"
				break
			case "QUIT":
				return
			case "GET":
				lineNo := params[0]
				status, line := obj.lineFile.GetLine(lineNo)
				responseMsg := "ERR"
				switch status {
				case 500:
					responseMsg = "Server error: " + line + "\n"
				case 200:
					responseMsg = "OK\n" + line + "\n"
				default:
					responseMsg = "ERR\n"
				}
				response := []byte(responseMsg)
				_, err2 := conn.Write(response)
				if err2 != nil {
					return
				}
			}
		} else {
			response := []byte("INVALID COMMAND!\n")
			_, err2 := conn.Write(response)
			if err2 != nil {
				return
			}
		}
	}
}

// validateCcommand - this is a helper method for processing client command message
func (obj *LineServer) validateCcommand(input [1024]byte, conn net.Conn) (isValid bool, cmd string, params []int) {
	r := bytes.NewReader(input[:])
	r2 := bufio.NewReader(r)
	chars, _, _ := r2.ReadLine()
	cmdString := ConvertBytesToString(chars)
	log.Printf("Received command %s from %s...", cmdString, conn.RemoteAddr().String())
	if cmdString == "" {
		return false, "", []int{}
	}
	cmdString = strings.ToUpper(strings.TrimSpace(cmdString))
	args := strings.Split(cmdString, " ")
	if args[0] == "QUIT" {
		return true, "QUIT", []int{}
	}
	if args[0] == "SHUTDOWN" {
		return true, "SHUTDOWN", []int{}
	}
	if args[0] == "GET" && len(args) >= 2 {
		line := args[1]
		lineNo, err := strconv.Atoi(line)
		if err == nil {
			return true, "GET", []int{int(lineNo)}
		}
	}
	return false, "", []int{}
}

func (obj *LineServer) checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
