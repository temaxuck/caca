package caca

import (
	"fmt"
	"time"

	"github.com/temaxuck/caca/internal/git"
)

func (cvs *Canvas) Preview() error {
	logSettings(*cvs)
	fmt.Printf("INFO: Canvas preview:\n%s", cvs.String())

	return nil
}

func (cvs *Canvas) Draw(verbose bool) error {
	gs, err := git.NewGitService(cvs.Metadata.RepositoryPath, cvs.Metadata.Author)
	if err != nil {
		return err
	}

	logSettings(*cvs)

	days := cvs.FlatCanvas()
	currentDate := cvs.Metadata.StartDate
	for _, day := range days {
		for i := range day {
			hash, err := gs.CommitEmpty(generateCommitMessage(i, currentDate), currentDate)
			if err != nil {
				return err
			}
			if verbose {
				fmt.Printf("INFO: [%d]: Created a commit[%.7s] on %s\n", i, hash, currentDate.Format(time.DateOnly))
			}
		}
		currentDate = currentDate.Add(time.Hour * 24)
	}

	return nil
}

func generateCommitMessage(commitNumber uint8, date time.Time) string {
	return fmt.Sprintf("caca: %s[%d] - drawing over contribution calendar", date.Format(time.DateOnly), commitNumber)
}

func logSettings(cvs Canvas) {
	fmt.Printf("INFO: Canvas settings:\n")
	fmt.Printf("      Starting commits from: %s\n", cvs.Metadata.StartDate.Format(time.DateOnly))
	fmt.Printf("      Target repository:     %q\n", cvs.Metadata.RepositoryPath)
	fmt.Printf("      Author:                \"%s %s\"\n", cvs.Metadata.Author.Name, cvs.Metadata.Author.Email)
}
