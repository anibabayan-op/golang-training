package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"Name"`
	Age  int    `json:"Age"`
}

func fromJSON(jsonStr string) (Person, error) {
	var p Person
	err := json.Unmarshal([]byte(jsonStr), &p)
	return p, err
}

func isAdult(p Person) bool {
	return p.Age >= 18
}

func main() {
	jsonStr := `{"Name": "Saleh", "Age": 20}`

	person, err := fromJSON(jsonStr)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
	if isAdult(person) {
		fmt.Println(person.Name, "is an adult.")
	} else {
		fmt.Println(person.Name, "is not an adult.")
	}
}
