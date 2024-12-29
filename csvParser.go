package csvparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"slices"
	"strings"
)

type AppSettings struct {
    WithHeader bool
    rowsHandled int
    headerRow []string
}

var appSettings *AppSettings

const ChunkSize = 4 * 1024
// const ChunkSize = 100

const (
    StateStart = iota
    StateInField
    StateInQuotedField
)

type CSVParser struct {
    state int
    currentField []byte
    currentRow []string
    inQuotes bool
}

func NewCSVParser() *CSVParser {
    return &CSVParser {
        state: StateStart,
    }
}

func (p *CSVParser) ProcessByte(b byte, emitRow func([]string)) {
    switch p.state {
        case StateStart:
            if b == ',' {
                p.currentRow = append(p.currentRow, "")
            } else if b == '"' {
                p.state = StateInQuotedField
            } else if b == '\n' {
                emitRow(p.currentRow)
                p.currentRow = nil
            } else {
                p.state = StateInField
                p.currentField = append(p.currentField, b)
            }

        case StateInField:
            if b == ',' {
                p.currentRow = append(p.currentRow, string(p.currentField))
                p.currentField = nil
                p.state = StateStart
            } else if b == '\n' {
                p.currentRow = append(p.currentRow, string(p.currentField))
                p.currentField = nil
                emitRow(p.currentRow)
                p.currentRow = nil
                p.state = StateStart
            } else {
                p.currentField = append(p.currentField, b)
            }

        case StateInQuotedField:
            if b == '"' {
                p.inQuotes = !p.inQuotes
                if !p.inQuotes {
                    p.state = StateInField
                } else {
                    p.currentField = append(p.currentField, b)
                }
            } else {
                p.currentField = append(p.currentField, b)
            }
    }
}

func (p *CSVParser) ProcessRemaining(emitRow func([]string)) {
    if len(p.currentField) > 0 || len(p.currentRow) > 0 {
        p.currentRow = append(p.currentRow, string(p.currentField))
        emitRow(p.currentRow)
    }
}

func handleRow(row []string) {
    appSettings.rowsHandled++
    if appSettings.WithHeader && appSettings.rowsHandled == 1 {
        headerRow := make([]string, len(row))
        headerLen := copy(headerRow, row)
        if headerLen != len(row) {
            log.Fatalln("Header row not correctly created!")
        }

        fmt.Printf("Header row: %v\n", headerRow)
    }
}

func concuctStruct[T any](s T, filePath string) (T, error) {
    structType := reflect.TypeOf(s)

    appSettings := AppSettings { WithHeader: true }
    appSettings.headerRow = []string{"test", "id", "id2", "something else"}

    // csvTags := make([]string, 0)
    for i := 0; i < structType.NumField(); i++ {
        field := structType.Field(i)
        fieldTag := string(field.Tag)
        if fieldTag == "" || !strings.HasPrefix(fieldTag, "csv:") {
            continue
        }

        if field.Type.Kind() != reflect.Slice {
            log.Fatalf("Only slice types allowed! Current type of column %v is: %v", field.Name, field.Type)
        }

        columnName := strings.Split(fieldTag, ":")[1]
        columnName = strings.Trim(columnName, "\"")
        fmt.Printf("Field: %v\n", columnName)
        fmt.Printf("Field type: %v\n", field.Type)

        // get corresponding header row position
        position := slices.Index(appSettings.headerRow, columnName)
        fmt.Printf("Field position in header row: %v\n", position)
    }

    return *new(T), nil
} 

func readFile(filePath string, settings *AppSettings) error {
    appSettings = settings

    file, err := os.Open(filePath)
    if err != nil {
        return errors.New(fmt.Sprintf("Failed to open file: %v", err))
    }
    defer file.Close()

    // Create a buffered reader
    reader := bufio.NewReader(file)

    // Allocate buffer for chunks
    buffer := make([]byte, ChunkSize)

    parser := NewCSVParser()

    for {
        n, err := reader.Read(buffer)
        if n > 0 {
            // Read was successful, do something
            for i := 0; i < n; i++ {
                parser.ProcessByte(buffer[i], handleRow)
            }
        }

        if err != nil {
            if err == io.EOF {
                break
            }

            return errors.New(fmt.Sprintf("Failed to open file: %v", err))
        }
    }

    parser.ProcessRemaining(handleRow)
    
    // fmt.Printf("Rows handled: %v", rowsHandled)

    return nil
}
