package main

import (
	"bufio"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	boundary := flag.Arg(0)
	buf := bufio.NewReader(os.Stdin)
	potentialPost, err := buf.Peek(4)
	if err != nil {
		log.Printf("Error reading! - %v", err)
	}
	if strings.HasPrefix(string(potentialPost), "POST") {
		buf.ReadLine() // ditch the POST request
		header, err := textproto.NewReader(buf).ReadMIMEHeader()
		if err != nil {
			log.Printf("Error reading headers - %v - either provide headers or specify the boundary as the first argument", err)
		}

		ct := header.Get("Content-Type")
		if strings.HasPrefix(ct, "multipart") {
			boundary = strings.Split(ct, "boundary=")[1]
		} else {
			log.Fatalf("Data provided is not multipart but - %s", ct)
		}
	}
	mr := multipart.NewReader(buf, boundary)
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			return // done!
		}
		if err != nil {
			log.Fatalf("Err opening part = %v", err)
		}
		file, err := ioutil.TempFile("", part.FileName())
		if err != nil {
			log.Fatalf("Err opening temp file = %v", err)
		}
		_, err = io.Copy(file, part)
		if err != nil {
			log.Fatalf("Err opening temp file = %v", err)
		}
		log.Printf("Wrote file - %s", file.Name())
		file.Close()
	}
}
