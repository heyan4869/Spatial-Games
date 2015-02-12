package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*===============================================================
 * Functions to manipulate a "field" of cells --- the main data
 * that must be managed by this program.
 *==============================================================*/

// The data stored in a single cell of a field
type Cell struct {
	kind  string
	score float64
}

// createField should create a new field of the ysize rows and xsize columns,
// so that field[r][c] gives the Cell at position (r,c).
func createField(rsize, csize int) [][]Cell {
	f := make([][]Cell, rsize)
	for i := range f {
		f[i] = make([]Cell, csize)
	}
	return f
}

// inField returns true if (row,col) is a valid cell in the field
func inField(field [][]Cell, row, col int) bool {
	return row >= 0 && row < len(field) && col >= 0 && col < len(field[0])
}

// readFieldFromFile should open the given file and read the initial
// values for the field. The first line of the file will contain
// two space-separated integers saying how many rows and columns
// the field should have:
//    10 15
// each subsequent line will consist of a string of Cs and Ds, which
// are the initial strategies for the cells:
//    CCCCCCDDDCCCCCC
//
// If there is ever an error reading, this function should cause the
// program to quit immediately.
func readFieldFromFile(filename string) [][]Cell {
    // Open the file 
    in, err := os.Open(filename)
        if err != nil {
            fmt.Println("Error: couldn鈥檛 open the file")
			os.Exit(3) 
		}

	var lines []string = make([]string, 0)

	// Read the file content
    scanner := bufio.NewScanner(in)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
	}

	// Check if the read is right
	if scanner.Err() != nil {
        fmt.Println("Error: there was a problem reading the file")
        os.Exit(3)
	}

	// Get the rows and columns value from the content
	var items []string = strings.Split(lines[0], " ")
	rows, err1 := strconv.Atoi(items[0])
	if err1 != nil {
		fmt.Println("Error: unable to convert.")
		os.Exit(3)
	}
	columns, err2 := strconv.Atoi(items[1])
	if err2 != nil {
		fmt.Println("Error: unable to convert.")
		os.Exit(3)
	}

	// Close the file
	in.Close()

	// Initial the state of each cell
	firstfield := createField(rows, columns)

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			var state string = lines[i+1]
			firstfield[i][j].kind = state[j:j+1]
		}
	}
	return firstfield
}

// drawField should draw a representation of the field on a canvas and save the
// canvas to a PNG file with a name given by the parameter filename.  Each cell
// in the field should be a 5-by-5 square, and cells of the "D" kind should be
// drawn red and cells of the "C" kind should be drawn blue.
func drawField(secondtofield [][]Cell, field [][]Cell, filename string) {
    // Draw the picture based on the state of each cell at the last step and the state
    // before the last state
    var row int = len(field)
    var column int = len(field[0])
    picSpatial := CreateNewCanvas(5*row, 5*column)
    picSpatial.SetLineWidth(1)
    for i := 0; i < row; i++ {
    	for j := 0; j < column; j++ {
    		picSpatial.MoveTo(float64(i*5.0), float64(j*5.0))
    		if field[i][j].kind == "C" && secondtofield[i][j].kind == "C" {
    			picSpatial.SetFillColor(MakeColor(0, 0, 255))
    		} else if field[i][j].kind == "D" && secondtofield[i][j].kind == "C" {
    			picSpatial.SetFillColor(MakeColor(255, 255, 0))
    		} else if field[i][j].kind == "C" && secondtofield[i][j].kind == "D" {
    			picSpatial.SetFillColor(MakeColor(0, 128, 0))
    		} else if field[i][j].kind == "D" && secondtofield[i][j].kind == "D" {
    			picSpatial.SetFillColor(MakeColor(255, 0, 0))
    		}
    		picSpatial.LineTo(float64(i*5.0+5.0), float64(j*5.0))
    		picSpatial.LineTo(float64(i*5.0+5.0), float64(j*5.0+5.0))
    		picSpatial.LineTo(float64(i*5.0), float64(j*5.0+5.0))
    		picSpatial.LineTo(float64(i*5.0), float64(j*5.0))
    		picSpatial.Fill()

    	}
    }
    picSpatial.SaveToPNG(filename)

}

/*===============================================================
 * Functions to simulate the spatial games
 *==============================================================*/

// play a game between a cell of type "me" and a cell of type "them" (both me
// and them should be either "C" or "D"). This returns the reward that "me"
// gets when playing against them.
func gameBetween(me, them string, b float64) float64 {
	if me == "C" && them == "C" {
		return 1
	} else if me == "C" && them == "D" {
		return 0
	} else if me == "D" && them == "C" {
		return b
	} else if me == "D" && them == "D" {
		return 0
	} else {
		fmt.Println("type ==", me, them)
		panic("This shouldn't happen")
	}
}

// updateScores goes through every cell, and plays the Prisoner's dilema game
// with each of it's in-field nieghbors (including itself). It updates the
// score of each cell to be the sum of that cell's winnings from the game.
func updateScores(field [][]Cell, b float64) {
    // Update the score in each cell
    var row int = len(field)
    var column int = len(field[0])
    for i := 0; i < row; i++ {
    	for j := 0; j < column; j++ {
    		field[i][j].score = gameScore(field, i, j, b)
    	}
    }
}

// gameScore calculate the score that a cell could get after it plays game with
// its neighbors
func gameScore(field [][]Cell, i int, j int, b float64) float64 {
	var self string = field[i][j].kind
	var oneScore float64 = 0

	// Check if the cell's neighbor is in the field and calculate the score
	if inField(field, i-1, j-1) == true {
		var other string = field[i-1][j-1].kind
		oneScore = oneScore + gameBetween(self, other, b) 
	}
	if inField(field, i-1, j) == true {
		var other string = field[i-1][j].kind
		oneScore = oneScore + gameBetween(self, other, b)
	}
	if inField(field, i-1, j+1) == true {
		var other string = field[i-1][j+1].kind
		oneScore = oneScore + gameBetween(self, other, b) 
	}
	if inField(field, i, j-1) == true {
		var other string = field[i][j-1].kind
		oneScore = oneScore + gameBetween(self, other, b) 
	}
	if inField(field, i, j) == true {
		var other string = field[i][j].kind
		oneScore = oneScore + gameBetween(self, other, b) 
	}
	if inField(field, i, j+1) == true {
		var other string = field[i][j+1].kind
		oneScore = oneScore + gameBetween(self, other, b) 
	}
	if inField(field, i+1, j-1) == true {
		var other string = field[i+1][j-1].kind
		oneScore = oneScore + gameBetween(self, other, b) 
	}
	if inField(field, i+1, j) == true {
		var other string = field[i+1][j].kind
		oneScore = oneScore + gameBetween(self, other, b) 
	}
	if inField(field, i+1, j+1) == true {
		var other string = field[i+1][j+1].kind
		oneScore = oneScore + gameBetween(self, other, b) 
	}
	return oneScore

}

// updateStrategies create a new field by going through every cell (r,c), and
// looking at each of the cells in its neighborhood (including itself) and the
// setting the kind of cell (r,c) in the new field to be the kind of the
// neighbor with the largest score
func updateStrategies(field [][]Cell) [][]Cell {
    // Update the strategy of each cell for the next step 
	var r int = len(field)
	var c int = len(field[0])
	newField := createField(r, c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			var highScore float64 = field[i][j].score
			newField[i][j].kind = field[i][j].kind

			// Check if the cell's neighbor is in the field and update the
			// cell's strategy
			if inField(field, i-1, j-1) == true {
				if field[i-1][j-1].score > highScore {
				highScore =  field[i-1][j-1].score
				newField[i][j].kind = field[i-1][j-1].kind
				}
			}
			if inField(field, i-1, j) == true {
				if field[i-1][j].score > highScore {
				highScore =  field[i-1][j].score
				newField[i][j].kind = field[i-1][j].kind
				}
			}
			if inField(field, i-1, j+1) == true {
				if field[i-1][j+1].score > highScore {
					highScore =  field[i-1][j+1].score
					newField[i][j].kind = field[i-1][j+1].kind
				}
			}
			if inField(field, i, j-1) == true {
				if field[i][j-1].score > highScore {
					highScore =  field[i][j-1].score
					newField[i][j].kind = field[i][j-1].kind
				}
			}
			if inField(field, i, j+1) == true {
				if field[i][j+1].score > highScore {
					highScore =  field[i][j+1].score
					newField[i][j].kind = field[i][j+1].kind
				}
			}
			if inField(field, i+1, j-1) == true {
				if field[i+1][j-1].score > highScore {
					highScore =  field[i+1][j-1].score
					newField[i][j].kind = field[i+1][j-1].kind
				}
			}
			if inField(field, i+1, j) == true {
				if field[i+1][j].score > highScore {
					highScore =  field[i+1][j].score
					newField[i][j].kind = field[i+1][j].kind
				}
			}
			if inField(field, i+1, j+1) == true {
				if field[i+1][j+1].score > highScore {
					highScore =  field[i+1][j+1].score
					newField[i][j].kind = field[i+1][j+1].kind
				}
			}
		}
	}

	return newField
}

// evolve takes an intial field and evolves it for nsteps according to the game
// rule. At each step, it should call "updateScores()" and the updateStrategies
func evolve(field [][]Cell, nsteps int, b float64) [][]Cell {
	for i := 0; i < nsteps; i++ {
		updateScores(field, b)
		field = updateStrategies(field)
	}
	return field
}

// Implements a Spatial Games version of prisoner's dilemma. The command-line
// usage is:
//     ./spatial field_file b nsteps
// where 'field_file' is the file continaing the initial arrangment of cells, b
// is the reward for defecting against a cooperator, and nsteps is the number
// of rounds to update stategies.
//
func main() {
	// Parse the command line
	if len(os.Args) != 4 {
		fmt.Println("Error: should input spatial field_file b nsteps.")
		return
	}

	fieldFile := os.Args[1]

	b, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil || b <= 0 {
		fmt.Println("Error: bad b parameter.")
		return
	}

	nsteps, err := strconv.Atoi(os.Args[3])
	if err != nil || nsteps < 0 {
		fmt.Println("Error: bad number of steps.")
		return
	}

    // Read the field
	field := readFieldFromFile(fieldFile)
    fmt.Println("Field dimensions are:", len(field), "by", len(field[0]))

    // Evolve the field and draw it as a PNG
	secondtofield := evolve(field, nsteps-1, b)
	field = evolve(secondtofield, 1, b)
	drawField(secondtofield, field, "Prisoners.png")
}