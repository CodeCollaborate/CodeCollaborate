package utils

import (
	"testing"
	"errors"
)

//
//func TestRead(t *testing.T) {
//	rdr, err := ioutil.ReadFile("test.txt")
//	if err != nil {
//		t.Fatal("Failed to read file")
//	}
//	s := string(rdr)
//	fmt.Println(s)
//	fr := FakeReader{}
//	utils.Read(fr)
//}

func TestFailOnError(t *testing.T) {
	err := errors.New("I'm an error")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Where's the panic?")
		}
	}()
	FailOnError(err, "Fail me")
}

func TestLogOnError(t *testing.T) {
	err := errors.New("I'm also an error")
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Why did you panic?")
		}
	}()
	LogOnError(err, "Fail me also")
}