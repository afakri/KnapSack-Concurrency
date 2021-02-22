package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Driver code
func main() {

	fmt.Println("Number of cores: ", runtime.NumCPU())
	W := 7
	weights := make([]int, 0)
	values := make([]int, 0)
	readFile("p1.txt", &W, &weights, &values)
	start := time.Now()
	results := make(chan int, 2)
	count := 0
	KnapSack(W, weights, values, &results, &count)
	fmt.Println(<-results)
	end := time.Now()
	fmt.Printf("Total runtime: %s\n", end.Sub(start))

}
func readFile(fileName string, W *int, wt *[]int, val *[]int) {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	j := scanner.Text()
	i, _ := strconv.Atoi(j)
	for scanner.Scan() {
		if i == 0 {
			*W, _ = strconv.Atoi(scanner.Text())

		} else {
			line := strings.Split(scanner.Text(), " ")
			value, _ := strconv.Atoi(line[1])
			*val = append(*val, value)
			weight, _ := strconv.Atoi(line[2])
			*wt = append(*wt, weight)
			i--
		}

	}

}
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func KnapSack(W int, wt []int, val []int, ch *chan int, count *int) {
	*count++

	// Base Case
	if len(wt) == 0 || W == 0 {
		*ch <- 0
		return
	}
	last := len(wt) - 1

	// If weight of the nth item is more
	// than Knapsack capacity W, then
	// this item cannot be included
	// in the optimal solution
	if wt[last] > W {
		if *count > len(wt)/4 {
			KnapSack(W, wt[:last], val[:last], ch, count)
		} else {
			go KnapSack(W, wt[:last], val[:last], ch, count)
		}

		// Return the maximum of two cases:
		// (1) nth item included
		// (2) item not included
	} else {
		included := make(chan int, 2)
		not := make(chan int, 1)
		x := val[last]
		if *count > len(wt)/4 {
			KnapSack(W-wt[last], wt[:last], val[:last], &included, count)
			KnapSack(W, wt[:last], val[:last], &not, count)
		} else {
			go KnapSack(W-wt[last], wt[:last], val[:last], &included, count)
			go KnapSack(W, wt[:last], val[:last], &not, count)
		}

		x += <-included
		y := <-not
		*ch <- Max(x, y)

	}
}
