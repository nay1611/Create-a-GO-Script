package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	// Launch a new browser instance
	url, err := launcher.New().Headless(false).Launch()
	if err != nil {
		log.Fatalf("Failed to launch browser: %v", err)
	}
	fmt.Println("Browser launched at URL:", url)

	// Connect to the browser
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// Create a new page
	page := browser.MustPage("https://www.google.com/")
	fmt.Println("Navigated to Google")

	// Wait for the page to load
	page.MustWaitLoad()
	fmt.Println("Page fully loaded")

	// Wait for the search box to appear
	err = rod.Try(func() {
		searchBox := page.Timeout(120 * time.Second).MustElement(`#APjFqb`)
		log.Println("Search box element found:", searchBox)
		searchBox.MustInput("what is go-rod")
		searchBox.MustType(input.Enter)
	})
	if err != nil {
		log.Fatalf("Failed to enter search query: %v", err)
	}

	// Wait for the results to load
	err = rod.Try(func() {
		page.Timeout(60 * time.Second).MustElement(`#search`)
	})
	if err != nil {
		log.Fatalf("Failed to load search results: %v", err)
	}
	fmt.Println("Search results loaded")

	// Fetch titles of all search results
	elements, err := page.Elements("h3")
	if err != nil {
		log.Fatalf("Failed to fetch search result elements: %v", err)
	}
	for i, el := range elements {
		title, err := el.Text()
		if err != nil {
			log.Printf("Failed to fetch text for element %d: %v", i+1, err)
			continue
		}
		fmt.Printf("Result %d: %s\n", i+1, title)
	}

	// Take a screenshot of the page
	screenshotData := page.MustScreenshot()
	if err := os.WriteFile("screenshot.png", screenshotData, 0644); err != nil {
		log.Fatalf("Failed to save screenshot: %v", err)
	} else {
		fmt.Println("Screenshot saved as screenshot.png")
	}
}
