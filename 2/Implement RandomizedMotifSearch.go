package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
)

var i = 0

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Profile(length int, motifs []string) map[string][]float64 {
	var number = float64(len(motifs)) + 4
	var profiles = map[string][]float64{
		"A": {},
		"C": {},
		"G": {},
		"T": {},
	}
	for index := 0; index < length; index++ {
		for key, profile := range profiles {
			profiles[key] = append(profile, 1)
		}
		for _, motif := range motifs {
			var key = string(motif[index])
			var profile, _ = profiles[key]
			profile[index]++
			profiles[key] = profile
		}
		for key, profile := range profiles {
			profile[index] = profile[index] / number
			profiles[key] = profile
		}
	}
	return profiles
}

func Motifs(profiles map[string][]float64, dnas []string, length int) []string {
	var motifs []string
	for _, dna := range dnas {
		var bestMotif = ""
		var bestResult = 0.0
		for index := 0; index <= len(dna)-length; index++ {
			var motif = dna[index : index+length]
			var result = 1.0
			for position := 0; position < length; position++ {
				var frequency, _ = profiles[string(motif[position])]
				result *= frequency[position]
			}
			if result > bestResult {
				bestMotif = motif
				bestResult = result
			}
		}
		motifs = append(motifs, bestMotif)
	}
	return motifs
}

func Score(motifs []string, length int) int {
	var score = 0
	var dnaNumber = len(motifs)
	for index := 0; index < length; index++ {
		var letters = make(map[uint8]int)
		for _, motif := range motifs {
			var key = motif[index]
			var letter, contains = letters[key]
			if contains {
				letters[key] = letter + 1
			} else {
				letters[key] = 1
			}
		}
		var maximal = 0
		for _, number := range letters {
			maximal = max(maximal, number)
		}
		score += dnaNumber - maximal
	}
	return score
}

func RandomizedMotifSearch(dnas []string, length int) []string {
	var motifs []string
	for _, dna := range dnas {
		var start = rand.Intn(len(dna) - length)
		var motif = dna[start : start+length]
		motifs = append(motifs, motif)
		i++
		i %= len(dna) - length
	}
	var bestMotifs = motifs
	var bestScore = Score(bestMotifs, length)
	for {
		var profiles = Profile(length, bestMotifs)
		motifs = Motifs(profiles, dnas, length)
		var score = Score(motifs, length)
		if score < bestScore {
			bestMotifs = motifs
			bestScore = score
		} else {
			return bestMotifs
		}
	}
}

func main() {
	fileIn, _ := os.Open("2/input.txt")
	in := bufio.NewReader(fileIn)
	fileOut, _ := os.Create("2/output.txt")
	out := bufio.NewWriter(fileOut)
	defer out.Flush()

	var (
		motifLength int
		dnaNumber   int
	)
	_, _ = fmt.Fscan(in, &motifLength, &dnaNumber)
	var dnas []string
	for index := 0; index < dnaNumber; index++ {
		var dna string
		_, _ = fmt.Fscan(in, &dna)
		dnas = append(dnas, dna)
	}
	var bestMotifs []string
	var bestScore = math.MaxInt
	for index := 0; index <= 1000; index++ {
		var motifs = RandomizedMotifSearch(dnas, motifLength)
		var score = Score(motifs, motifLength)
		if score < bestScore {
			bestMotifs = motifs
			bestScore = score
		}
	}
	for _, motif := range bestMotifs {
		_, _ = fmt.Fprintln(out, motif)
	}
	_, _ = fmt.Fprintln(out, bestScore)
}
