package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/temaxuck/caca"
)

type Date struct {
	t     time.Time
	isSet bool
}
type User struct {
	name, email string
	isSet       bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS] <canvas file>\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	help := flag.Bool("h", false, "Help message")
	verbose := flag.Bool("v", false, "Enable verbose mode")
	preview := flag.Bool("p", false, "Enable preview mode\nWith this option enabled no commits are made")
	repository := flag.String("repository", "", "Target repository")
	user := User{}
	flag.Var(&user, "user", "User on behalf of whom to create commits. Format: '<name> <email>'\nIf not specified global config user setting is used")
	startDate := Date{}
	flag.Var(&startDate, "start-date", "Set a start date for the canvas\nYou, probably, want it to be a Sunday")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Println("ERROR: Specify canvas file")
		flag.Usage()
		os.Exit(1)
	}

	canvasFile := os.Args[len(os.Args)-flag.NArg()]
	canvas, err := caca.ParseCanvasFile(canvasFile)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	if *repository != "" {
		canvas.SetRepository(*repository)
	}
	if startDate.isSet {
		canvas.SetStartDate(time.Time(startDate.t))
	}
	if user.isSet {
		canvas.SetAuthor(user.name, user.email)
	}

	if *preview {
		if err := canvas.Preview(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err := canvas.Draw(*verbose); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

}

func (d *Date) String() string {
	return time.Time(d.t).Format(time.DateOnly)
}

func (d *Date) Set(value string) error {
	t, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return err
	}

	t = t.Add(time.Hour * 12)
	d.t = t
	d.isSet = true

	return nil
}

func (u *User) String() string {
	if u.isSet {
		return u.name + " " + u.email
	}
	return "null"
}

func (u *User) Set(value string) error {
	tokens := strings.Split(value, " ")
	if len(tokens) == 1 {
		return errors.New("failed to parse user")
	}

	u.name, u.email = strings.Join(tokens[:len(tokens)-1], " "), tokens[len(tokens)-1]
	u.isSet = true
	return nil
}
