package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
)

var Gens = map[int]string{
	0: "A",
	1: "C",
	2: "T",
	3: "G",
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func pow(base, degree int) int {
	if degree == 0 {
		return 0
	}
	if degree == 1 {
		return base
	}
	var part = pow(base, degree/2)
	var answer = part * part
	if degree%2 == 1 {
		answer *= base
	}
	return answer
}

func HammingDistance(first string, second string) (int, error) {
	if len(first) != len(second) {
		return 0, errors.New("input strings lengths must be equal")
	}
	var distance = 0
	for index := 0; index < len(first); index++ {
		if first[index] != second[index] {
			distance++
		}
	}
	return distance, nil
}

func generate(index int, length int) string {
	var builder strings.Builder
	for i := 0; i < length; i++ {
		builder.WriteString(Gens[index%4])
		index /= 4
	}
	return builder.String()
}

func d(pattern string, dna string) int {
	var length = len(pattern)
	var minimal = math.MaxInt
	for index := 0; index <= len(dna)-length; index++ {
		var part = dna[index : index+length]
		var distance, _ = HammingDistance(pattern, part)
		minimal = min(minimal, distance)
	}
	return minimal
}

func main() {
	file, _ := os.Open("1/input.txt")
	in := bufio.NewReader(file)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var patternLength int
	_, _ = fmt.Fscan(in, &patternLength)

	var dnas []string
	for {
		var dna string
		var _, err = fmt.Fscan(in, &dna)
		if err != nil {
			break
		}
		dnas = append(dnas, dna)
	}
	var bestPattern = ""
	var bestDistance = math.MaxInt
	var variants = pow(4, patternLength)
	for variant := 0; variant < variants; variant++ {
		var pattern = generate(variant, patternLength)
		var dist = 0
		for _, dna := range dnas {
			dist += d(pattern, dna)
		}
		if bestDistance > dist {
			bestDistance = dist
			bestPattern = pattern
		}
	}
	_, _ = fmt.Fprintln(out, bestPattern)
}
