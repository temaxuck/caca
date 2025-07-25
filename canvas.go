package caca

import (
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
	Count    int

	Metadata *CanvasMetadata
}

func (cvs *Canvas) FlatCanvas() []uint8 {
	flat := make([]uint8, cvs.Count, cvs.Count)
	for i := range len(cvs.Canvas2D) {
		for j := range len(cvs.Canvas2D[i]) {
			flat[j*len(cvs.Canvas2D)+i] = cvs.Canvas2D[i][j]
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

func (cvs *Canvas) SetAuthor(author *git.GitAuthor) {
	cvs.Metadata.Author = author
}

func (cvs *Canvas) SetRepository(repoPath string) {
	cvs.Metadata.RepositoryPath = repoPath
}

func (cvs *Canvas) SetStartDate(date time.Time) {
	cvs.Metadata.StartDate = date
}
