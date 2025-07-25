# caca

*caca (Contribution Activity Canvas)* - use your GitHub profile's [contributions
calendar](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-github-profile/managing-contribution-settings-on-your-profile/viewing-contributions-on-your-profile#contributions-calendar)
as a canvas, you can draw over.

This package provides [library](#api) and [CLI tool](#cli).

> [!WARNING]
> This project is just a proof of concept. 
> The API and CLI are still WIP.

## Proof of concept

It is possible to use GitHub Contribution Activity Calendar as a canvas:
[poc](https://github.com/temaxuck/caca/tree/main/demo/poc.png)

## Canvas file format

Canvas file is a text file containing a grid, where each element $I_{ij} \in \{ 0, 1, 2, 3 \}$ with $i \in [1, 7]$ and $j \in \mathbb{N}$, represents the contribution intensity level for $i$-th day of week in the $j$-th week starting from the **start date**.

### Example

```
# Comments are allowed
#
# Each numeric value is an intensity level:
# 0 - no contributions
# 1..4 - from the lowest to highest contribution activity level for each day
#
# Rows represent day of the week, starting from start date

00000000000000000000000000000000000000000000000000000
00000000000044440000333300000022220000111100000000000
00000000004400000033000033002200000011000011000000000
00000000004400000033000033002200000011000011000000000
00000000004400000033333333002200000011111111000000000
00000000000044440033000033000022220011000011000000000
00000000000000000000000000000000000000000000000000000
```

## API

### Installation

`go get github.com/temaxuck/caca`

### Quick start

```go
    canvas, _ := ReadCanvas(pathToCanvasFile)
    
    // Basically name doesn't matter much, but email must be set accordingly to your GitHub account
    canvas.SetAuthor("Your name", "your-email@nowhere.com")
    // Set repository where the commits are going to be made
    canvas.SetRepository("path-to-repo")
    // Set repository where the commits are going to be made
    canvas.SetStartDate(time.Date(2020, time.December, 27, 12, 0, 0, 0, time.UTC))
    
    canvas.Draw(true)
```

### Canvas type API

```go
type Canvas struct {
	Canvas2D [][]uint8 // Row: Day of the week, Column: Number of the week
	Count    int

    // Use SetAuthor, SetRepository, SetStartDate to set metadata
	Metadata *CanvasMetadata 
}
```
### Methods API

```go
// Read Canvas from file
func ReadCanvas(path string) (*Canvas, error)

// Set metadata
func (cvs *Canvas) SetAuthor(name, email string) 
func (cvs *Canvas) SetRepository(repoPath string)
func (cvs *Canvas) SetStartDate(date time.Time)

// Start drawing over contribution calendar
func (cvs *Canvas) Draw(verbose bool) error
```

## CLI

### Install

`go install github.com/temaxuck/caca/cmd/caca`

Also make sure your `PATH` variable is set to include installed go command line tools:

`PATH=$PATH:~/go/bin`

### Usage

```shell
Usage: caca [OPTIONS] <canvas file>
Options:
  -h	Help message
  -p	Enable preview mode
    	With this option enabled no commits are made
  -repository string
    	Target repository (default ".")
  -start-date value
    	Set a start date for the canvas
    	You, probably, want it to be a Sunday (default 2025-07-25)
  -user value
    	User on behalf of whom to create commits. Format: '<name> <email>'
    	If not specified global config user setting is used
  -v	Enable verbose mode
```

### Examples 

1. Get help message
```shell
$ caca -h
Usage: caca [OPTIONS] <canvas file>
Options:
  -h	Help message
  -p	Enable preview mode
    	With this option enabled no commits are made
  -repository string
    	Target repository (default ".")
  -start-date value
    	Set a start date for the canvas
    	You, probably, want it to be a Sunday (default 2025-07-25)
  -user value
    	User on behalf of whom to create commits. Format: '<name> <email>'
    	If not specified global config user setting is used
  -v	Enable verbose mode
```

2. Set start date:
```shell
$ caca --start-date 2020-12-27 canvas.txt
INFO: Starting commits from: 2020-12-27
```

3. See preview:
```shell
$ caca -p canvas.txt
```
This command produces output like this:
![preview command](https://github.com/temaxuck/caca/tree/main/demo/preview.png)

4. Set user and repository:
```shell
$ caca --user 'Your Name your-email@nowhere.com' --repository path/to/your/repo canvas.txt
```

## Contribution and Support

Feel free to fire off a PR, or open an Issue if you want to report a bug or request a feature.

## License

This project comes with ISC License, see [LICENSE](#TODO-ADD-LICENSE-LINK) for details.
