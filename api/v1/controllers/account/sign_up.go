package account

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"
	"crypto/md5"
	"io"
	"github.com/Kedarnag13/Online_test/api/v1/models"
	"github.com/asaskevich/govalidator"
	"github.com/Kedarnag13/Online_test/api/v1/controllers"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"regexp"
	"log"
	"encoding/hex"
)

type registrationController struct{}

var Registration registrationController

func (r registrationController) Create(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	flag := 1
	var u models.Register

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", "password=password host=localhost dbname=online_test_dev sslmode=disable")
	if err != nil {
		panic(err)
	}

	res, err := db.Query("SELECT email, phone_number FROM users ")
	if err != nil {
		panic(err)
	}
	defer res.Close()

	fetch_id, err := db.Query("SELECT coalesce(max(id), 0) FROM users")
	if err != nil {
		panic(err)
	}
	defer fetch_id.Close()

	if flag == 1 {
		email := `^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,4}$`
		exp, err := regexp.Compile(email)
		if err != nil {
			panic(err)
		}

		if !exp.MatchString(u.Email){
			flag = 0
			_, err := govalidator.ValidateStruct(u)
			if err != nil {
				println("error: " + err.Error())
			}

			b, err := json.Marshal(models.EmailErrorMessage{
				Success: "false",
				Email_error: "Give a valid email",
			})
			if err != nil {
				panic(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto create_user_end
		}
		if u.Password != u.Password_confirmation {
			flag = 0
			b, err := json.Marshal(models.PasswordErrorMessage{
				Success: "false",
				Password_error:   "Password and confirm password do not match!",
			})
			if err != nil {
				panic(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto create_user_end
		}
	}

	if flag == 1 {
		for res.Next() { // email already exist condition
			var email string
			var phone_number string
			err = res.Scan(&email, &phone_number)
			if err != nil {
				panic(err)
			}

			if email == u.Email && phone_number == u.Phone_number {
				b, err := json.Marshal(models.EmailPasswordErrorMessage{
					Success: "false",
					Email_error: "Email already exist",
					Phone_number_error: "Phone number already exist",
				})
				if err != nil {
					panic(err)
				}
				rw.Header().Set("Content-Type", "application/json")
				rw.Write(b)
				flag = 0
				goto create_user_end
				}else if email == u.Email {
					b, err := json.Marshal(models.EmailErrorMessage{
						Success: "false",
						Email_error:   "Email id already exist",
					})
					if err != nil {
						panic(err)
					}
					rw.Header().Set("Content-Type", "application/json")
					rw.Write(b)
					flag = 0
					goto create_user_end
					}else if phone_number == u.Phone_number {
						b, err := json.Marshal(models.PhoneNumberErrorMessage{
							Success: "false",
							Phone_number_error: "Phone number already exist",
						})
						if err != nil {
							panic(err)
						}
						rw.Header().Set("Content-Type", "application/json")
						rw.Write(b)
						flag = 0
						goto create_user_end
					}
				}

				// Insert into users table ======================================

				for fetch_id.Next() {
					var id int
					err = fetch_id.Scan(&id)

					if err != nil {
						panic(err)
					}
					id = id + 1
					created_at := time.Now()
					var sStmt string = "insert into users (id, first_name, last_name, email, college, branch, phone_number, year_of_passing, city, password, password_confirmation, created_at, batch, role) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)"
					db, err := sql.Open("postgres", "password=password host=localhost dbname=online_test_dev sslmode=disable")
					if err != nil {
						panic(err)
					}
					stmt, err := db.Prepare(sStmt)
					if err != nil {
						panic(err)
					}
					defer stmt.Close()


					key := []byte("traveling is fun")
					password := []byte(u.Password)
					confirm_password := []byte(u.Password_confirmation)
					encrypt_password := controllers.Encrypt(key, password)
					encrypt_password_confirmation := controllers.Encrypt(key, confirm_password)

					user_res, err := stmt.Exec(id, u.First_name, u.Last_name, u.Email, u.College, u.Branch, u.Phone_number, u.Year_of_passing, u.City, encrypt_password, encrypt_password_confirmation, created_at, u.Batch,"user")
					if err != nil || user_res == nil {
						panic(err)
					}
					defer stmt.Close()

					// Create Session for the User =========================================

					auth_string := strconv.FormatInt(time.Now().Unix(), 10)
					h := md5.New()
					io.WriteString(h, auth_string)
					auth_token := hex.EncodeToString(h.Sum(nil))
					var session string = "insert into sessions (start_time, user_id, auth_token) values ($1,$2,$3)"
					prepare_session, err := db.Prepare(session)
					if err != nil {
						panic(err)
					}
					defer prepare_session.Close()
					start_time := time.Now()
					session_res, err := prepare_session.Exec(start_time, id, string(auth_token))
					if err != nil || session_res == nil {
						panic(err)
					}
					db.Close()
					user := models.Register{id, u.First_name, u.Last_name, u.Email, u.Password, u.Password_confirmation, u.College, u.Branch, u.Year_of_passing, u.City, u.Phone_number, u.Batch}

					b, err := json.Marshal(models.SignUp{
						Success: "true",
						Message: "User created Successfully!",
						User:    user,
						Session: models.Session{id, start_time, string(auth_token)},
					})

					if err != nil || res == nil {
						panic(err)
					}

					rw.Header().Set("Content-Type", "application/json")
					rw.Write(b)
				}
				// defer fetch_id.Close()
			}
			create_user_end:
			db.Close()
		}

		func (r registrationController) CreateAdmin(rw http.ResponseWriter, req *http.Request) {

			db, err := sql.Open("postgres", "password=password host=localhost dbname=online_test_dev sslmode=disable")
			if err != nil {
				panic(err)
			}
			flag := 1
			CheckAdmin, err := db.Query("SELECT phone_number from USERS where phone_number = $1 AND role = $2","1111111111","admin")
			if err != nil {
				panic(err)
			}
			defer CheckAdmin.Close()
			for CheckAdmin.Next(){
				flag = 0
				b, err := json.Marshal(models.ErrorMessage{
					Success: "false",
					Error:   "Admin already exist",
				})
				if err != nil {
					panic(err)
				}
				rw.Header().Set("Content-Type", "application/json")
				rw.Write(b)

				goto AdminEnd
			}
			if flag == 1{

				fetch_id, err := db.Query("SELECT coalesce(max(id), 0) FROM users")
				if err != nil {
					panic(err)
				}
				defer fetch_id.Close()
				var user_id int
				for fetch_id.Next(){
					err = fetch_id.Scan(&user_id)
					if err != nil {
						panic(err)
					}
					user_id = user_id + 1
					key := []byte("traveling is fun")
					password := []byte("Qwinix123")
					encrypt_password := controllers.Encrypt(key, password)
					log.Printf("user_id:%v",user_id)
					_, err = db.Query("INSERT into USERS(id, first_name, last_name, email, password, phone_number, role) VALUES($1, 'Qwinix', 'admin', 'qwinixadmin@yopmail.com', $2, '1111111111', 'admin')", user_id, encrypt_password)
					if err != nil {
						panic(err)
					}
					b, err := json.Marshal(models.AdminSuccessMessage{
						Success: "True",
						Message:   "Admin Created Successfully",
					})
					if err != nil {
						panic(err)
					}
					rw.Header().Set("Content-Type", "application/json")
					rw.Write(b)
				}
			}
			AdminEnd:
			db.Close()
		}
