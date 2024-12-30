package csvparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

const ChunkSize = 4 * 1024

type AppSettings struct {
    WithHeader bool
}

var appSettings *AppSettings

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

func (p *CSVParser) ProcessByte(b byte, emitRow func([]string) error) error {
    switch p.state {
        case StateStart:
            if b == ',' {
                p.currentRow = append(p.currentRow, "")
            } else if b == '"' {
                p.state = StateInQuotedField
            } else if b == '\n' {
                if err := emitRow(p.currentRow); err != nil {
                    return err
                }
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
                
                if err := emitRow(p.currentRow); err != nil {
                    return err
                }

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

    return nil
}

func (p *CSVParser) ProcessRemaining(emitRow func([]string) error) error {
    if len(p.currentField) > 0 || len(p.currentRow) > 0 {
        p.currentRow = append(p.currentRow, string(p.currentField))
        if err := emitRow(p.currentRow); err != nil {
            return err
        }
    }

    return nil
}

func extractHeaderRow[T any](structPtr *T, headerRow []string) (map[string]int, error) {
    structType := reflect.TypeOf(structPtr).Elem().Elem()
    headerRowPositions := make(map[string]int)

    for i := 0; i < structType.NumField(); i++ {
        field := structType.Field(i)
        fieldTag := string(field.Tag)
        if fieldTag == "" || !strings.HasPrefix(fieldTag, "csv:") {
            continue
        }

        if field.Type.Kind() != reflect.Slice {
            return nil, errors.New("Only slice types allowed!")
        }

        columnName := strings.Split(fieldTag, ":")[1]
        columnName = strings.Trim(columnName, "\"")
        position := slices.Index(headerRow, columnName)
        headerRowPositions[field.Name] = position
    }

    return headerRowPositions, nil
}

func populateStruct[T any](row []string, structPtr *T, headerRowPositions map[string]int) error {
    structValue := reflect.ValueOf(structPtr).Elem()
    if structValue.Kind() != reflect.Ptr || structValue.Elem().Kind() != reflect.Struct {
        return errors.New("structPtr must be a pointer to a struct")
	}

    if appSettings.WithHeader {
        for _, pos := range headerRowPositions {
            if pos < 0 || pos >= len(row) {
                continue
            }
            if err := populateStructField(&structValue, row, pos); err != nil {
                return err
            }
        }
    } else {
        for i := range row {
            if err := populateStructField(&structValue, row, i); err != nil {
                return err
            }
        }
    }

    return nil
}

func populateStructField(structValue *reflect.Value, row []string, pos int) error {
    if pos >= structValue.Elem().NumField() {
        return nil
    }

    field := structValue.Elem().FieldByIndex([]int{pos})
    if !field.IsValid() {
        return errors.New(fmt.Sprintf("field at position [%v] not valid\n", pos))
    }

    if field.Kind() != reflect.Slice {
        return errors.New(fmt.Sprintf("field at position [%v] not a slice\n", pos))
    }

    if !field.CanSet() {
        return errors.New(fmt.Sprintf("field at position [%v] can't be set\n", pos))
    }

    valueToAppend := reflect.ValueOf(row[pos])

    switch field.Type().Elem().Kind() {
    case reflect.String:
        valueToAppend = reflect.ValueOf(row[pos])
        newSlice := reflect.Append(field, valueToAppend)
        field.Set(newSlice)
    case reflect.Int:
        valueAsInt, err := strconv.Atoi(row[pos])
        if err != nil {
            return errors.New(fmt.Sprintf("Failed to convert value to int: %v\n", err))
        }

        valueToAppend = reflect.ValueOf(valueAsInt)
        newSlice := reflect.Append(field, valueToAppend)
        field.Set(newSlice)
    case reflect.Float32:
        valueAsFloat, err := strconv.ParseFloat(row[pos], 32)
        if err != nil {
            return errors.New(fmt.Sprintf("Failed to convert value to float: %v\n", err))
        }

        valueToAppend = reflect.ValueOf(valueAsFloat)
        newSlice := reflect.Append(field, valueToAppend)
        field.Set(newSlice)
    case reflect.Float64:
        valueAsFloat, err := strconv.ParseFloat(row[pos], 64)
        if err != nil {
            return errors.New(fmt.Sprintf("Failed to convert value to float: %v\n", err))
        }

        valueToAppend = reflect.ValueOf(valueAsFloat)
        newSlice := reflect.Append(field, valueToAppend)
        field.Set(newSlice)
    default:
        return errors.New(fmt.Sprintf("Value type not supported: %v\n", field.Type().Elem().Kind()))
    }

    return nil
}

func concuctStruct[T any](structPtr *T, filePath string, settings *AppSettings) error {
    if settings == nil {
        return errors.New("AppSettings not set!")
    }

    headerRowPositions := make(map[string]int)
    rowsHandled := 0

    handleRow := func(row []string) error {
        var err error

        rowsHandled++
        if settings.WithHeader && rowsHandled == 1 {
            headerRowPositions, err = extractHeaderRow(&structPtr, row)
            if err != nil {
                return err
            }

            return nil
        }

        if err := populateStruct(row, &structPtr, headerRowPositions); err != nil {
            return err
        }

        return nil
    }

    if err := readFile(filePath, settings, handleRow); err != nil {
        return err
    }

    return nil
} 

func readFile(filePath string, settings *AppSettings, emitRow func([]string) error) error {
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
                err := parser.ProcessByte(buffer[i], emitRow)
                if err != nil {
                    return err
                }
            }
        }

        if err != nil {
            if err == io.EOF {
                break
            }

            return errors.New(fmt.Sprintf("Failed to open file: %v", err))
        }
    }

    err = parser.ProcessRemaining(emitRow)
    if err != nil {
        return err
    }
    
    return nil
}
