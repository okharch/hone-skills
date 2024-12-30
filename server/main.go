package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// ExtractTextFromURL extracts text from a ChatGPT conversation page.
func ExtractTextFromURL(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
	)
	allocatorCtx, cancelAllocator := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAllocator()

	browserCtx, cancelBrowser := chromedp.NewContext(allocatorCtx)
	defer cancelBrowser()

	var extractedText string
	err := chromedp.Run(browserCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.Text(`body`, &extractedText),
	)
	if err != nil {
		return "", err
	}

	return extractedText, nil
}

// HandleRequest processes the API request.
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimSpace(r.URL.Query().Get("url"))
	if url == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	text, err := ExtractTextFromURL(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to extract text: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"text": text}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/extract", HandleRequest)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
