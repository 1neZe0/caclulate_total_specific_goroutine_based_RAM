package main

import (
	"crypto/rand"
	"fmt"
	"runtime"
	"sync"
	"syscall"
)

func main() {
	// Create a WaitGroup
	var wg sync.WaitGroup
	// Increment the WaitGroup counter to indicate that we have one goroutine
	wg.Add(1)
	// Start the goroutine
	//data := make([]byte, 0)
	go func() {
		defer wg.Done()
		// Get the number of gigabytes of RAM to fill
		var gigabytes int
		//fmt.Print("Enter the number of gigabytes to fill: ")
		//_, err := fmt.Scan(&gigabytes)
		//if err != nil {
		//	panic(err)
		//}
		gigabytes = 1

		// Calculate the number of bytes to fill
		bytes := int64(gigabytes) * int64(1<<30)

		// Allocate the bytes in a loop
		for i := int64(0); i < bytes; i += int64(1 << 20) {
			// Allocate a slice of 1MB
			b := make([]byte, 1<<20)
			// Fill the slice with random data
			_, err := rand.Read(b)
			if err != nil {
				panic(err)
			}
			// Keep the slice in memory by assigning it to a global variable
			// Note: This is not a good way to allocate large amounts of memory, as it will not be garbage collected.
			// It is only being used here for testing purposes.
			_ = b
		}

		// Print the current memory usage
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		fmt.Printf("Total memory usage: %d MB\n", mem.TotalAlloc/1024/1024)
	}()

	// Wait for the goroutine to finish
	wg.Wait()

	// Get the current memory stats
	var memStats runtime.MemStats
	var sysinfo syscall.Sysinfo_t
	_ = syscall.Sysinfo(&sysinfo)

	runtime.ReadMemStats(&memStats)

	// Calculate the memory usage of the goroutine
	goroutineMemoryUsage := memStats.StackInuse + memStats.MSpanInuse + memStats.MCacheInuse

	fmt.Printf("Memory usage of goroutine: %d bytes\n", goroutineMemoryUsage)

	// Calculate the maximum number of goroutines that can be run
	// based on the total amount of memory and the memory usage of the goroutine
	fmt.Printf("Total memory: %d bytes\n", sysinfo.Totalram)
	maxGoroutines := sysinfo.Totalram / goroutineMemoryUsage
	fmt.Printf("Memory usage of goroutine: %d bytes\n", memStats.Sys-memStats.HeapReleased)
	fmt.Printf("Maximum number of goroutines: %d\n", maxGoroutines)
}
