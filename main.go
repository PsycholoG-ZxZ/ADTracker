package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"database/sql"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Types
type data struct {
	Email string `json:"Email"`
	URL   string `json:"URL"`
}

type allData []data

var someInf = allData{
	{
		Email: "dimka_volnyi@mail.ru",
		URL:   "google.com",
	},
}

type item struct {
	id    string
	URL   string
	price string
}

var database *sql.DB

var ch = 1

//MainHandler ...
func MainHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello")
}

func GetVer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	sql_string := "update adtracker.users SET verification = 1 where verification_code = " + id + " ;"

	_, err := database.Exec(sql_string)
	if err != nil {
		log.Println(err)
		io.WriteString(w, "verification is not done! Try again")
	} else {
		io.WriteString(w, "verification done! Welcome")
	}

}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newData data
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Task Data")
	}

	json.Unmarshal(reqBody, &newData)
	someInf = append(someInf, newData)

	user_string := "select user_id from adtracker.users where email = '" + someInf[ch].Email + "'"
	row := database.QueryRow(user_string)

	user_id_sql := ""
	min := 1000000000
	max := 9999999999

	code := rand.Intn(max-min) + min
	fmt.Println("code")
	fmt.Println(code)

	err = row.Scan(&user_id_sql)
	if err != nil {
		log.Println(err)
		// Создаем пользователя
		user_string_ins := "insert into adtracker.users (email, verification_code) values ('" + someInf[ch].Email + "' , '" + strconv.Itoa(code) + "')"
		_, err = database.Exec(user_string_ins)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("New member add : " + someInf[ch].Email)
			row := database.QueryRow(user_string)
			err2 := row.Scan(&user_id_sql)
			if err2 != nil {
				log.Println(err)
			}
		} //TODO :: отправить письмо верификации

	} else {
		//мейл есть
	}

	user_string = "select verification, verification_code from adtracker.users where email = '" + someInf[ch].Email + "'"
	row = database.QueryRow(user_string)
	boolSql := ""
	codeSql := ""
	err = row.Scan(&boolSql, &codeSql)
	if err != nil {
		log.Println(err)
	}
	log.Println("BOOOL : " + boolSql)

	if boolSql == "0" {

		bodyString := "If you want see notifications, folow this link: localhost:8282/" + codeSql

		mainMail(someInf[ch].Email, bodyString)
	}

	url_string := "select item_id from adtracker.items where url = '" + someInf[ch].URL + "'"
	row = database.QueryRow(url_string)
	item_id_sql := ""

	err = row.Scan(&item_id_sql)
	if err != nil {
		log.Println(err)
		// Создаем ссылку
		url_string_ins := "insert into adtracker.items (URL, price) values ('" + someInf[ch].URL + "' , '" + (getPrice(someInf[ch].URL)) + "')"
		_, err = database.Exec(url_string_ins)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("New URL add : " + someInf[ch].URL)
			row = database.QueryRow(url_string)
			err2 := row.Scan(&item_id_sql)
			if err2 != nil {
				log.Println(err)
			}
		}

	} else {
		//ссылка есть
	}

	inc_string := "insert into adtracker.users_items (user_id, item_id) values ('" + user_id_sql + "' , '" + item_id_sql + "')"
	_, err = database.Exec(inc_string)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("New pair add : " + user_id_sql + " " + item_id_sql)
	}

	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusCreated)
	io.WriteString(w, someInf[ch].Email+"  "+someInf[ch].URL)

	fmt.Println(someInf[ch].Email)

	ch++

}

func getPrice(URL string) string {
	response, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(response)
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("Error loading HTTP response body. ", err)
	}

	res := ""
	// Find and print image URLs
	document.Find(".js-item-price").Each(func(index int, element *goquery.Selection) {
		imgSrc := element.Text()
		res = imgSrc
	})
	fmt.Println(res)
	return res
}

func f() {
	for {
		fmt.Println("f on work")

		sql_string := "Select * from items"

		rows, err := database.Query(sql_string)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()

		//items := []item{}

		for rows.Next() {
			it := item{}
			urlString := ""
			err := rows.Scan(&it.id, &it.price, &it.URL)
			urlString = it.URL
			if err != nil {
				fmt.Println(err)
				continue
			}
			price := getPrice(urlString)

			sql_string := "Select price from items where URL = '" + urlString + "'"

			row := database.QueryRow(sql_string)
			sql_price := ""

			err = row.Scan(&sql_price)
			if err != nil {
				log.Println(err)
			}

			if sql_price != price {
				//TODO отправляем email про изменение
				log.Println("Email Sender:")

				bodyString := "New price for this item: " + urlString + " is: " + price

				//Select email from users join users_items ON (users_items.user_id = users.user_id) AND (users_items.item_id = 6)
				id_string := "Select email from users join users_items ON (users_items.user_id = users.user_id) AND (users_items.item_id = " + it.id + ")"

				rows_id, err := database.Query(id_string)
				if err != nil {
					log.Println(err)
				}
				defer rows_id.Close()

				for rows_id.Next() {
					emailString := ""
					err := rows_id.Scan(&emailString)
					if err != nil {
						fmt.Println(err)
						continue
					}

					check_str := "select verification from adtracker.users where email = '" + emailString + "'"
					check := database.QueryRow(check_str)
					err = check.Scan(&check_str)
					if err != nil {
						fmt.Println(err)
					}
					if check_str == "1" {
						mainMail(emailString, bodyString)
					}

					log.Println(" Email Send for " + emailString)
				}

				//log.Println("fake Email Send")

				url_string_ins := "update adtracker.items SET price = '" + price + "' Where URL = '" + urlString + "'"

				_, err = database.Exec(url_string_ins)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("price updated")
				}

			} else {
				log.Println("price still the same")
			}

		}

		time.Sleep(10 * time.Second)

	}
}

func main() {

	db, err := sql.Open("mysql", "root:qweasd123@/adtracker")

	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	//https://www.avito.ru/moskva/audio_i_video/televizor_samsung_40_dyuymov_2026758687

	// Make HTTP request

	router := mux.NewRouter()

	router.HandleFunc("/", MainHandler)
	router.HandleFunc("/create", createTask).Methods("POST")
	router.HandleFunc("/{id:[0-9]+}", GetVer)

	http.Handle("/", router)

	fmt.Println("Server is listening...")
	go f()
	fmt.Println(http.ListenAndServe(":8282", nil))

}
