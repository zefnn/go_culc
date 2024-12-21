package application

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/zefnn/go_calc/pkg/calc"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

type RequestBody struct {
	Expression string `json:"expression"`
}

type ResponseBody struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

// Функция запуска приложения
// тут будем читать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение
func (a *Application) Run() error {
	for {
		// читаем выражение для вычисления из командной строки
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
		}
		// убираем пробелы, чтобы оставить только вычислемое выражение
		text = strings.TrimSpace(text)
		// выходим, если ввели команду "exit"
		if text == "exit" {
			log.Println("aplication was successfully closed")
			return nil
		}
		//вычисляем выражение
		result, err := calc.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed wit error: ", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

func handleCalc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ResponseBody{Error: "Метод не поддерживается!!!"})
		return

	}

	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseBody{Error: "Недопустимый запрос!!!"})
		return
	}
	////////////////////////////////////////////////////////////////////////
	//Логирование в файл
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "ошибка при открытии файла", http.StatusInternalServerError)
	}
	defer file.Close()
	log.SetOutput(file)

	log.Println("input expression: ", reqBody.Expression)

	///////////////////////////////////////////////////////////////////////
	result, err := calc.Calc(reqBody.Expression)
	if err != nil {
		switch err.Error() {
		case calc.ErrDivisionByZero.Error():
			w.WriteHeader(http.StatusUnprocessableEntity) //422
			json.NewEncoder(w).Encode(ResponseBody{Error: "Expression is not valid: division by zero"})
			log.Println("error: expression is not valid: division by zero")
		case calc.ErrInvalidNumberFormat.Error():
			w.WriteHeader(http.StatusUnprocessableEntity) //422
			json.NewEncoder(w).Encode(ResponseBody{Error: "Expression is not valid: invalid number format"})
			log.Println("error: expression is not valid: invalid number format")
		case calc.ErrUnsupportedSymbol.Error():
			w.WriteHeader(http.StatusInternalServerError) //500
			json.NewEncoder(w).Encode(ResponseBody{Error: "Internal server error: expression is not valid: unsupported symbol"})
			log.Println("error: internal server error: expression is not valid: unsupported symbol")
		case calc.ErrUnbalancedBrackets.Error():
			w.WriteHeader(http.StatusUnprocessableEntity) //422
			json.NewEncoder(w).Encode(ResponseBody{Error: "Expression is not valid: unbalanced brackets"})
			log.Println("error: expression is not valid: unbalanced brackets")
		case calc.ErrInvalidExpression.Error():
			w.WriteHeader(http.StatusUnprocessableEntity) //422
			json.NewEncoder(w).Encode(ResponseBody{Error: "Expression is not valid"})
			log.Println("error: expression is not valid")
		case calc.ErrNotEnoughOperands.Error():
			w.WriteHeader(http.StatusUnprocessableEntity) //422
			json.NewEncoder(w).Encode(ResponseBody{Error: "Expression is not valid: not enough operands"})
			log.Println("error: expression is not valid: not enough operands")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ResponseBody{Error: "Internal server error"})
			log.Println("error: internal server error")
		}
		return
	} else {
		log.Println(reqBody.Expression, "=", result)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResponseBody{Result: result})
	}

}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", handleCalc)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
