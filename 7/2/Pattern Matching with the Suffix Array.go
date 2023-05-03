package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type SuffixArray struct {
	suffix []int
	dna    string
}

func CreateSuffixArray(dna string) SuffixArray {
	dna = dna + "$"
	var suffix []int
	var compare []int
	for index := 0; index < len(dna); index++ {
		suffix = append(suffix, index)
		compare = append(compare, int(dna[index]))
	}
	sort.Slice(suffix, func(i, j int) bool {
		return dna[i] < dna[j]
	})
	var total = len(dna)
	for length := 1; length < total; length *= 2 {
		sort.Slice(suffix, func(i, j int) bool {
			var first = suffix[i]
			var second = suffix[j]
			if compare[first] != compare[second] {
				return compare[first] < compare[second]
			} else {
				return compare[(first+length)%total] < compare[(second+length)%total]
			}
		})
		var compareNew = make([]int, len(compare))
		copy(compareNew, compare)
		compareNew[suffix[0]] = 0
		for index := 1; index < total; index++ {
			var left1 = suffix[index-1]
			var right1 = (suffix[index-1] + length) % total
			var left2 = suffix[index]
			var right2 = (suffix[index] + length) % total
			if compare[left1] != compare[left2] || compare[right1] != compare[right2] {
				compareNew[suffix[index]] = compareNew[suffix[index-1]] + 1
			} else {
				compareNew[suffix[index]] = compareNew[suffix[index-1]]
			}
		}
		compare = compareNew
	}
	return SuffixArray{dna: dna, suffix: suffix}
}

func (receiver *SuffixArray) Find(pattern string) []int {
	var left = 0
	var right = len(receiver.dna) - 1
	for index := 0; index < len(pattern); index++ {
		var l1 = left - 1
		var r1 = right
		for r1-l1 > 1 {
			var middle = (l1 + r1) / 2
			var position = receiver.suffix[middle] + index
			if receiver.dna[position] >= pattern[index] {
				r1 = middle
			} else {
				l1 = middle
			}
		}

		var l2 = left
		var r2 = right + 1
		for r2-l2 > 1 {
			var middle = (l2 + r2) / 2
			var position = receiver.suffix[middle] + index
			if receiver.dna[position] <= pattern[index] {
				l2 = middle
			} else {
				r2 = middle
			}
		}

		left = r1
		right = l2
	}
	var positions []int
	for index := left; index <= right; index++ {
		positions = append(positions, receiver.suffix[index])
	}
	return positions
}

func main() {
	fileIn, _ := os.Open("7/2/input.txt")
	in := bufio.NewReader(fileIn)
	fileOut, _ := os.Create("7/2/output.txt")
	out := bufio.NewWriter(fileOut)
	defer fileIn.Close()
	defer fileOut.Close()
	defer out.Flush()

	var dna string
	_, _ = fmt.Fscan(in, &dna)
	var suffixArray = CreateSuffixArray(dna)
	var positions []int
	for {
		var pattern string
		var _, err = fmt.Fscan(in, &pattern)
		if err != nil {
			break
		}
		var patternPositions = suffixArray.Find(pattern)
		for _, position := range patternPositions {
			positions = append(positions, position)
		}
	}
	for _, position := range positions {
		_, _ = fmt.Fprint(out, position)
		_, _ = fmt.Fprint(out, " ")
	}
	_, _ = fmt.Fprintln(out)
}
