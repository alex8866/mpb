package mpb_test

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/vbauerster/mpb"
)

func Example() {
	// Star mpb's rendering goroutine.
	p := mpb.New()
	// Set custom width for every bar, which mpb will contain
	// The default one is 80
	p.SetWidth(100)
	// Set custom format for every bar, the default one is "[=>-]"
	p.Format("╢▌▌░╟")
	// Set custom refresh rate, the default one is 100 ms
	p.RefreshRate(120 * time.Millisecond)

	// Add a bar. You're not limited to just one bar, add many if you need.
	bar := p.AddBar(100).
		PrependName("Single Bar:", 0, 0).
		AppendPercentage(5, 0)

	for i := 0; i < 100; i++ {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		bar.Incr(1) // increment progress bar
	}

	// Don't forget to stop mpb's rendering goroutine
	p.Stop()
}

func ExampleBar_InProgress() {
	p := mpb.New()
	bar := p.AddBar(100).AppendPercentage(5, 0)

	for bar.InProgress() {
		time.Sleep(time.Millisecond * 20)
		bar.Incr(1)
	}
}

func ExampleBar_PrependFunc() {
	decor := func(s *mpb.Statistics, myWidth chan<- int, maxWidth <-chan int) string {
		str := fmt.Sprintf("%3d/%3d", s.Current, s.Total)
		// send width to Progress' goroutine
		myWidth <- utf8.RuneCountInString(str)
		// receive max width
		max := <-maxWidth
		return fmt.Sprintf(fmt.Sprintf("%%%ds", max+1), str)
	}

	totalItem := 100
	var wg sync.WaitGroup
	p := mpb.New()
	wg.Add(3) // add wg delta
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("Bar#%d:", i)
		bar := p.AddBar(int64(totalItem)).
			PrependName(name, len(name), 0).
			PrependFunc(decor)
		go func() {
			defer wg.Done()
			for i := 0; i < totalItem; i++ {
				time.Sleep(time.Duration(rand.Intn(totalItem)) * time.Millisecond)
				bar.Incr(1)
			}
		}()
	}
	wg.Wait() // Wait for goroutines to finish
	p.Stop()  // Stop mpb's rendering goroutine
}

func ExampleBar_AppendFunc() {
	decor := func(s *mpb.Statistics, myWidth chan<- int, maxWidth <-chan int) string {
		str := fmt.Sprintf("%3d/%3d", s.Current, s.Total)
		// send width to Progress' goroutine
		myWidth <- utf8.RuneCountInString(str)
		// receive max width
		max := <-maxWidth
		return fmt.Sprintf(fmt.Sprintf("%%%ds", max+1), str)
	}

	totalItem := 100
	var wg sync.WaitGroup
	p := mpb.New()
	wg.Add(3) // add wg delta
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("Bar#%d:", i)
		bar := p.AddBar(int64(totalItem)).
			PrependName(name, len(name), 0).
			AppendFunc(decor)
		go func() {
			defer wg.Done()
			for i := 0; i < totalItem; i++ {
				time.Sleep(time.Duration(rand.Intn(totalItem)) * time.Millisecond)
				bar.Incr(1)
			}
		}()
	}
	wg.Wait() // Wait for goroutines to finish
	p.Stop()  // Stop mpb's rendering goroutine
}
