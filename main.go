// This is the name of our package
// Everything with this package name can see everything
// else inside the same package, regardless of the file they are in
package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (
	// "fmt" has methods for formatted I/O operations (like printing to the console)
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"

	// The "net/http" library has methods to implement HTTP clients and servers
	"net/http"
)

func main() {
	// The "HandleFunc" method accepts a path and a function as arguments
	// (Yes, we can pass functions as arguments, and even treat them like variables in Go)
	// However, the handler function has to have the appropriate signature (as described by the "handler" function below)
	http.HandleFunc("/", handler)

	// After defining our server, we finally "listen and serve" on port 8080
	// The second argument is the handler, which we will come to later on, but for now it is left as nil,
	// and the handler defined above (in "HandleFunc") is used
	addr := "127.0.0.0"
	if runtime.GOOS == "darwin" {
		addr = "localhost"
	}
	http.ListenAndServe(fmt.Sprintf("%s:8080", addr), nil)
}

// "handler" is our handler function. It has to follow the function signature of a ResponseWriter and Request type
// as the arguments.
func handler(w http.ResponseWriter, r *http.Request) {
	filePath := os.Getenv("FILE_PATH")
	if len(filePath) == 0 {
		filePath = "/tmp/test.txt"
	}

	entry := readAndAppend(filePath)
	fmt.Fprintf(w, "Hello World! => "+entry)
}

func readAndAppend(fileName string) string {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	fileDescriptor, err := os.Stat(fileName)
	lineBreak := []byte("\n")
	lineBreakSize := int64(len(lineBreak))
	current := fileDescriptor.Size() - lineBreakSize
	entry := 0
	if current > 0 {
		lineBuf := make([]byte, lineBreakSize)

		for {
			n, _ := file.ReadAt(lineBuf, current)

			if bytes.Equal(lineBuf[:n], lineBreak) {
				break
			}

			current -= lineBreakSize
		}

		numberBuf := make([]byte, fileDescriptor.Size()-current-lineBreakSize)
		_, _ = file.ReadAt(numberBuf, current+lineBreakSize)

		entry, err = strconv.Atoi(string(numberBuf))
		if err != nil {
			panic(err)
		}
	}

	entry += 1
	if _, err := file.WriteString(fmt.Sprintf("\n%d", entry)); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%d", entry)
}
