// main.go
package main

import "os"
import "io"
import "net/http"
import "html/template"
import "log"

func main() {
	conf, err := NewConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		log.Print(r.URL.Path)
		if r.URL.Path[1:] != "" {
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, conf.NotFoundDoc)
			return
		}
		t, _ := template.ParseFiles(conf.UploadDoc)
		t.Execute(w, nil)
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			t, _ := template.ParseFiles(conf.UploadDoc)
			t.Execute(w, nil)
			return
		}
		r.ParseMultipartForm(32 << 24)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Printf("(1) %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()
		
		f, err := os.OpenFile(conf.UploadPath + handler.Filename, os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil {
			log.Printf("(2) %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	})

	log.Printf("Starting on port %s...", conf.ListenPort)
	if err := http.ListenAndServe(":" + conf.ListenPort, nil); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
}
