package login

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func readPasswordFromStdin(prompt string) (string, error) {
	var out string
	var err error
	fmt.Fprint(os.Stdout, prompt)
	stdinFD := int(os.Stdin.Fd())
	if term.IsTerminal(stdinFD) {
		var inBytes []byte
		inBytes, err = term.ReadPassword(int(os.Stdin.Fd()))
		out = string(inBytes)
	} else {
		out, err = bufio.NewReader(os.Stdin).ReadString('\n')
	}
	if err != nil {
		return "", fmt.Errorf("could not read from stdin: %w", err)
	}
	fmt.Println()
	return strings.TrimRight(out, "\r\n"), nil
}
