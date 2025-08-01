package caca

import (
	"fmt"
	"time"

	"github.com/temaxuck/caca/colors"
	"github.com/temaxuck/caca/internal/git"
)

type CanvasMetadata struct {
	Author         *git.GitAuthor
	StartDate      time.Time
	RepositoryPath string
}

type Canvas struct {
	Canvas2D [][]uint8 // Rows: Weekdays, Columns: Number of the week
	Metadata *CanvasMetadata
}

func NewCanvas() *Canvas {
	u, _ := git.GetGlobalUser()
	if u == nil {
		u = &git.GitAuthor{}
	}

	return &Canvas{
		Canvas2D: make([][]uint8, 0, 7),
		Metadata: &CanvasMetadata{
			Author:         u,
			RepositoryPath: ".",
			StartDate:      time.Now(),
		},
	}
}

func (cvs *Canvas) FlatCanvas() []uint8 {
	if len(cvs.Canvas2D) == 0 || len(cvs.Canvas2D[0]) == 0 {
		return nil
	}

	maxCols := 0
	for _, row := range cvs.Canvas2D {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	rows := len(cvs.Canvas2D) // weekdays
	flat := make([]uint8, 0, rows*maxCols)

	for col := 0; col < maxCols; col++ {
		for row := 0; row < rows; row++ {
			if col < len(cvs.Canvas2D[row]) {
				flat[col*rows+row] = cvs.Canvas2D[row][col]
			}
		}
	}
	return flat
}

func (cvs *Canvas) String() string {
	repr := ""
	for weekDay := range len(cvs.Canvas2D) {
		for week := range len(cvs.Canvas2D[weekDay]) {
			i := cvs.Canvas2D[weekDay][week]
			color := colors.FromIntensity(i)
			if color == nil {
				repr += " "
			} else {
				repr += colors.TextBackground(" ", *color)
			}
		}
		repr += "\n"
	}

	return repr
}

func (cvs *Canvas) SetAuthor(name, email string) {
	cvs.Metadata.Author = &git.GitAuthor{name, email}
}

func (cvs *Canvas) SetRepository(repoPath string) {
	cvs.Metadata.RepositoryPath = repoPath
}

func (cvs *Canvas) SetStartDate(date time.Time) {
	cvs.Metadata.StartDate = date
}

func (m *CanvasMetadata) String() string {
	s := fmt.Sprintf("Starting commits from: %s\n", m.StartDate.Format(time.DateOnly))
	s += fmt.Sprintf("Target repository:     %q\n", m.RepositoryPath)
	s += fmt.Sprintf("Author:                \"%s %s\"", m.Author.Name, m.Author.Email)

	return s
}
