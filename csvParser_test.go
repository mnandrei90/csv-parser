package csvparser

import (
	"log"
	"testing"
)

func TestSmall(t *testing.T) {
    err := readFile("test_data/test.csv")
    if err != nil {
        log.Fatalln(err)
    }
}

func TestLarge(t *testing.T) {
    err := readFile("test_data/test_large.csv")
    if err != nil {
        log.Fatalln(err)
    }
}

func TestExtraLarge(t *testing.T) {
    err := readFile("test_data/test_xl.csv")
    if err != nil {
        log.Fatalln(err)
    }
}
