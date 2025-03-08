package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"calculator_golangV3/config/calculator"
	structs "calculator_golangV3/config/data"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func HandleComputation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var input structs.Request
		err = json.Unmarshal(data, &input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := rand.Int()

		go calculator.CalculateExpression(input.Expression, id)

		response, _ := json.Marshal(structs.ResponseOK{Id: id})
		fmt.Fprint(w, string(response))
		log.Println("POST", input, string(response), 201)

	} else {
		w.WriteHeader(405)
		log.Println(r.Method, 405)
	}
}

func HandleFetch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(404)
		response, _ := json.Marshal(map[string]structs.ResponseResult{"res": structs.ResponseResult{id, "err", 404}})
		fmt.Fprint(w, string(response))
		log.Println(string(response))
		return
	}

	file, err := os.Open("database/results.jsonl")
	if err != nil {
		w.WriteHeader(500)
		response, _ := json.Marshal(map[string]structs.ResponseResult{"res": structs.ResponseResult{id, "err", 500}})
		fmt.Fprint(w, string(response))
		log.Println(string(response))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var res structs.ResponseResult
		if err := json.Unmarshal(scanner.Bytes(), &res); err != nil {
			continue
		}
		if res.Id == id {
			w.WriteHeader(http.StatusOK)
			out, _ := json.Marshal(map[string]structs.ResponseResult{"res": res})
			fmt.Fprint(w, string(out))
			log.Println(string(out))
			return
		}
	}

	w.WriteHeader(404)
	out, _ := json.Marshal(map[string]structs.ResponseResult{"res": structs.ResponseResult{id, "not found", 404}})
	fmt.Fprint(w, string(out))
	log.Println(string(out))
}

func HandleFetchAll(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("database/results.jsonl")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "err")
		return
	}
	defer file.Close()

	var results []structs.ResponseResult
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var r structs.ResponseResult
		if err := json.Unmarshal(scanner.Bytes(), &r); err != nil {
			continue
		}
		results = append(results, r)
	}

	w.WriteHeader(http.StatusOK)
	out, _ := json.Marshal(map[string][]structs.ResponseResult{"res": results})
	fmt.Fprint(w, string(out))
	log.Println(string(out))
}

var ActiveTasks int = 0
var MaxTasks int = 1000
var taskLock sync.Mutex

func Initialize() {
	godotenv.Load(".env")
	val := os.Getenv("MAX_ROUTINES")
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("err")
			os.Exit(0)
		}
		MaxTasks = intVal
	} else {
		MaxTasks = 1000
	}
}

func HandleTaskOrchestration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		for {
			if ActiveTasks < MaxTasks {
				taskLock.Lock()
				ActiveTasks++
				taskLock.Unlock()
				data, _ := ioutil.ReadAll(r.Body)
				defer r.Body.Close()
				var input structs.AgentResponse
				json.Unmarshal(data, &input)
				timer := time.NewTimer(time.Duration(input.Operation_time) * time.Millisecond)
				result := 0.0
				if input.Operation == "+" {
					result = input.Arg1 + input.Arg2
				} else if input.Operation == "-" {
					result = input.Arg1 - input.Arg2
				} else if input.Operation == "*" {
					result = input.Arg1 * input.Arg2
				} else if input.Operation == "/" {
					result = input.Arg1 / input.Arg2
				}
				<-timer.C
				w.WriteHeader(http.StatusOK)
				out, _ := json.Marshal(structs.AgentResult{result})
				fmt.Fprint(w, string(out))
				taskLock.Lock()
				ActiveTasks--
				taskLock.Unlock()
				break
			}
		}
	}
}
