package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
func main() {

	// make a buffer chan between sender and receiver
	ch := make(chan int, 3)
	// make a slice of task id
	tasks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 3, 44, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70}
	// mask a map of running task id
	waiting := make(map[int]int)
	// make a slice to store task done
	done := []int{}
	// define max worker count
	const workers = 50

	// make a map to track worker done status
	workerIds := make(map[int]int, workers)
	// sender will send 1 package per 2 second via child thread
	for id := 0; id < workers; id++ {
		go func(id int) {
			workerIds[id] = id
			for {
				if len(tasks) == 0 {
					delete(workerIds, id)
					if len(workerIds) == 0 {
						close(ch)
					}
					break
				}
				var i int = tasks[0]
				for ok := true; ok; {
					i = tasks[0]
					tasks = tasks[1:]
					_, ok = waiting[i]
				}
				waiting[i] = i
				speed := randInt(10, 30)
				for p := 0; p < 100; p++ {
					fmt.Printf("Worker:[id=%2d] Task:[id=%2d]\t %-100s[%3d%%]\n", id, i, string(strings.Repeat("-", p)+">"), p)
					time.Sleep(time.Duration(speed) * time.Millisecond)
				}

				fmt.Printf("Worker:[id=%2d] Task:[id=%2d]\t %-100s[%3d%%]\n", id, i, string(strings.Repeat("-", 100)+">"), 100)
				delete(waiting, i)
				ch <- i

			}
		}(id)
	}

	// receiver will receive 1 package per 1 millisecond via main thread
	for {
		if task, ok := <-ch; ok {
			done = append(done, task)
			sort.Ints(done)
			fmt.Println("main thread status\ntasks:", tasks, "\ndone", done)
		} else {
			fmt.Println("main thread done\n")
			break
		}
		time.Sleep(time.Millisecond)
	}
}
