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
		printMessageAndFail(t, expected, actual, message)
	}
}

func AreNotEqual(t *testing.T, expected interface{}, actual interface{}, message string) {
	if expected == actual {
		printMessageAndFail(t, expected, actual, message)
	}
}

func AreAllEqual(t *testing.T, expecteds interface{}, actuals interface{}, message string) {
	expected := fmt.Sprintf("%s - % x", expecteds, expecteds)
	actual := fmt.Sprintf("%s - % x", actuals, actuals)
	if expected != actual {
		printMessageAndFail(t, expected, actual, message)
	}
}

func printMessageAndFail(t *testing.T, expected interface{}, actual interface{}, message string) {
	log.Printf("Expected [%s] but was [%s]: %s", expected, actual, message)
	log.Printf("Expected [% x] but was [% x]: %s", expected, actual, message)
	t.Fail()
}
