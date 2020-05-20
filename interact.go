package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
)

type actionFunc func(answer string) error

type action struct {
	Question string
	Validate actionFunc
	Action   actionFunc
}

func userAction(a *action) error {
	printMsg(a.Question)

	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		inp := scn.Text()
		if err := a.Validate(inp); err != nil {
			printErr(err.Error())
			continue
		}

		if err := a.Action(inp); err != nil {
			return err
		}

		break
	}

	return nil
}

func printErr(msg string) {
	fmt.Fprint(os.Stderr, color.RedString("%s\r\n", msg))
}

func printMsg(msg string) {
	fmt.Fprint(os.Stdin, fmt.Sprintf("%s\r\n", msg))
}

func printSuccess(msg string) {
	fmt.Fprint(os.Stdin, color.GreenString("%s\r\n", msg))
}
