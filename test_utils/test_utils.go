package test_utils

import (
	"fmt"
	"log"
	//"reflect" // No biggie in test
	"testing"
)

func ignore() { log.Println("") }

func AreEqual(t *testing.T, expected interface{}, actual interface{}, message string) {
	if expected != actual {
		log.Printf("Expected [%s] but was [%s]: %s", expected, actual, message)
		t.Fail()
	}
}

func AreNotEqual(t *testing.T, expected interface{}, actual interface{}, message string) {
	if expected == actual {
		log.Printf("Expected [%s] but was [%s]: %s", expected, actual, message)
		t.Fail()
	}
}

func AreAllEqual(t *testing.T, expecteds interface{}, actuals interface{}, message string) {
	expected := fmt.Sprintf("%s - % x", expecteds, expecteds)
	actual := fmt.Sprintf("%s - % x", actuals, actuals)
	if expected != actual {
		log.Printf("Expected [%s] but was [%s]: %s", expected, actual, message)
		t.Fail()
	}
}
