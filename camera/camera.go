package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var targetDir string

func main() {

	flag.StringVar(&targetDir, "targetDir", "/var/lib/motion/", "The location in which to store captured images")
	flag.Parse()
	http.HandleFunc("/take/", take)
	http.HandleFunc("/last/", last)
	log.Println("Listening on 8082")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func take(w http.ResponseWriter, r *http.Request) {

	_, err := http.Get("localhost:8080/0/action/snapshot")

	if err != nil {
		fmt.Println("Error connecting with motion. Try to restart motion or else")
		fmt.Println(" check the log files at /tmp/motion.log")
	}
	last(w, r)
}

func last(w http.ResponseWriter, r *http.Request) {

	img, err := ioutil.ReadFile(targetDir + "lastsnap.jpg")
	if err != nil {
		fmt.Println("File read error")
	}
	w.Write(img)
}
