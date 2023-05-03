package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Direction int

const (
	Head Direction = iota
	Tail
)

type Colour int

const (
	Red Colour = iota
	Blue
)

type Node struct {
	value     int
	direction Direction
	neighbour int
	red       int
	blue      int
}

type Graph struct {
	nodes  []Node
	number int
}

func Indexes(nodeValue, nextValue int) (int, int) {
	var node, next int
	var nodeIndex = Abs(nodeValue) - 1
	var nextIndex = Abs(nextValue) - 1

	if nodeValue > 0 {
		node = nodeIndex*2 + 1
	} else {
		node = nodeIndex * 2
	}

	if nextValue > 0 {
		next = nextIndex * 2
	} else {
		next = nextIndex*2 + 1
	}
	return node, next
}

func CreateGraph(red []int, blue []int) Graph {
	var nodes []Node
	var number = len(red)
	for index := 0; index < number; index++ {
		var head = Node{value: index + 1, direction: Head, neighbour: index*2 + 1}
		var tail = Node{value: index + 1, direction: Tail, neighbour: index * 2}
		nodes = append(nodes, head)
		nodes = append(nodes, tail)
	}
	for index, nodeValue := range red {
		var nextValue = red[(index+1)%number]
		var node, next = Indexes(nodeValue, nextValue)
		nodes[node].red = next
		nodes[next].red = node
	}
	for index, nodeValue := range blue {
		var nextValue = blue[(index+1)%number]
		var node, next = Indexes(nodeValue, nextValue)
		nodes[node].blue = next
		nodes[next].blue = node
	}
	return Graph{nodes: nodes, number: number}
}

func (receiver *Graph) FindCycles() [][]int {
	var used []bool
	var length = len(receiver.nodes)
	for index := 0; index < length; index++ {
		used = append(used, false)
	}
	var cycles [][]int
	for index, _ := range receiver.nodes {
		var cycle []int
		var currentIndex = index
		var colour = Blue
		for !used[currentIndex] {
			var current = receiver.nodes[currentIndex]
			cycle = append(cycle, currentIndex)
			used[currentIndex] = true
			if colour == Red {
				currentIndex = current.red
				colour = Blue
			} else {
				currentIndex = current.blue
				colour = Red
			}
		}
		if len(cycle) > 2 {
			cycles = append(cycles, cycle)
		}
	}
	return cycles
}

func (receiver *Graph) Break(cycle []int) {
	var start = cycle[0]
	var finish = cycle[1]
	var left = cycle[len(cycle)-1]
	var right = cycle[2]

	receiver.nodes[start].red = finish
	receiver.nodes[finish].red = start

	receiver.nodes[left].red = right
	receiver.nodes[right].red = left
}

func (receiver *Graph) Println(writer *bufio.Writer) {
	var processed []bool
	for index := 0; index < receiver.number; index++ {
		processed = append(processed, false)
	}
	var answer = ""
	for index := 0; index < receiver.number; index++ {
		if processed[index] {
			continue
		}
		var start = receiver.nodes[index*2]
		var finish = receiver.nodes[index*2+1]
		var cycle = ""
		for !processed[start.value-1] {
			var sign = ""
			if start.direction == Head {
				sign = "+"
			} else {
				sign = "-"
			}
			cycle = cycle + sign + strconv.Itoa(start.value) + " "
			processed[start.value-1] = true
			start = receiver.nodes[finish.red]
			finish = receiver.nodes[start.neighbour]
		}
		answer = answer + "(" + strings.Trim(cycle, " ") + ")"
	}
	_, _ = fmt.Fprintln(writer, answer)
}

func main() {
	fileIn, _ := os.Open("5/input.txt")
	scanner := bufio.NewScanner(fileIn)
	fileOut, _ := os.Create("5/output.txt")
	out := bufio.NewWriter(fileOut)
	defer fileIn.Close()
	defer fileOut.Close()
	defer out.Flush()

	scanner.Scan()
	var start = scanner.Text()
	scanner.Scan()
	var finish = scanner.Text()
	var red, blue []int

	redReader := strings.NewReader(start[1 : len(start)-1])
	for {
		var redValue int
		var _, err = fmt.Fscan(redReader, &redValue)
		if err != nil {
			break
		}
		red = append(red, redValue)
	}

	blueReader := strings.NewReader(finish[1 : len(finish)-1])
	for {
		var blueValue int
		var _, err = fmt.Fscan(blueReader, &blueValue)
		if err != nil {
			break
		}
		blue = append(blue, blueValue)
	}

	var graph = CreateGraph(red, blue)
	graph.Println(out)
	for {
		var cycles = graph.FindCycles()
		if len(cycles) == 0 {
			break
		}
		for _, cycle := range cycles {
			graph.Break(cycle)
			graph.Println(out)
		}
	}
}
