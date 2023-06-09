package test

import "fmt"

func foo() error {
	return nil
}

func bar() (string, error) {
	return "", nil
}

func qux() error {
	var err error
	err = foo()
	if err != nil {
		return fmt.Errorf("hello there: %w", err)
	}
	err = foo()
	if err != nil {
		return err
	}
	bazString, err := bar()
	if err != nil {
		return fmt.Errorf("general kenobi: %w", err)
	}
	fmt.Println(bazString)
	return nil
}
