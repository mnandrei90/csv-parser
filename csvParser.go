package csvparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

//const ChunkSize = 4 * 1024
const ChunkSize = 100

func readFile() error {
    file, err := os.Open("test_data/test.csv")
    if err != nil {
        return errors.New(fmt.Sprintf("Failed to open file: %v", err))
    }
    defer file.Close()

    // Create a buffered reader
    reader := bufio.NewReader(file)

    // Allocate buffer for chunks
    buffer := make([]byte, ChunkSize)

    var partialRow string

    for {
        n, err := reader.Read(buffer)
        if n > 0 {
            // Read was successful, do something
            fmt.Printf("Read: %v\n", string(buffer[:n]))
        }

        if err != nil {
            if err == io.EOF {
                break
            }

            return errors.New(fmt.Sprintf("Failed to open file: %v", err))
        }
    }
    
    return nil
}
