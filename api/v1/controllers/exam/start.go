package exam

import (
	"database/sql"
	"encoding/json"
	"github.com/Kedarnag13/Online_test/api/v1/models"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"github.com/gorilla/mux"
		"strconv"
)

type examController struct{}

var Exam examController

func (e examController) Questions(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	section_id, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", "password=password host=localhost dbname=online_test_dev sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	questions, err:= db.Exec("CREATE TABLE IF NOT EXISTS questions(id int, title text, option_1 varchar(100), option_2 varchar(100), option_3 varchar(100), option_4 varchar(100), answer varchar(100), section_id int, CONSTRAINT section_id_key FOREIGN KEY(section_id) REFERENCES sections (id), PRIMARY KEY(id))")
	if err != nil || questions == nil {
		log.Fatal(err)
	}
	get_questions, err := db.Query("SELECT id, title, option_1, option_2, option_3, option_4, answer FROM questions WHERE section_id=$1", section_id)
	if err != nil || get_questions == nil {
		log.Fatal(err)
	}
	defer get_questions.Close()

	questions_section := []models.Question{}

	for get_questions.Next() {
	var id int
	var title string
	var question_details models.Question
	var option_1 string
	var option_2 string
	var option_3 string
	var option_4 string
	var answer string
	err := get_questions.Scan(&id, &title, &option_1, &option_2, &option_3, &option_4, &answer)
	if err != nil {
			log.Fatal(err)
		}
		question_details = models.Question{id, title, option_1, option_2, option_3, option_4, answer}
		questions_section = append(questions_section, question_details)
	}
	b, err := json.Marshal(models.QuestionResponse{
			Success:     "true",
			Message:     "Questions per section",
			QuestionList:	questions_section,
			})
		if err != nil {
			log.Fatal(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
}
