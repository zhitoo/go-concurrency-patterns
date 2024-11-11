package main

import (
	"fmt"
	"sync"
	"time"
)

/*
What is the Worker Pool Pattern?

The worker pool pattern involves creating a group of worker goroutines to process tasks concurrently,
limiting the number of simultaneous operations => in this example limit to 3 worker
This pattern is valuable when you have a large number of tasks to execute.

Some examples of using the Worker Pool Pattern in Real-world Applications:
	- Handling incoming HTTP requests in a web server.
	- Processing images concurrently.
*/

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(time.Second * 1)
		results <- job * 2
	}
}

func main() {
	numJobs := 10
	numWorkers := 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	/*
		Using a WaitGroup in the original code allows you to:

			- Wait until all workers finish processing.
			- Safely close the results channel,
			  which signals to the for range loop in the main function
			  that no more results will come.
	*/
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			worker(workerID, jobs, results)
		}(i)
	}

	// Enqueue jobs
	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	//the main goroutine will keep trying to read from results until it is closed
	//Notice: The reason for the deadlock is because the results channel is never closed,
	// and the main function is trying to range over it indefinitely.
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}
