package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func NewPersonFromJSON(jsonStr string) (Person, error) {
	var p Person
	err := json.Unmarshal([]byte(jsonStr), &p)
	return p, err
}

func (p Person) IsAdult() bool {
	return p.Age >= 18
}

func main() {
	jsonStr := `{"Name": "Saleh", "Age": 20}`

	person, err := NewPersonFromJSON(jsonStr)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)

	if person.IsAdult() {
		fmt.Println(person.Name, "is an adult.")
	} else {
		fmt.Println(person.Name, "is not an adult.")
	}
}
