package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)



func main() {
	// Create a new HTTP client
	client := &http.Client{}
	
	// Create a new request
	req, err := http.NewRequest("GET", "https://my.api.mockaroo.com/users.json", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	
	// Add the API key header
	req.Header.Add("X-API-Key", "15d6b050")
	
	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	
	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	
	// Parse JSON into a single Person object
	var person Person
	err = json.Unmarshal(body, &person)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	
	// Print the result
	fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
}
