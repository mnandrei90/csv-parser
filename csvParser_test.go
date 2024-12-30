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

func TestConcuct(t *testing.T) {
    testStruct := MyStruct{}
    // testStruct := MyStruct2{}
    err := concuctStruct(&testStruct, "test_data/test_w_header.csv", &AppSettings{WithHeader: true})
    if err != nil {
        log.Fatalln(err)
    }
    
    fmt.Println(testStruct)
    // concuctStruct(MyStruct2{}, "")
}

func TestConcuct2(t *testing.T) {
    // testStruct := MyStruct{}
    testStruct := MyStruct2{}
    err := concuctStruct(&testStruct, "test_data/test_wo_header.csv", &AppSettings{WithHeader: false})
    if err != nil {
        log.Fatalln(err)
    }
    
    fmt.Println(testStruct)
    // concuctStruct(MyStruct2{}, "")
}

func TestConcuct3(t *testing.T) {
    // testStruct := MyStruct{}
    testStruct := MyStruct{}
    err := concuctStruct(&testStruct, "test_data/test_xl.csv", &AppSettings{WithHeader: true})
    if err != nil {
        log.Fatalln(err)
    }
    
    fmt.Println(testStruct.Name[99999])
    // concuctStruct(MyStruct2{}, "")
}
