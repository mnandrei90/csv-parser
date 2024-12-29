package csvparser

import (
	"log"
	"testing"
)

func TestSmall(t *testing.T) {
    settings := AppSettings { WithHeader: false }
    err := readFile("test_data/test.csv", &settings)
    if err != nil {
        log.Fatalln(err)
    }
}

func TestLarge(t *testing.T) {
    settings := AppSettings { WithHeader: false }
    err := readFile("test_data/test_large.csv", &settings)
    if err != nil {
        log.Fatalln(err)
    }
}

func TestExtraLarge(t *testing.T) {
    settings := AppSettings { WithHeader: true }
    err := readFile("test_data/test_xl.csv", &settings)
    if err != nil {
        log.Fatalln(err)
    }
}

type MyStruct struct {
    Id []int `csv:"id"`
    Name []string `csv:"first name"`
}

type MyStruct2 struct {
    Id int `csv:"id2"`
    Name string `csv:"first name2"`
}

func TestConcuct(t *testing.T) {
    concuctStruct(MyStruct{}, "")
    concuctStruct(MyStruct2{}, "")
}
