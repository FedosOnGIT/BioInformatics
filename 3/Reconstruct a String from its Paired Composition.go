package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
)

//region Stack

type Stack struct {
	values *[]*Node
	size   int
}

func (receiver *Stack) Push(node *Node) {
	if receiver.size == len(*receiver.values) {
		*receiver.values = append(*receiver.values, node)
	} else {
		(*receiver.values)[receiver.size] = node
	}
	receiver.size++
}

func (receiver *Stack) Peek() *Node {
	return (*receiver.values)[receiver.size-1]
}

func (receiver *Stack) Empty() bool {
	return receiver.size == 0
}

func (receiver *Stack) Pop() Node {
	var answer = (*receiver.values)[receiver.size-1]
	receiver.size--
	if receiver.size*2 < len(*receiver.values) {
		*receiver.values = (*receiver.values)[0:receiver.size]
	}
	return *answer
}

//endregion

//region Node

type Node struct {
	parts [2]string
	to    *[]Node
	from  *[]Node
	index int
}

func (receiver *Node) Degree() int {
	return len(*receiver.to) - receiver.index + len(*receiver.from)
}

func (receiver *Node) HasPath() bool {
	return len(*receiver.to) != receiver.index
}

func (receiver *Node) Remove() (Node, error) {
	if !receiver.HasPath() {
		return Node{}, errors.New("no more edges")
	}
	var position = receiver.index
	receiver.index++
	return (*receiver.to)[position], nil
}

func (receiver *Node) Connect(neighbour *Node) {
	*receiver.to = append(*receiver.to, *neighbour)
	*neighbour.from = append(*neighbour.from, *receiver)
}

//endregion

func FindEulerPath(graph map[[2]string]Node) Stack {
	var stack = Stack{values: &[]*Node{}, size: 0}
	for _, node := range graph {
		if node.Degree()%2 == 1 && node.HasPath() {
			stack.Push(&node)
			break
		}
	}
	var answer = Stack{values: &[]*Node{}, size: 0}
	for !stack.Empty() {
		var node = stack.Peek()
		if node.HasPath() {
			var neighbour, _ = node.Remove()
			stack.Push(&neighbour)
		} else {
			var last = stack.Pop()
			answer.Push(&last)
		}
	}
	return answer
}

func main() {
	fileIn, _ := os.Open("3/input.txt")
	in := bufio.NewReader(fileIn)
	fileOut, _ := os.Create("3/output.txt")
	out := bufio.NewWriter(fileOut)
	defer out.Flush()

	var length, gap int
	_, _ = fmt.Fscan(in, &length, &gap)

	var graph = make(map[[2]string]Node)
	var splitter = regexp.MustCompile("\\|")
	var number = 0
	for {
		var dna string
		var _, err = fmt.Fscan(in, &dna)
		if err != nil {
			break
		}
		var parts = splitter.Split(dna, 2)
		var prefix = [2]string{
			parts[0][0 : length-1],
			parts[1][0 : length-1],
		}
		var suffix = [2]string{
			parts[0][1:],
			parts[1][1:],
		}
		var prefixNode, prefixContains = graph[prefix]
		var suffixNode, suffixContains = graph[suffix]

		if !prefixContains {
			prefixNode = Node{parts: prefix, to: &[]Node{}, from: &[]Node{}, index: 0}
			graph[prefix] = prefixNode
		}
		if !suffixContains {
			suffixNode = Node{parts: suffix, to: &[]Node{}, from: &[]Node{}, index: 0}
			graph[suffix] = suffixNode
		}
		prefixNode.Connect(&suffixNode)
		number++
	}
	var path = FindEulerPath(graph)
	var answerLength = 2*length + gap + number - 1
	var answer []string
	for index := 0; index < answerLength; index++ {
		answer = append(answer, "X")
	}
	var index = 0
	for !path.Empty() {
		var node = path.Pop()
		for position, character := range node.parts[0] {
			answer[index+position] = string(character)
		}
		for position, character := range node.parts[1] {
			answer[index+length+gap+position] = string(character)
		}
		index++
	}
	for _, character := range answer {
		_, _ = fmt.Fprint(out, character)
	}
	_, _ = fmt.Fprintln(out)
}
