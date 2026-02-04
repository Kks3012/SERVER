package main
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)
const logFile = "network_log.txt"
func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		content, err := ioutil.ReadFile(logFile)
		if err != nil {
			http.Error(w, "Лог-файл порожній або не знайдений", http.StatusNotFound)
			return
		}
		fmt.Println("--- Поточний лог пристроїв ---")
		fmt.Println(string(content))
		fmt.Fprintf(w, "Лог успішно виведено в термінал сервера")
	case http.MethodPost:
		body, _ := ioutil.ReadAll(r.Body)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		data := fmt.Sprintf("[%s] Дані пристрою: %s\n", timestamp, string(body))
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.WriteString(data)
		fmt.Fprintf(w, "Запис додано: %s", data)
	case http.MethodDelete:
		err := os.Remove(logFile)
		if err != nil {
			http.Error(w, "Не вдалося видалити файл", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Лог-файл очищено")
	default:
		http.Error(w, "Метод не підтримується", http.StatusMethodNotAllowed)
	}
}
func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Сервер запущено на http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}