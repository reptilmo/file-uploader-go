// main.go
package main

import (
	"crypto/sha256"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"mime/multipart"
	"fmt"
)

var conf *Config

func main() {
	var err error

	conf, err = NewConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[1:] != "" {
			http.Error(w, "404 Nothing Here", http.StatusNotFound)
			return
		}
		t, _ := template.ParseFiles(conf.UploadDoc)
		t.Execute(w, nil)
	})

	http.HandleFunc("/upload", readMultiPart)

	log.Printf("Starting on port %s...", conf.ListenPort)
	if err := http.ListenAndServe(":"+conf.ListenPort, nil); err != nil {
		log.Fatal(err)
	}
}

func readMultiPart(w http.ResponseWriter, r *http.Request) {
	var err error
	var n int
	var reader *multipart.Reader
	var part *multipart.Part

	if reader, err = r.MultipartReader(); err != nil {
		log.Printf("Error (1) %s", err.Error())
		http.Error(w, "Server error", 500)
		return
	}

	var buffer []byte = make([]byte, 4096)
	for {
		if part, err = reader.NextPart(); err != nil {
			if err != io.EOF {
				log.Printf("Error (2) %s", err.Error())
				http.Error(w, "Server error", 500)
				return
			}

			w.WriteHeader(200)
			fmt.Fprintf(w, "Finished")
			return
		}

		var file *os.File
		var uploaded bool

		//log.Printf("Uploading %s %s", part.Header, part.FileName())
		file, err = os.OpenFile(conf.UploadPath+part.FileName(), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Printf("Error (3) opening local file %s", err.Error())
			http.Error(w, "Server error", 500)
			return
		}
		defer file.Close()
		sha := sha256.New()

		for !uploaded {
			if n, err = part.Read(buffer); err != nil {
				if err != io.EOF {
					log.Printf("Error (4) reading part %s", err.Error())
					http.Error(w, "Server error", 500)
					return
				}
				uploaded = true
			}

			if n, err = file.Write(buffer[:n]); err != nil {
				log.Printf("Error (5) writing part %s", err.Error())
				http.Error(w, "Server error", 500)
				return
			}

			sha.Write(buffer[:n])
		}

		log.Printf("Uploaded %s, %x", part.FileName(), sha.Sum(nil))
	}
}
