package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// retriveUserInput is a function that can retrive user input in form of string. By default,
// it will prompt the user. In test, you can replace this with code that returns the appropriate response.
var retrieveUserInput = func() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	answer = strings.Replace(answer, "\r", "", 1)
	answer = strings.Replace(answer, "\n", "", 1)

	return answer, nil
}

// AskForConfirm parses and verifies user input for confirmation.
func AskForConfirm(message string) error {
	warnConfirm(message)
	answer, err := retrieveUserInput()
	if err != nil {
		return fmt.Errorf("unable to parse users input: %s", err)
	}

	a := strings.ToLower(answer)
	if a != "y" && a != "ye" && a != "yes" {
		return fmt.Errorf("invalid user input")
	}

	return nil
}
