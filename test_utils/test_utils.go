package test_utils

import (
	"fmt"
	"log"
	//"reflect" // No biggie in test
	"testing"
)

func ignore() { log.Println("") }

type PanicFunc func()

func ExpectPanic(t *testing.T, message string) PanicFunc {
	return func() {
		if r := recover(); r == nil {
			log.Printf("Expected panic: %s", message)
			t.Fail()
		}
	}
}

func IsNil(t *testing.T, target interface{}, message string) {
	if target != nil {
		printMessageAndFail(t, "nil", target, message)
	}
}

func IsNotNil(t *testing.T, target interface{}, message string) {
	if target == nil {
		printMessageAndFail(t, "nil", target, message)
	}
}

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
	log.Printf("**\tExpected [%s] but was [%s]: %s", expected, actual, message)
	log.Printf("\thex-> [% x] but was [% x]: %s\n\n", expected, actual, message)
	t.Fail()
}
