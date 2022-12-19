package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Create a WaitGroup
	var wg sync.WaitGroup
	// Increment the WaitGroup counter to indicate that we have one goroutine
	wg.Add(1)
	// Start the goroutine
	go func() {

		defer wg.Done()
		f, err := os.OpenFile("test.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		auth := []string{"hello:hello", "bye:bye", "salam:salam"}
		proxy := ""
		// Try connecting to the proxy without a username and password
		_, err = ssh.Dial("tcp", proxy, &ssh.ClientConfig{
			User: "",
			Auth: []ssh.AuthMethod{ssh.Password("")},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: time.Second * 30,
		})
		if err == nil {
			fmt.Printf("good ssh = " + proxy + "\n")
			// If the connection is successful, write the proxy's address to the good_proxies file
			_, err = f.WriteString(proxy + "\n")
			if err != nil {
				fmt.Printf("Error writing to good_proxies file: %s\n", err)
			}
			return
		}
		fmt.Println("1")
		for _, a := range auth {
			pair := strings.Split(a, ":")
			username := pair[0]
			password := pair[1]

			_, err := ssh.Dial("tcp", proxy, &ssh.ClientConfig{
				User: username,
				Auth: []ssh.AuthMethod{ssh.Password(password)},
				HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
					return nil
				},
				Timeout: time.Second * 30,
			})
			if err == nil {
				fmt.Printf("good ssh = " + proxy + ":" + username + ":" + password + "\n")
				// If the connection is successful, write the proxy's address to the good_proxies file
				_, _ = f.WriteString(proxy + ":" + username + ":" + password + "\n")
				return
			}
			time.Sleep(1)
		}
		fmt.Println("2")
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

	fmt.Printf("Maximum number of goroutines: %d\n", maxGoroutines)
}
