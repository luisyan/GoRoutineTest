package main

import (
	"time"
	"fmt"
	"sync"
)

func main() {
	testRepeat := 10
	maxGoRoutines := 100*10000
	resultWithDispatcher := make(map[int][]time.Duration)
	resultWithoutDispatcher := make(map[int][]time.Duration)
	averageWithDispatcher := make(map[int]time.Duration)
	averageWithoutDispatcher := make(map[int]time.Duration)
	mainRoutineBlocker := make(chan int)
	for i:=1;i<= maxGoRoutines;i*=10 {
		resultWithDispatcher[i] = make([]time.Duration,0)
		resultWithoutDispatcher[i] = make([]time.Duration,0)
		averageWithDispatcher[i] = time.Duration(0)
		averageWithoutDispatcher[i] = time.Duration(0)
	}
	execTest(test1, resultWithDispatcher, false, mainRoutineBlocker, testRepeat, maxGoRoutines, "with dispatcher")
	execTest(test2, resultWithoutDispatcher, true, mainRoutineBlocker, testRepeat, maxGoRoutines, "without dispatcher")
	<-mainRoutineBlocker
	for key,results := range resultWithDispatcher {
		sum1 := time.Duration(0)
		for _,v := range results {
			sum1 += v
		}
		average := sum1/time.Duration(testRepeat)
		averageWithDispatcher[key] = average
	}
	for key,results := range resultWithoutDispatcher {
		sum2 := time.Duration(0)
		for _,v := range results {
			sum2 += v
		}
		average := sum2/time.Duration(testRepeat)
		averageWithoutDispatcher[key] = average
	}
	for i:=1;i<= maxGoRoutines;i*=10 {
		fmt.Printf("\n")
		fmt.Printf("%v go routines: dispatcher=%v, no dispatcher=%v\n", i, averageWithDispatcher[i], averageWithoutDispatcher[i])
	}
}

func execTest(f func(int, map[int][]time.Duration, int, bool, bool, chan int), resultMap map[int][]time.Duration, exitAtEnd bool, exitChan chan int, numOuterLoop int, numInnerLoop int, testName string) {
	fmt.Printf("\n%v:\n", testName)
	for k:=1;k<=numOuterLoop;k++ {
		fmt.Printf("round %v ", k)
		for i:=1;i<=numInnerLoop;i*=10 {
			f(i, resultMap, numInnerLoop, exitAtEnd, k==numOuterLoop, exitChan)
		}
	}
}

func test1(nThreads int, resultMap1 map[int][]time.Duration, lastLoop int, exitAtEndLoop bool, lastOuterLoop bool, exitChan chan int) {
	wg := &sync.WaitGroup{}
	workers := make(chan Worker, nThreads)
	for i := 1; i <= nThreads; i++ {
		worker := newWorker(wg)
		worker.startWorking(i, workers)
	}
	timeStart := time.Now()
	go func(i int) {
		wg.Wait()
		duration := time.Since(timeStart)
		resultMap1[i] = append(resultMap1[i], duration)
		if i == lastLoop && exitAtEndLoop && lastOuterLoop {
			exitChan <- 1
		}
	}(nThreads)

	go func() {
		workerList := make([]Worker,0)
		for i := 1; i <= nThreads; i++ {
			selectedWorker := <- workers
			workerList = append(workerList, selectedWorker)
			selectedWorker.jobChan <- struct{}{}
		}
		for _,v := range workerList {
			v.exitChan <- struct{}{}
		}
	}()
}


func test2(nThreads int, resultMap2 map[int][]time.Duration, lastLoop int, exitAtEndLoop bool, lastOuterLoop bool, exitChan chan int) {
	wg := &sync.WaitGroup{}
	taskChan := make(chan struct{}, nThreads)
	robotQueue := make(chan Robot, nThreads)
	for i := 1; i <= nThreads; i++ {
		robot := newRobot(wg, taskChan)
		robot.startWorking(i, robotQueue)
	}
	timeStart := time.Now()
	go func(i int) {
		wg.Wait()
		duration := time.Since(timeStart)
		resultMap2[i] = append(resultMap2[i], duration)
		if i == lastLoop && exitAtEndLoop && lastOuterLoop {
			exitChan <- 1
		}
	}(nThreads)

	go func() {
		robotList := make([]Robot,0)
		for i := 1; i <= nThreads; i++ {
			selectedRobot := <- robotQueue
			robotList = append(robotList, selectedRobot)
			selectedRobot.jobChan <- struct{}{}
		}
		for _,v := range robotList {
			v.exitChan <- struct{}{}
		}
	}()
}



type Worker struct {
	jobChan chan struct{}
	exitChan chan struct{}
	waitG *sync.WaitGroup
}

type Robot struct {
	jobChan chan struct{}
	exitChan chan struct{}
	waitG *sync.WaitGroup
}

func (worker *Worker) startWorking(i int, workerQueue chan Worker) {
	worker.waitG.Add(1)
	workerQueue <- *worker
	go func() {
		for {
			select {
			case <-worker.jobChan:
			case <-worker.exitChan:
				worker.waitG.Done()
				return
			}
		}
	}()
}

func (robot *Robot) startWorking(i int, robotQueue chan Robot) {
	robot.waitG.Add(1)
	robotQueue <- *robot
	go func() {
		for {
			select {
			case <-robot.jobChan:
			case <-robot.exitChan:
				robot.waitG.Done()
				return
			}
		}
	}()
}

func newWorker(group *sync.WaitGroup) Worker {
	return Worker{make(chan struct{}), make(chan struct{}), group}
}

func newRobot(group *sync.WaitGroup, taskChan chan struct{}) Robot {
	return Robot{taskChan, make(chan struct{}), group}
}