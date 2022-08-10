package task

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles(path.Join("views", "index.html"), path.Join("views", "layout.html"))
	if err != nil {
		log.Println(err)
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	// data := map[string]interface{}{
	// 	"title":   "golang Web",
	// 	"content": "Yusuf Adi Mirzaman",
	// }

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}
}

var db *sql.DB

var err error

// conect db and set template
func init() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/tugas_mandiri")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}

// deklarasi users variabel properti
type Task struct {
	ID       int
	Kegiatan string
	Tanggal  string
	Status   int
}

func (t Task) TaskStatus() string {
	var keterangan string
	if t.Status < 1 {
		keterangan = "Belum Selesai"
	} else {
		keterangan = "Selesai"
	}

	return keterangan
}

// list user
func Index(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles(path.Join("views", "task.html"), path.Join("views", "layout.html"))
	if err != nil {
		log.Println(err)
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	rows, e := db.Query(`SELECT * FROM task;`)

	if e != nil {
		log.Println(e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	tasks := make([]Task, 0)
	for rows.Next() {
		tsk := Task{}
		rows.Scan(&tsk.ID, &tsk.Kegiatan, &tsk.Tanggal, &tsk.Status)
		tasks = append(tasks, tsk)
	}
	err = tmpl.Execute(w, tasks)
	if err != nil {
		log.Println(err)
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}
}

// form create user
func Form(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(path.Join("views", "form.html"), path.Join("views", "layout.html"))
	if err != nil {
		log.Println(err)
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}
}

// action create users
func Proses(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		tsk := Task{}
		tsk.Kegiatan = req.FormValue("kegiatan")
		tsk.Tanggal = req.FormValue("tanggal")
		st, err := strconv.Atoi(req.FormValue("status"))
		if err != nil {
			log.Println(err)
			http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
			return
		}
		tsk.Status = st

		_, err = db.Exec("INSERT INTO task (kegiatan, tanggal, status) VALUES (?,?,?)",
			tsk.Kegiatan,
			tsk.Tanggal,
			tsk.Status,
		)

		if err != nil {
			log.Println(err)
			http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/task", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
}

// form edut users
func Ubah(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles(path.Join("views", "edit.html"), path.Join("views", "layout.html"))
	if err != nil {
		log.Println(err)
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}
	id := req.FormValue("id")
	rows, err := db.Query(`SELECT * FROM task WHERE id = ` + id + `;`)

	if err != nil {
		log.Println(err)
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	tsk := Task{}
	for rows.Next() {
		rows.Scan(&tsk.ID, &tsk.Kegiatan, &tsk.Tanggal, &tsk.Status)
	}
	err = tmpl.Execute(w, tsk)
	if err != nil {
		log.Println(err)
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}
}

// action edit users
func ProsesUbah(w http.ResponseWriter, req *http.Request) {
	_, err := db.Exec("UPDATE task SET kegiatan = ?, tanggal = ?, status = ? WHERE id = ?",
		req.FormValue("kegiatan"),
		req.FormValue("tanggal"),
		req.FormValue("status"),
		req.FormValue("id"),
	)

	if err != nil {
		log.Println(err)
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/task", http.StatusSeeOther)
}

// action deleted users
func Hapus(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	if id == "" {
		http.Error(w, "ID tidak ada", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM task WHERE id = ?", id)

	if err != nil {
		log.Println(err)
		http.Error(w, "terjadi kesalahan", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/task", http.StatusSeeOther)
}
