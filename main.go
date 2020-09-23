package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pbkheiron/sweinterview/calc"
)

var (
	interactive = flag.Bool("i", true, "run in interactive mode")
	port        = flag.Int("p", 0, "http port number; Forces HTTP server mode")
	notation    = flag.String("n", "infix", "notation (infix | prefix)")
)

func main() {
	parseFlags()
	if *interactive {
		go runInteractiveMode(*notation)
		waitForSigint()
	} else {
		runHttpMode(*port)
	}
}

func parseFlags() {
	flag.Parse()
	if *port != 0 {
		*interactive = false
	}
}

func exitOnError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func waitForSigint() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	<-c
	fmt.Fprintf(os.Stderr, "SIGINT. Exit now.\n")
}

func runInteractiveMode(notation string) {
	evalFunc, err := chooseEvalFunc(notation)
	exitOnError(err)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "> ")
		expr, _ := reader.ReadString('\n')
		result, err := evalFunc(expr)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
		} else {
			fmt.Fprintln(os.Stdout, result)
		}
	}
}

type evalFunc func(s string) (float64, error)

func chooseEvalFunc(notation string) (evalFunc, error) {
	switch notation {
	case "infix":
		return calc.InfixEval, nil
	case "prefix":
		return calc.PrefixEval, nil
	default:
		return nil, fmt.Errorf("unsupported notation: %s", notation)
	}
}

func runHttpMode(port int) {
	http.Handle("/calc", http.HandlerFunc(calcHandlerFunc))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	exitOnError(err)
}

type CalcReq struct {
	Notation string `json:"notation"`
	Expr     string `json:"expr"`
}

type CalcResp struct {
	Result float64 `json:"result"`
	Err    string  `json:"err,omitempty"`
}

func calcHandlerFunc(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	sendJsonResp := func(status int, obj interface{}) {
		b, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"err"": "server marshal error"}`))
			return
		}
		w.WriteHeader(status)
		w.Write(b)
	}

	calcReq := CalcReq{Notation: "infix"}
	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(&calcReq)
	if err != nil {
		resp := CalcResp{Err: err.Error()}
		sendJsonResp(http.StatusBadRequest, resp)
		return
	}

	evalFun, err := chooseEvalFunc(calcReq.Notation)
	if err != nil {
		resp := CalcResp{Err: err.Error()}
		sendJsonResp(http.StatusBadRequest, resp)
		return
	}

	result, err := evalFun(calcReq.Expr)
	if err != nil {
		resp := CalcResp{Err: err.Error()}
		sendJsonResp(http.StatusBadRequest, resp)
	} else {
		resp := CalcResp{Result: result}
		sendJsonResp(http.StatusOK, resp)
	}
}
