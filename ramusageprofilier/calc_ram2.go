package main

// #include "calc_ram2.c"
import "C"
import (
	"crypto/rand"
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

type MemoryUsageData struct {
	cData *C.struct_MemoryUsageData
}

// NewMemoryUsageData creates a new MemoryUsageData struct and returns a pointer to it.
func NewMemoryUsageData() *MemoryUsageData {
	cData := C.malloc(C.sizeof_struct_MemoryUsageData)
	return &MemoryUsageData{cData: (*C.struct_MemoryUsageData)(cData)}
}

// Free frees the memory allocated for the MemoryUsageData struct.
func (data *MemoryUsageData) Free() {
	C.free(unsafe.Pointer(data.cData))
}

// StartMemoryUsageThread starts the memory usage thread.
// Returns 0 on success and a non-zero value on failure.
func (data *MemoryUsageData) StartMemoryUsageThread() int {
	return int(C.startMemoryUsageThread(data.cData))
}

// StopMemoryUsageThread stops the memory usage thread.
// Returns 0 on success and a non-zero value on failure.
func (data *MemoryUsageData) StopMemoryUsageThread() int {
	return int(C.stopMemoryUsageThread(data.cData))
}

// CalculateMemoryUsage calculates the maximum memory usage.
func (data *MemoryUsageData) CalculateMemoryUsage() uint64 {
	return uint64(C.calculate_memory_usage(data.cData))
}

// CalculatePercentMemoryUsage calculates the percentage of memory usage.
func (data *MemoryUsageData) CalculatePercentMemoryUsage() float64 {
	return float64(C.calculate_percent_memory_usage(data.cData))
}

// GetTotalMemory returns the total memory available on the system.
func (data *MemoryUsageData) GetTotalMemory() uint64 {
	return uint64(C.get_total_memory(data.cData))
}

// GetFreeMemory returns the amount of free memory on the system.
func (data *MemoryUsageData) GetFreeMemory() uint64 {
	return uint64(C.get_free_memory(data.cData))
}

// GetUsedMemory returns the amount of used memory on the system.
func (data *MemoryUsageData) GetUsedMemory() uint64 {
	return uint64(C.get_used_memory(data.cData))
}

// CalculatePercentFreeMemory calculates the percentage of free memory on the system.
func (data *MemoryUsageData) CalculatePercentFreeMemory() float64 {
	return float64(C.calculate_percent_free_memory(data.cData))
}

// CalculatePercentUsedMemory calculates the percentage of used memory on the system.
func (data *MemoryUsageData) CalculatePercentUsedMemory() float64 {
	return float64(C.calculate_percent_used_memory(data.cData))
}

func main() {

	ram := NewMemoryUsageData()
	defer ram.Free()

	ram.StartMemoryUsageThread()

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
	ram.StopMemoryUsageThread()

	fmt.Println("Memory usage:", ram.CalculateMemoryUsage())
	fmt.Println("Percent memory usage:", ram.CalculatePercentMemoryUsage())
	fmt.Println("Total memory:", ram.GetTotalMemory())
	fmt.Println("Free memory:", ram.GetFreeMemory())
	fmt.Println("Used memory:", ram.GetUsedMemory())
	fmt.Println("Percent free memory:", ram.CalculatePercentFreeMemory())
	fmt.Println("Percent used memory:", ram.CalculatePercentUsedMemory())
}
