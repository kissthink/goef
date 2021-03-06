package goef

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/siongui/mypkg"
)

func TestImport(t *testing.T) {
	a1, err := ioutil.ReadFile("testdir/hello.txt")
	if err != nil {
		t.Error(err)
		return
	}
	a2, err := ioutil.ReadFile("testdir/backtick.txt")
	if err != nil {
		t.Error(err)
		return
	}
	a3, err := ioutil.ReadFile("testdir/subdir/hello2.txt")
	if err != nil {
		t.Error(err)
		return
	}

	b, err := mypkg.ReadFile("hello.txt")
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(a1, b) {
		t.Error("hello.txt content not correct")
		return
	}

	b, err = mypkg.ReadFile("backtick.txt")
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(a2, b) {
		t.Error("backtick.txt content not correct")
		return
	}

	b, err = mypkg.ReadFile("subdir/hello2.txt")
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(a3, b) {
		t.Error("subdir/hello2.txt content not correct")
		return
	}

	_, err = mypkg.ReadFile("hello3.txt")
	if err != os.ErrNotExist {
		t.Error("hello3.txt should not exit!")
		return
	}
}
