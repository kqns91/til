package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/sugawarayuuta/sonnet"
)

const (
	inPath        = "sample.json"
	outPathStd    = "std.json"
	outPathSonnet = "sonnet.json"
)

type Response struct {
	Total  int                 `json:"total"`
	Result []Result `json:"result"`
}

type Result struct {
	ArtiCode  string   `json:"arti_code"`
	Title     string   `json:"title"`
	Member    string   `json:"member"`
	Date      string   `json:"date"`
	Link      string   `json:"link"`
	Images    []string `json:"images"`
	Highlight []string `json:"highlight"`
}

func main() {
	err := Std()
	if err != nil {
		log.Fatal(err)
	}

	err = Sonnet()
	if err != nil {
		log.Fatal(err)
	}
}

func Std() error {
	in, err := os.OpenFile(inPath, os.O_RDONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out, err := os.Create(outPathStd)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	var res Response

	err = json.NewDecoder(in).Decode(&res)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(out).Encode(res)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func Sonnet() error {
	in, err := os.OpenFile(inPath, os.O_RDONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out, err := os.Create(outPathSonnet)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	var res Response

	dec := sonnet.NewDecoder(in)

	err = dec.Decode(&res)
	if err != nil {
		log.Fatal(err)
	}

	enc := sonnet.NewEncoder(out)

	enc.SetEscapeHTML(true)

	err = enc.Encode(res)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
