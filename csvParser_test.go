package csvparser

import (
	"log"
	"testing"
)

func TestHelloWorld(t *testing.T) {
    err := readFile()
    if err != nil {
        log.Fatalln(err)
    }
}
