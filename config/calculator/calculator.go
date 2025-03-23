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

	"calculator_golangV3/config/structs"

	"github.com/joho/godotenv"
)

type Op struct {
	Act  string
	Pos  int
	Prio int
}
type Res struct {
	Val string
	Pos int
}
type Job struct {
	Op Op
	Ch chan Res
}

var T_ADD int = 0
var T_SUB int = 0
var T_MUL int = 0
var T_DIV int = 0

func Trim(s string) string {
	var out []string
	for _, c := range s {
		if c != ' ' {
			out = append(out, string(c))
		}
	}
	return strings.Join(out, "")
}
func Do(j Job, e []string) {
	x, _ := strconv.ParseFloat(e[j.Op.Pos-1], 64)
	y, _ := strconv.ParseFloat(e[j.Op.Pos+1], 64)
	var r structs.AgentResult
	url := "http://localhost:8080/internal/task"

	if j.Op.Act == "+" {
		d, _ := json.Marshal(structs.AgentResponse{x, y, j.Op.Act, T_ADD})
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(d))
		b, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(b, &r)
	} else if j.Op.Act == "-" {
		d, _ := json.Marshal(structs.AgentResponse{x, y, j.Op.Act, T_SUB})
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(d))
		b, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(b, &r)
	} else if j.Op.Act == "*" {
		d, _ := json.Marshal(structs.AgentResponse{x, y, j.Op.Act, T_MUL})
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(d))
		b, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(b, &r)
	} else if j.Op.Act == "/" {
		if y != 0 {
			d, _ := json.Marshal(structs.AgentResponse{x, y, j.Op.Act, T_DIV})
			resp, _ := http.Post(url, "application/json", bytes.NewBuffer(d))
			b, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(b, &r)
		} else {
			j.Ch <- Res{
				"err",
				j.Op.Pos,
			}
			return
		}
	}
	j.Ch <- Res{
		fmt.Sprintf("%f", r.Result),
		j.Op.Pos,
	}
}
func Eval(e string) (string, error) {
	e = Trim(e)
	e = strings.Replace(e, "/", " / ", -1)
	e = strings.Replace(e, "*", " * ", -1)
	e = strings.Replace(e, "+", " + ", -1)
	e = strings.Replace(e, "-", " - ", -1)
	p := strings.Split(e, " ")
	if strings.Contains("+-/*", p[0]) {
		return "", errors.New("invalid")
	}
	if strings.Contains("+-/*", p[len(p)-1]) {
		return "", errors.New("invalid")
	}
	if len(p) == 1 {
		return e, nil
	}
	n := 0
	o := 0
	prio := -1
	prev := -1
	oList := make([]Op, 0)
	for i, el := range p {
		if strings.Contains("+-/*", el) {
			if prev != -1 {
				if prev == 1 {
					fmt.Println(p)
					return "", errors.New("invalid")
				}
			}
			o++
			prev = 1
			if prio == -1 && strings.Contains("/*", el) {
				prio = i
			}
			oList = append(oList, Op{
				Act:  el,
				Pos:  i,
				Prio: prio,
			})
		} else {
			for _, c := range el {
				if !strings.Contains("1234567890.", string(c)) {
					return "", errors.New("invalid")
				}
			}
			if prev != -1 {
				if prev == 0 {
					return "", errors.New("invalid")
				}
			}
			prev = 0
			n++
		}
	}
	if len(oList) > 1 {
		parOps := make([]Op, 0)
		i := 0
		for {
			if i >= len(oList) {
				break
			}
			if i == 0 {
				if oList[i].Prio >= oList[i+1].Prio {
					parOps = append(parOps, oList[i])
					i += 2
				} else {
					i++
				}
			} else if i == len(oList)-1 {
				if oList[i].Prio >= oList[i-1].Prio {
					parOps = append(parOps, oList[i])
					i += 2
				} else {
					i++
				}
			} else {
				if oList[i-1].Prio <= oList[i].Prio && oList[i].Prio >= oList[i+1].Prio {
					parOps = append(parOps, oList[i])
					i += 2
				} else {
					i++
				}
			}
		}
		ch := make(chan Res)
		for _, op := range parOps {
			go Do(Job{
				op,
				ch,
			}, p)
		}
		for i := 0; i < len(parOps); i++ {
			select {
			case x, ok := <-ch:
				if ok {
					if x.Val == "err" {
						return "", errors.New("div by zero")
					}
					p[x.Pos] = x.Val
				}
			}
		}
		newP := make([]string, 0)
		if strings.Contains("+-/*", p[1]) {
			newP = append(newP, p[0])
		}
		for i := 1; i+1 < len(p); i++ {
			if !strings.Contains("+-/*", p[i-1]) && !strings.Contains("+-/*", p[i+1]) {
				newP = append(newP, p[i])
			} else if strings.Contains("+-/*", p[i-1]) && strings.Contains("+-/*", p[i+1]) {
				newP = append(newP, p[i])
			} else if !strings.Contains("+-/*", p[i]) && !strings.Contains("+-/*", p[i-1]) && !strings.Contains("+-/*", p[i+1]) {
				newP = append(newP, p[i])
			}
		}
		if strings.Contains("+-/*", p[len(p)-2]) {
			newP = append(newP, p[len(p)-1])
		}
		return Eval(strings.Join(newP, " "))
	}
	if n-o != 1 {
		return "", errors.New("invalid")
	}
	out := 0.0
	if prio != -1 {
		a, _ := strconv.ParseFloat(p[prio-1], 64)
		b, _ := strconv.ParseFloat(p[prio+1], 64)
		if p[prio] == "*" {
			timer := time.NewTimer(time.Duration(T_MUL) * time.Millisecond)
			out = a * b
			<-timer.C
		} else {
			if b != 0 {
				timer := time.NewTimer(time.Duration(T_DIV) * time.Millisecond)
				out = a / b
				<-timer.C
			} else {
				return "", errors.New("div by zero")
			}
		}
		if len(p)-2 != 1 {
			return Eval(fmt.Sprintf("%s%f%s", strings.Join(p[:prio-1], ""), out, strings.Join(p[prio+2:], "")))
		}
	} else {
		a, _ := strconv.ParseFloat(p[0], 64)
		b, _ := strconv.ParseFloat(p[2], 64)
		if p[1] == "+" {
			timer := time.NewTimer(time.Duration(T_ADD) * time.Millisecond)
			out = a + b
			<-timer.C
		} else {
			timer := time.NewTimer(time.Duration(T_SUB) * time.Millisecond)
			out = a - b
			<-timer.C
		}
		if len(p)-2 != 1 {
			return Eval(fmt.Sprintf("%f%s", out, strings.Join(p[3:], "")))
		}
	}
	return fmt.Sprintf("%f", out), nil
}

func Calc(e string, id string) (float64, error) {
	open := 0
	start := -1
	end := -1
	for i, c := range e {
		if c == '(' {
			open++
			start = i
		} else if c == ')' {
			open--
			end = i
			if open == -1 {
				return 0, errors.New("invalid")
			}
			if end-start == 1 {
				return 0, errors.New("invalid")
			}
			res, err := Eval(e[start+1 : end])
			if err != nil {
				return 0, err
			}
			return Calc(e[:start]+res+e[end+1:], id)
		}
	}

	if open > 0 {
		return 0, errors.New("invalid")
	}
	out, err := Eval(e)
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
		Id:         id,
		Status:     "ok",
		Expression: e,
		Result:     out1,
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(newRes); err != nil {
		return 0, err
	}

	return out1, nil
}

func Init() {
	godotenv.Load(".env")
	val := os.Getenv("TIME_ADDITION_MS")
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("err")
			os.Exit(0)
		}
		T_ADD = intVal
	} else {
		T_ADD = 0
	}

	val = os.Getenv("TIME_SUBTRACTION_MS")
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("err")
			os.Exit(0)
		}
		T_SUB = intVal
	} else {
		T_SUB = 0
	}

	val = os.Getenv("TIME_MULTIPLICATIONS_MS")
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("err")
			os.Exit(0)
		}
		T_MUL = intVal
	} else {
		T_MUL = 0
	}

	val = os.Getenv("TIME_DIVISIONS_MS")
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("err")
			os.Exit(0)
		}
		T_DIV = intVal
	} else {
		T_DIV = 0
	}
	fmt.Printf("T_ADD: %d\n", T_ADD)
	fmt.Printf("T_SUB: %d\n", T_SUB)
	fmt.Printf("T_MUL: %d\n", T_MUL)
	fmt.Printf("T_DIV: %d\n", T_DIV)

}
