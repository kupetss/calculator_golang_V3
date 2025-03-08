package calculator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	structs "calculator_golangV3/config/data"

	"github.com/joho/godotenv"
)

type Operation struct {
	Action   string
	Index    int
	Priority int
}
type Result struct {
	Value string
	Index int
}
type Task struct {
	Operation Operation
	Channel   chan Result
}

var TimeAdd int = 0
var TimeSub int = 0
var TimeMul int = 0
var TimeDiv int = 0

func CleanString(s string) string {
	var result []string
	for _, char := range s {
		if char != ' ' {
			result = append(result, string(char))
		}
	}
	return strings.Join(result, "")
}

func ExecuteTask(t Task, elements []string) {
	x, _ := strconv.ParseFloat(elements[t.Operation.Index-1], 64)
	y, _ := strconv.ParseFloat(elements[t.Operation.Index+1], 64)
	var response structs.AgentResult
	url := "http://localhost:8080/internal/task"

	if t.Operation.Action == "+" {
		data, _ := json.Marshal(structs.AgentResponse{x, y, t.Operation.Action, TimeAdd})
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, &response)
	} else if t.Operation.Action == "-" {
		data, _ := json.Marshal(structs.AgentResponse{x, y, t.Operation.Action, TimeSub})
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, &response)
	} else if t.Operation.Action == "*" {
		data, _ := json.Marshal(structs.AgentResponse{x, y, t.Operation.Action, TimeMul})
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, &response)
	} else if t.Operation.Action == "/" {
		if y != 0 {
			data, _ := json.Marshal(structs.AgentResponse{x, y, t.Operation.Action, TimeDiv})
			resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(body, &response)
		} else {
			t.Channel <- Result{
				"err",
				t.Operation.Index,
			}
			return
		}
	}
	t.Channel <- Result{
		fmt.Sprintf("%f", response.Result),
		t.Operation.Index,
	}
}

func EvaluateExpression(expr string) (string, error) {
	expr = CleanString(expr)
	expr = strings.Replace(expr, "/", " / ", -1)
	expr = strings.Replace(expr, "*", " * ", -1)
	expr = strings.Replace(expr, "+", " + ", -1)
	expr = strings.Replace(expr, "-", " - ", -1)
	parts := strings.Split(expr, " ")
	if strings.Contains("+-/*", parts[0]) {
		return "", errors.New("invalid")
	}
	if strings.Contains("+-/*", parts[len(parts)-1]) {
		return "", errors.New("invalid")
	}
	if len(parts) == 1 {
		return expr, nil
	}
	numCount := 0
	opCount := 0
	priority := -1
	prev := -1
	opList := make([]Operation, 0)
	for i, element := range parts {
		if strings.Contains("+-/*", element) {
			if prev != -1 {
				if prev == 1 {
					fmt.Println(parts)
					return "", errors.New("invalid")
				}
			}
			opCount++
			prev = 1
			if priority == -1 && strings.Contains("/*", element) {
				priority = i
			}
			opList = append(opList, Operation{
				Action:   element,
				Index:    i,
				Priority: priority,
			})
		} else {
			for _, char := range element {
				if !strings.Contains("1234567890.", string(char)) {
					return "", errors.New("invalid")
				}
			}
			if prev != -1 {
				if prev == 0 {
					return "", errors.New("invalid")
				}
			}
			prev = 0
			numCount++
		}
	}
	if len(opList) > 1 {
		parallelOps := make([]Operation, 0)
		i := 0
		for {
			if i >= len(opList) {
				break
			}
			if i == 0 {
				if opList[i].Priority >= opList[i+1].Priority {
					parallelOps = append(parallelOps, opList[i])
					i += 2
				} else {
					i++
				}
			} else if i == len(opList)-1 {
				if opList[i].Priority >= opList[i-1].Priority {
					parallelOps = append(parallelOps, opList[i])
					i += 2
				} else {
					i++
				}
			} else {
				if opList[i-1].Priority <= opList[i].Priority && opList[i].Priority >= opList[i+1].Priority {
					parallelOps = append(parallelOps, opList[i])
					i += 2
				} else {
					i++
				}
			}
		}
		channel := make(chan Result)
		for _, op := range parallelOps {
			go ExecuteTask(Task{
				op,
				channel,
			}, parts)
		}
		for i := 0; i < len(parallelOps); i++ {
			select {
			case x, ok := <-channel:
				if ok {
					if x.Value == "err" {
						return "", errors.New("div by zero")
					}
					parts[x.Index] = x.Value
				}
			}
		}
		newParts := make([]string, 0)
		if strings.Contains("+-/*", parts[1]) {
			newParts = append(newParts, parts[0])
		}
		for i := 1; i+1 < len(parts); i++ {
			if !strings.Contains("+-/*", parts[i-1]) && !strings.Contains("+-/*", parts[i+1]) {
				newParts = append(newParts, parts[i])
			} else if strings.Contains("+-/*", parts[i-1]) && strings.Contains("+-/*", parts[i+1]) {
				newParts = append(newParts, parts[i])
			} else if !strings.Contains("+-/*", parts[i]) && !strings.Contains("+-/*", parts[i-1]) && !strings.Contains("+-/*", parts[i+1]) {
				newParts = append(newParts, parts[i])
			}
		}
		if strings.Contains("+-/*", parts[len(parts)-2]) {
			newParts = append(newParts, parts[len(parts)-1])
		}
		return EvaluateExpression(strings.Join(newParts, " "))
	}
	if numCount-opCount != 1 {
		return "", errors.New("invalid")
	}
	result := 0.0
	if priority != -1 {
		a, _ := strconv.ParseFloat(parts[priority-1], 64)
		b, _ := strconv.ParseFloat(parts[priority+1], 64)
		if parts[priority] == "*" {
			timer := time.NewTimer(time.Duration(TimeMul) * time.Millisecond)
			result = a * b
			<-timer.C
		} else {
			if b != 0 {
				timer := time.NewTimer(time.Duration(TimeDiv) * time.Millisecond)
				result = a / b
				<-timer.C
			} else {
				return "", errors.New("div by zero")
			}
		}
		if len(parts)-2 != 1 {
			return EvaluateExpression(fmt.Sprintf("%s%f%s", strings.Join(parts[:priority-1], ""), result, strings.Join(parts[priority+2:], "")))
		}
	} else {
		a, _ := strconv.ParseFloat(parts[0], 64)
		b, _ := strconv.ParseFloat(parts[2], 64)
		if parts[1] == "+" {
			timer := time.NewTimer(time.Duration(TimeAdd) * time.Millisecond)
			result = a + b
			<-timer.C
		} else {
			timer := time.NewTimer(time.Duration(TimeSub) * time.Millisecond)
			result = a - b
			<-timer.C
		}
		if len(parts)-2 != 1 {
			return EvaluateExpression(fmt.Sprintf("%f%s", result, strings.Join(parts[3:], "")))
		}
	}
	return fmt.Sprintf("%f", result), nil
}

func CalculateExpression(expr string, id int) (float64, error) {
	open := 0
	start := -1
	end := -1
	for i, char := range expr {
		if char == '(' {
			open++
			start = i
		} else if char == ')' {
			open--
			end = i
			if open == -1 {
				return 0, errors.New("invalid")
			}
			if end-start == 1 {
				return 0, errors.New("invalid")
			}
			res, err := EvaluateExpression(expr[start+1 : end])
			if err != nil {
				return 0, err
			}
			return CalculateExpression(expr[:start]+res+expr[end+1:], id)
		}
	}

	if open > 0 {
		return 0, errors.New("invalid")
	}
	out, err := EvaluateExpression(expr)
	if err != nil {
		return 0, err
	}
	out1, _ := strconv.ParseFloat(out, 64)

	file, err := os.OpenFile("database/results.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	newRes := structs.ResponseResult{
		Id:     id,
		Status: "ok",
		Result: out1,
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(newRes); err != nil {
		return 0, err
	}

	return out1, nil
}

func Initialize() {
	godotenv.Load(".env")
	val := os.Getenv("TIME_ADDITION_MS")
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("err")
			os.Exit(0)
		}
		TimeAdd = intVal
	} else {
		TimeAdd = 0
	}

	val = os.Getenv("TIME_SUBTRACTION_MS")
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("err")
			os.Exit(0)
		}
		TimeSub = intVal
	} else {
		TimeSub = 0
	}

	val = os.Getenv("TIME_MULTIPLICATIONS_MS")
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("err")
			os.Exit(0)
		}
		TimeMul = intVal
	} else {
		TimeMul = 0
	}

	val = os.Getenv("TIME_DIVISIONS_MS")
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("err")
			os.Exit(0)
		}
		TimeDiv = intVal
	} else {
		TimeDiv = 0
	}
	fmt.Printf("TimeAdd: %d\n", TimeAdd)
	fmt.Printf("TimeSub: %d\n", TimeSub)
	fmt.Printf("TimeMul: %d\n", TimeMul)
	fmt.Printf("TimeDiv: %d\n", TimeDiv)
}
