package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileSourceParse(t *testing.T) {
	path := "../../config.yaml"
	file, err := os.Open(path)
	if err != nil {
		_ = fmt.Errorf("%v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		_ = fmt.Errorf("%v", err)
	}
	info, err := file.Stat()
	if err != nil {
		_ = fmt.Errorf("%v", err)
	}
	fmt.Println(info.Name())
	fmt.Println(data, " ===> \n", string(data))
}
