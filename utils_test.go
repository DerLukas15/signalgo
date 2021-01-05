package signalgo

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCommandValid(t *testing.T) {
	_, err := commandValid("pathnotvalid")
	if err == nil {
		t.Errorf("Not existing command passed")
	}
	valid, err := commandValid("/usr/bin/ls")
	if err != nil {
		t.Errorf("ls command threw error: %s", err)
	}
	if !valid {
		t.Errorf("ls command is not valid")
	}

	tempPath, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tempPath)

	fileName := tempPath + "/sampleNonExecutable"
	defer os.Remove(fileName)
	ioutil.WriteFile(fileName, []byte("Temp"), 0644)
	valid, err = commandValid(fileName)
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if valid {
		t.Errorf("sampleNonExecutable is valid")
	}

	fileName = tempPath + "/sampleUserExecutable"
	defer os.Remove(fileName)
	ioutil.WriteFile(fileName, []byte("Temp"), 0744)
	valid, err = commandValid(fileName)
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if !valid {
		t.Errorf("sampleUserExecutable is not valid")
	}

	fileName = tempPath + "/sampleGroupExecutable"
	defer os.Remove(fileName)
	ioutil.WriteFile(fileName, []byte("Temp"), 0654)
	valid, err = commandValid(fileName)
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if !valid {
		t.Errorf("sampleGroupExecutable is not valid")
	}

	fileName = tempPath + "/sampleAllExecutable"
	defer os.Remove(fileName)
	ioutil.WriteFile(fileName, []byte("Temp"), 0645)
	valid, err = commandValid(fileName)
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if !valid {
		t.Errorf("sampleAllExecutable is not valid")
	}

}
