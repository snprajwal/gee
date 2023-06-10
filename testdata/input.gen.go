/*
 Assuming some Generated file with // go:generate
 This File should not be edited
*/

package test

import "fmt"

func foo() error {
	return nil
}

func bar() (string, error) {
	return "", nil
}

func qux() error {
	var _ error
	//gee:hello there
	_ = foo()
	//gee:
	_ = foo()
	//gee:general kenobi
	bazString, _ := bar()
	fmt.Println(bazString)
	return nil
}
