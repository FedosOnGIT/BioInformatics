package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type Action int

const (
	Insertion Action = iota
	Deletion
	Letter
)

type Cell struct {
	action Action
	value  int
	i      int
	j      int
}

var (
	Open     = 11
	Continue = 1
)

func readMatrix(scanner *bufio.Scanner) map[string]map[string]int {
	var matrix = make(map[string]map[string]int)
	var line string
	scanner.Scan()
	line = scanner.Text()
	var letters []string
	var lineReader = strings.NewReader(line)
	for {
		var letter string
		var _, err = fmt.Fscan(lineReader, &letter)
		if err != nil {
			break
		}
		letters = append(letters, letter)
	}
	for scanner.Scan() {
		line = scanner.Text()
		lineReader = strings.NewReader(line)
		var firstLetter string
		_, _ = fmt.Fscan(lineReader, &firstLetter)
		matrix[firstLetter] = make(map[string]int)
		for _, secondLetter := range letters {
			var score int
			_, _ = fmt.Fscan(lineReader, &score)
			matrix[firstLetter][secondLetter] = score
		}
	}
	return matrix
}

func main() {
	costFile, _ := os.Open("4/cost_matrix.txt")
	costReader := bufio.NewScanner(costFile)
	var costMatrix = readMatrix(costReader)

	fileIn, _ := os.Open("4/input.txt")
	in := bufio.NewReader(fileIn)
	fileOut, _ := os.Create("4/output.txt")
	out := bufio.NewWriter(fileOut)
	defer fileIn.Close()
	defer fileOut.Close()
	defer out.Flush()
	defer costFile.Close()

	var first, second string
	_, _ = fmt.Fscan(in, &first, &second)

	var insertion, deletion, matrix [][]Cell
	for i := 0; i <= len(first); i++ {
		var insertionRow, deletionRow, matrixRow []Cell
		for j := 0; j <= len(second); j++ {
			insertionRow = append(insertionRow, Cell{Insertion, math.MinInt / 2, i - 1, j})
			deletionRow = append(deletionRow, Cell{Deletion, math.MinInt / 2, i, j - 1})
			matrixRow = append(matrixRow, Cell{Letter, math.MinInt / 2, i - 1, j - 1})
		}
		insertion = append(insertion, insertionRow)
		deletion = append(deletion, deletionRow)
		matrix = append(matrix, matrixRow)
	}
	matrix[0][0].value = 0
	for i := 0; i <= len(first); i++ {
		for j := 0; j <= len(second); j++ {
			if i != 0 {
				var one = insertion[i-1][j].value - Continue
				var two = matrix[i-1][j].value - Open
				if one > two {
					insertion[i][j] = Cell{Insertion, one, i - 1, j}
				} else {
					insertion[i][j] = Cell{Letter, two, i - 1, j}
				}
			}
			if j != 0 {
				var one = deletion[i][j-1].value - Continue
				var two = matrix[i][j-1].value - Open
				if one > two {
					deletion[i][j] = Cell{Deletion, one, i, j - 1}
				} else {
					deletion[i][j] = Cell{Letter, two, i, j - 1}
				}
			}
			if i != 0 || j != 0 {
				var one = insertion[i][j].value
				var two = deletion[i][j].value
				if one > two {
					matrix[i][j] = Cell{Insertion, one, i, j}
				} else {
					matrix[i][j] = Cell{Deletion, two, i, j}
				}
			}
			if i != 0 && j != 0 {
				var firstLetter = string(first[i-1])
				var secondLetter = string(second[j-1])
				var three = matrix[i-1][j-1].value + costMatrix[firstLetter][secondLetter]
				if three > matrix[i][j].value {
					matrix[i][j] = Cell{Letter, three, i - 1, j - 1}
				}
			}
		}
	}
	var mapping = map[Action][][]Cell{
		Insertion: insertion,
		Deletion:  deletion,
		Letter:    matrix,
	}
	var current = matrix[len(first)][len(second)]
	_, _ = fmt.Fprintln(out, current.value)
	var currentAction = Letter
	var firstAnswer, secondAnswer []string
	var firstIndex = len(first) - 1
	var secondIndex = len(second) - 1
	for current.i >= 0 || current.j >= 0 {
		switch currentAction {
		case Insertion:
			currentAction = current.action
			current = mapping[currentAction][current.i][current.j]
			firstAnswer = append(firstAnswer, string(first[firstIndex]))
			secondAnswer = append(secondAnswer, "-")
			firstIndex--
		case Deletion:
			currentAction = current.action
			current = mapping[currentAction][current.i][current.j]
			firstAnswer = append(firstAnswer, "-")
			secondAnswer = append(secondAnswer, string(second[secondIndex]))
			secondIndex--
		case Letter:
			currentAction = current.action
			current = mapping[currentAction][current.i][current.j]
			if currentAction == Letter {
				firstAnswer = append(firstAnswer, string(first[firstIndex]))
				secondAnswer = append(secondAnswer, string(second[secondIndex]))
				firstIndex--
				secondIndex--
			}
		}
	}
	for index := len(firstAnswer) - 1; index >= 0; index-- {
		_, _ = fmt.Fprint(out, firstAnswer[index])
	}
	_, _ = fmt.Fprintln(out)
	for index := len(secondAnswer) - 1; index >= 0; index-- {
		_, _ = fmt.Fprint(out, secondAnswer[index])
	}
	_, _ = fmt.Fprintln(out)
}
