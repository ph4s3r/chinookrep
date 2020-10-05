package main

import (
	"fmt"
	"github.com/gorilla/mux" //REST API lib
	"net/http"
	"regexp"
)

//WEBSERVER FUNCTIONS BEGIN

func doReportHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	vals := r.URL.Query()
	fromDate, fromDate_ok := vals["fromDate"]
	toDate, toDate_ok := vals["toDate"]
	re := regexp.MustCompile("((19|20)\\d\\d)(0?[1-9]|1[012])(0?[1-9]|[12][0-9]|3[01])")
	ok := fromDate_ok && re.MatchString(fromDate[0]) && toDate_ok && re.MatchString(toDate[0])

	if ok {
		fmt.Fprintf(w, "doReport handler sent a task to celery with params fromDate: %v toDate: %v\n", fromDate[0], toDate[0])
		go celeryClient("doreport_a", fromDate[0], toDate[0])
		go celeryClient("doreport_b", fromDate[0], toDate[0])
		go celeryClient("doreport_c", fromDate[0], toDate[0])
	} else {
		fmt.Fprintf(w, "please provide the fromDate and toDate as a query param in the following format: http://<server>/doReport?fromDate=YYYYMMDD&toDate=YYYYMMDD\r\n")
	}

}

func listReportsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "listReportHandler called")
}

//WEBSERVER FUNCTIONS END

func main() {
	//doing the first few sql preparative tasks
	sqlprep()
	//creating 3 workers
	go createWorker(3)
	fmt.Println("celery worker started, task has been registered")
	//starting 2 handlers of webserver
	r := mux.NewRouter()
	r.HandleFunc("/doReport", doReportHandler)
	r.HandleFunc("/listReports", listReportsHandler)

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(501)
		w.Write([]byte(`available endpoints: /doReport and /listReports`))
	})

	fmt.Println("gmux starting to serve requests on port 8000")
	err := http.ListenAndServe(":8000", r)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//this is used to continue running if there is an error
func moveonErr(err error) {
	if err == nil {
		return
	}
}
