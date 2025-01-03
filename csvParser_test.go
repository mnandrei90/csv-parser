package csvparser

import (
	"fmt"
	"log"
	"testing"
)

type MyStruct struct {
    Id []int `csv:"id"`
    Name []string `csv:"firstname"`
    LastName []string `csv:"lastname"`
    Email []string `csv:"email"`
    SecondaryEmail []string `csv:"email2"`
    Profession []string `csv:"profession"`
}

type MyStruct2 struct {
    Id []int
    Name []string
    LastName []string
    Email []string
    SecondaryEmail []string
    Profession []string
}

type MyStruct3 struct {
    First []string
    Second []int
    Third []float32
    Fourth []string
}

func TestSpecial(t *testing.T) {
    testStruct := MyStruct3{}
    err := ConcuctStruct(&testStruct, "test_data/test.csv", &AppSettings{WithHeader: false})
    if err != nil {
        log.Fatalln(err)
    }
    
    fmt.Println(testStruct)
}

func TestConcuct(t *testing.T) {
    testStruct := MyStruct{}
    err := ConcuctStruct(&testStruct, "test_data/test_w_header.csv", &AppSettings{WithHeader: true})
    if err != nil {
        log.Fatalln(err)
    }
    
    fmt.Println(testStruct)
}

func TestConcuct2(t *testing.T) {
    testStruct := MyStruct2{}
    err := ConcuctStruct(&testStruct, "test_data/test_wo_header.csv", &AppSettings{WithHeader: false})
    if err != nil {
        log.Fatalln(err)
    }
    
    fmt.Println(testStruct)
}

func TestConcuct3(t *testing.T) {
    testStruct := MyStruct{}
    err := ConcuctStruct(&testStruct, "test_data/test_xl.csv", &AppSettings{WithHeader: true})
    if err != nil {
        log.Fatalln(err)
    }
    
    fmt.Println(testStruct.Name[99999])
}
