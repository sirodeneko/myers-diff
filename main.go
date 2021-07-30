package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

type Operation uint

type Diff struct {
	Op  Operation
	str string
}

const (
	INSERT Operation = 1
	DELETE           = 2
	MOVE             = 3
)

func (op Operation) String() string {
	switch op {
	case INSERT:
		return "INS"
	case DELETE:
		return "DEL"
	case MOVE:
		return "MOV"
	default:
		return "UNKNOWN"
	}
}

var colors = map[Operation]string{
	INSERT: "\033[32m",
	DELETE: "\033[31m",
	MOVE:   "\033[39m",
}

func main() {
	flag.Parse()

	var src []string
	var dst []string

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: myers-diff [src_file] [dst_file]")
		os.Exit(1)
	}

	var err error

	src, err = getFileLines(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	dst, err = getFileLines(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	generateDiffAndPrint(src, dst)
}

func generateDiffAndPrint(src, dst []string) {
	script := shortestEditScript(src, dst)

	srcIndex, dstIndex := 0, 0

	for _, op := range script {
		switch op {
		case INSERT:
			fmt.Println(colors[op] + "+" + dst[dstIndex])
			dstIndex += 1

		case MOVE:
			fmt.Println(colors[op] + " " + src[srcIndex])
			srcIndex += 1
			dstIndex += 1

		case DELETE:
			fmt.Println(colors[op] + "-" + src[srcIndex])
			srcIndex += 1
		}
	}
}

func generateDiff(src, dst []string) []Diff {
	script := shortestEditScript(src, dst)

	diff := make([]Diff, len(script))
	srcIndex, dstIndex, diffIndex := 0, 0, 0

	for _, op := range script {
		diff[diffIndex].Op = op
		switch op {
		case INSERT:
			diff[diffIndex].str = dst[dstIndex]
			dstIndex += 1

		case MOVE:
			diff[diffIndex].str = src[srcIndex]
			srcIndex += 1
			dstIndex += 1

		case DELETE:
			diff[diffIndex].str = src[srcIndex]
			srcIndex += 1
		}
		diffIndex++
	}
	return diff
}

// 生成最短的编辑脚本
func shortestEditScript(src, dst []string) []Operation {
	n := len(src)
	m := len(dst)
	max := n + m
	var trace []map[int]int
	var x, y int
	// 反向回溯
	var script []Operation

	if max <= 0 {
		return script
	}

loop:
	for d := 0; d <= max; d++ {
		// 最多只有 d+1 个 k
		v := make(map[int]int, d+2)
		trace = append(trace, v)

		// 需要注意处理对角线
		if d == 0 {
			t := 0
			for len(src) > t && len(dst) > t && src[t] == dst[t] {
				t++
			}
			v[0] = t
			if t == len(src) && t == len(dst) {
				break loop
			}
			continue
		}

		lastV := trace[d-1]

		for k := -d; k <= d; k += 2 {
			// 向下
			if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
				x = lastV[k+1]
			} else { // 向右
				x = lastV[k-1] + 1
			}

			y = x - k

			// 处理对角线
			for x < n && y < m && src[x] == dst[y] {
				x, y = x+1, y+1
			}

			v[k] = x

			if x == n && y == m {
				break loop
			}
		}
	}

	x = n
	y = m
	var k, prevK, prevX, prevY int

	for d := len(trace) - 1; d > 0; d-- {
		k = x - y
		lastV := trace[d-1]

		if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX = lastV[prevK]
		prevY = prevX - prevK

		for x > prevX && y > prevY {
			script = append(script, MOVE)
			x -= 1
			y -= 1
		}

		if x == prevX {
			script = append(script, INSERT)
		} else {
			script = append(script, DELETE)
		}

		x, y = prevX, prevY
	}

	if trace[0][0] != 0 {
		for i := 0; i < trace[0][0]; i++ {
			script = append(script, MOVE)
		}
	}

	return reverse(script)
}

func printTrace(trace []map[int]int) {
	for d := 0; d < len(trace); d++ {
		fmt.Printf("d = %d:\n", d)
		v := trace[d]
		for k := -d; k <= d; k += 2 {
			x := v[k]
			y := x - k
			fmt.Printf("  k = %2d: (%d, %d)\n", k, x, y)
		}
	}
}

func reverse(s []Operation) []Operation {
	result := make([]Operation, len(s))

	for i, v := range s {
		result[len(s)-1-i] = v
	}

	return result
}

func getFileLines(p string) ([]string, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
