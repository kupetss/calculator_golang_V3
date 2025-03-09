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
	"calculator_golangV3/config/structs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func HandleCompute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var in structs.Request
		err = json.Unmarshal(d, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := rand.Int()

		go calculator.Calc(in.Expression, id)

		res, _ := json.Marshal(structs.ResponseOK{Id: id})
		fmt.Fprint(w, string(res))
		log.Println("POST", in, string(res), 201)

	} else {
		w.WriteHeader(405)
		log.Println(r.Method, 405)
	}
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	idStr := p["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(404)
		res, _ := json.Marshal(map[string]structs.ResponseResult{"res": structs.ResponseResult{id, "err", 404}})
		fmt.Fprint(w, string(res))
		log.Println(string(res))
		return
	}

	file, err := os.Open("database/results.jsonl")
	if err != nil {
		w.WriteHeader(500)
		res, _ := json.Marshal(map[string]structs.ResponseResult{"res": structs.ResponseResult{id, "err", 500}})
		fmt.Fprint(w, string(res))
		log.Println(string(res))
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

func HandleList(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("database/results.jsonl")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "err")
		return
	}
	defer file.Close()

	var res []structs.ResponseResult
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var r structs.ResponseResult
		if err := json.Unmarshal(scanner.Bytes(), &r); err != nil {
			continue
		}
		res = append(res, r)
	}

	w.WriteHeader(http.StatusOK)
	out, _ := json.Marshal(map[string][]structs.ResponseResult{"res": res})
	fmt.Fprint(w, string(out))
	log.Println(string(out))
}

var Active int = 0
var Max int = 1000
var lock sync.Mutex

func Init() {
	godotenv.Load(".env")
	val := os.Getenv("MAX_ROUTINES")
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("err")
			os.Exit(0)
		}
		Max = intVal
	} else {
		Max = 1000
	}
}

func HandleOrchestrate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		for {
			if Active < Max {
				lock.Lock()
				Active++
				lock.Unlock()
				d, _ := ioutil.ReadAll(r.Body)
				defer r.Body.Close()
				var in structs.AgentResponse
				json.Unmarshal(d, &in)
				timer := time.NewTimer(time.Duration(in.Operation_time) * time.Millisecond)
				res := 0.0
				if in.Operation == "+" {
					res = in.Arg1 + in.Arg2
				} else if in.Operation == "-" {
					res = in.Arg1 - in.Arg2
				} else if in.Operation == "*" {
					res = in.Arg1 * in.Arg2
				} else if in.Operation == "/" {
					res = in.Arg1 / in.Arg2
				}
				<-timer.C
				w.WriteHeader(http.StatusOK)
				out, _ := json.Marshal(structs.AgentResult{res})
				fmt.Fprint(w, string(out))
				lock.Lock()
				Active--
				lock.Unlock()
				break
			}
		}
	}
}
