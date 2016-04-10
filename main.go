package main

import (
	"github.com/Kedarnag13/Online_test/api/v1/controllers/account"
	"github.com/Kedarnag13/Online_test/api/v1/controllers/exam"
	"github.com/Kedarnag13/Online_test/api/v1/controllers/feedback"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sign_up", account.Registration.Create).Methods("POST")
	r.HandleFunc("/create_admin", account.Registration.CreateAdmin).Methods("POST")
	r.HandleFunc("/log_in", account.Session.Create).Methods("POST")
	r.HandleFunc("/delete/{auth_token:[A-Za-z0-9]+}", account.Session.Destroy).Methods("GET")
	r.HandleFunc("/section/evaluate", exam.Result.Create).Methods("POST")
	r.HandleFunc("/questions/{id:[0-9]+}", exam.Exam.Questions).Methods("GET")
	r.HandleFunc("/feedback", feedback.Feedback.Create).Methods("POST")
	r.HandleFunc("/export_csv", exam.Result.Export).Methods("GET")

	// HTTP Listening Port

	handler := cors.Default().Handler(r)
	http.Handle("/", handler)
	log.Println("main : Started : Listening on: http://localhost:3010 ...")
	log.Fatal(http.ListenAndServe("0.0.0.0:3010", nil))
}
