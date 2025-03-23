package structs

type Request struct {
	Expression string `json:"expression"`
}

type ResponseOK struct {
	Id string `json:"id"`
}

type ResponseERROR struct {
	Error string `json:"error"`
}

type ResponseResult struct {
	Id         string  `json:"id"`
	Status     string  `json:"status"`
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

type ResponseExpression struct {
	Expression ResponseResult `json:"expression"`
}

type AgentResponse struct {
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time int     `json:"operation_time"`
}

type AgentResult struct {
	Result float64 `json:"result"`
}
