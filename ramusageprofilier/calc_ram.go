package main

/*
#include <stdlib.h>
#include "calc_ram.c"
*/
import "C"
import (
	"crypto/rand"
	"fmt"
	"runtime"
	"sync"
)

type CalcRAM struct {
	cRam *C.struct_CalcRAM
}

func NewCalcRAM() *CalcRAM {
	ram := &CalcRAM{cRam: C.new_CalcRAM()}
	return ram
}

func (ram *CalcRAM) StartMeasuring() {
	C.start_measuring(ram.cRam)
}

func (ram *CalcRAM) StopMeasuring() {
	C.stop_measuring(ram.cRam)
}

func (ram *CalcRAM) CalculateMemoryUsage() uint64 {
	return uint64(C.calculate_memory_usage(ram.cRam))
}

func (ram *CalcRAM) CalculatePercentMemoryUsage() float64 {
	return float64(C.calculate_percent_memory_usage(ram.cRam))
}

func (ram *CalcRAM) GetTotalMemory() uint64 {
	return uint64(C.get_total_memory(ram.cRam))
}

func (ram *CalcRAM) GetFreeMemory() uint64 {
	return uint64(C.get_free_memory(ram.cRam))
}

func (ram *CalcRAM) GetUsedMemory() uint64 {
	return uint64(C.get_used_memory(ram.cRam))
}

func (ram *CalcRAM) CalculatePercentFreeMemory() float64 {
	return float64(C.calculate_percent_free_memory(ram.cRam))
}

func (ram *CalcRAM) CalculatePercentUsedMemory() float64 {
	return float64(C.calculate_percent_used_memory(ram.cRam))
}

func (ram *CalcRAM) StopCalc() {
	C.stop_calc(ram.cRam)
}

func main() {
	ram := NewCalcRAM()
	defer ram.StopCalc()

	ram.StartMeasuring()

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
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Printf("Total memory usage: %d MB\n", mem.TotalAlloc/1024/1024)
	ram.StopMeasuring()

	fmt.Println("Memory usage:", ram.CalculateMemoryUsage())
	fmt.Println("Percent memory usage:", ram.CalculatePercentMemoryUsage())
	fmt.Println("Total memory:", ram.GetTotalMemory())
	fmt.Println("Free memory:", ram.GetFreeMemory())
	fmt.Println("Used memory:", ram.GetUsedMemory())
	fmt.Println("Percent free memory:", ram.CalculatePercentFreeMemory())
	fmt.Println("Percent used memory:", ram.CalculatePercentUsedMemory())
}
