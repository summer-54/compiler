package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func compile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	name := uuid.New().URN()[9:]
	file, err := os.Create(filepath.Join(wd, name+".cpp"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	io.Copy(file, r.Body)
	cmd := exec.Command("/usr/local/bin/g++", "-O3", "-std=c++23", filepath.Join(wd, name+".cpp"), "-o", name)
	outfile, err := os.Create("./err-" + name)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
	cmd.Stdout = outfile
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	compiled, err := os.Open(filepath.Join(wd, name))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	io.Copy(w, compiled)
	os.Remove(filepath.Join(wd, name+".cpp"))
	os.Remove(filepath.Join(wd, name))
}

func main() {
	http.HandleFunc("/", compile)
	server := &http.Server{
		Addr:         ":54006",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	server.ListenAndServe()
}
