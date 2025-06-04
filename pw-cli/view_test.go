package main

import (
	"fmt"
	"testing"
)

func TestRandomGeneratePassword(t *testing.T) {
	pwd, err := randomGeneratePassword(10)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pwd)
	if len(pwd) != 10 {
		t.Errorf("Expected %d, receive %d: pwd:%v", 10, len(pwd), pwd)
	}

}
