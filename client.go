package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/data"
	"main/tools"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

func waitServer(url string, duration time.Duration) bool {
	deadline := time.Now().Add(duration)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)

		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return true

		}
	}
	return false
}

func main() {
	if !waitServer("http://localhost:9876", 5*time.Second) {
		fmt.Println("Server g ada")
		return
	}

	var choice int
	for {
		fmt.Println("Main Menu")
		fmt.Println("1. Get message")
		fmt.Println("2. Send file")
		fmt.Println("3. Quit")
		fmt.Print(">> ")
		fmt.Scanf("%d\n", &choice)
		if choice == 1 {
			getMessage()
		} else if choice == 2 {
			sendFile()
		} else if choice == 3 {
			break
		} else {
			fmt.Println("Invalid choice")
		}
	}
}

func getMessage() {
	resp, err := http.Get("http://localhost:9876")
	tools.ErrorHandler(err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	tools.ErrorHandler(err)

	fmt.Println("Server : ", string(data))
}

func sendFile() {
	var name string
	var age int

	// mirip scanner di java
	scanner := bufio.NewReader(os.Stdin)

	fmt.Printf("Input Name:")
	name, _ = scanner.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Input Age : ")
	fmt.Scanf("%d\n", &age)

	person := data.Person{Name: name, Age: age}
	// data person di encode JSON
	jsonData, err := json.Marshal(person)
	tools.ErrorHandler(err)

	//data penampung
	temp := new(bytes.Buffer)

	w := multipart.NewWriter(temp)

	personField, err := w.CreateFormField("Person")
	tools.ErrorHandler(err)

	_, err = personField.Write(jsonData)
	tools.ErrorHandler(err)

	//buka file txt
	file, err := os.Open("./file.txt")
	tools.ErrorHandler(err)

	defer file.Close()
	//buat field filenya
	fileField, err := w.CreateFormFile("File", file.Name())
	tools.ErrorHandler(err)

	// isi dari file akan dicopy ke fileField
	_, err = io.Copy(fileField, file)
	tools.ErrorHandler(err)

	// setelah kelar masukin data, multipartnya di close
	err = w.Close()
	tools.ErrorHandler(err)

	req, err := http.NewRequest("POST", "http://localhost:9876/sendFile", temp)
	tools.ErrorHandler(err)

	//set header utk ksh informasi kalau data yang dikirim itu merupakan multipart
	req.Header.Set("Content-Type", w.FormDataContentType())

	//kirim req ke server
	client := &http.Client{}
	resp, err := client.Do(req)
	tools.ErrorHandler(err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	tools.ErrorHandler(err)

	fmt.Println("Server : ", string(data))

}
