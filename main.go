package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	if err := Main(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main(args []string) error {
	if len(args) < 3 {
		return errors.New("Usage: dirset list|each set-name [commands]")
	}
	subcmd := args[1]
	setname := args[2]
	dirs, err := Dirs(setname)
	if err != nil {
		return err
	}

	switch subcmd {
	case "list":
		return List(dirs)
	case "each":
		return Each(dirs, args[3:])
	default:
		return fmt.Errorf("%s is unknown subcommand", subcmd)
	}
}

func List(dirs []string) error {
	for _, d := range dirs {
		fmt.Println(d)
	}
	return nil
}

func Each(dirs []string, commands []string) error {
	for _, d := range dirs {
		fmt.Printf("*** In %s\n", d)
		cmd := exec.Command(commands[0], commands[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = d
		cmd.Run()
	}
	return nil
}

func Dirs(name string) ([]string, error) {
	confPath, err := homedir.Expand("~/.config/dirset/" + name)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	s := strings.TrimSpace(string(b))

	return strings.Split(s, "\n"), nil
}
