package main

import (
	"fmt"
	"flag"
	"net/http"
	"encoding/json"
	"strconv"
	"bytes"
)

type Todo struct {
	Id int `json:"id"`
	Text string `json:"text"`
	Status bool `json:"status"`
}

const url = "http://localhost:3000/todos"

// todo -l
// todo -i {id} -g
// todo -a {text}
// todo -i {id} -e {text}
// todo -i {id} -d
func main() {

	listFlag := flag.Bool("l", false, "todo list")
	idFlag := flag.Int("i", 0, "todo id")
	getFlag := flag.Bool("g", false, "Get todo id")
	addFlag := flag.String("a", "", "Add todo text")
	editFlag := flag.String("e", "", "Edit todo text")
	deleteFlag := flag.Bool("d", false, "Delete todo")
	flag.Parse()

	executed := false

	if !executed && *listFlag {
		Print()
		executed = true
	}

	if !executed && *idFlag > 0 && *getFlag {
		_, err := http.Get(url + "/" + strconv.Itoa(*idFlag))
		if err != nil {
			fmt.Print(err)
		} else {
			Print()
		}
		executed = true
	}

	if !executed && *addFlag != "" {
		todo := Todo{Text:*addFlag}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(todo)
		_, err := http.Post(url + "/", "application/json; charset=utf-8", buf)
		if err != nil {
			fmt.Print(err)
		} else {
			Print()
		}
		executed = true
	}

	if !executed && *idFlag > 0 && *editFlag != "" {
		todo := Todo{Id:*idFlag, Text:*editFlag, Status:false}
		fmt.Println(url + "/" + strconv.Itoa(*idFlag))
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(todo)
		req, _ := http.NewRequest("PATCH", url + "/" + strconv.Itoa(*idFlag), buf)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		client := &http.Client{}
		_, err := client.Do(req)
		if err != nil {
			fmt.Print(err)
		} else {
			Print()
		}
		executed = true
	}

	if !executed && *idFlag > 0 && *deleteFlag {
		buf := new(bytes.Buffer)
		req, _ := http.NewRequest("DELETE", url + "/" + strconv.Itoa(*idFlag), buf)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		client := &http.Client{}
		_, err := client.Do(req)
		if err != nil {
			fmt.Print(err)
		} else {
			Print()
		}

		executed = true
	}

	if !executed {
		flag.Usage()
	}
}

func Print() {
	res, err := http.Get(url + "/")
	if err != nil {
		fmt.Print(err)
	} else {
		todos := []Todo{}
		json.NewDecoder(res.Body).Decode(&todos)
		for _, todo := range todos {
			ck := "o"
			if todo.Status {
				ck = "x"
			}
			fmt.Println("[" +ck + "]id:" + strconv.Itoa(todo.Id) + " " + todo.Text)
		}
	}
}
