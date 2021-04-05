package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	imageType = 0
	videoType = 1
)

var template string = `<!doctype html><html>
<head><meta charset="utf-8"><meta http-equiv="refresh" content="10;url={{ link }}">
<style>html, body {background:#343837;margin:0;padding:0;height:100%;overflow:hidden;place-items:stretch center}
a {display:block;width:100%;height:100%;text-align:center}
img, video {object-fit:contain;margin:0 auto}</style></head>
<body><a href="{{ link }}">{{ content }}</a></body>
</html>`
var imageTemplate string = `<img src="{{ source }}" style="max-width:100%;height:100%">`
var videoTemplate string = `<video src="{{ source }}" style="max-width:100%;height:100%" autoplay loop></video>`

var images = []string{".gif", ".jpg", ".jpeg", ".png", ".webp"}
var videos = []string{".mp4", ".webm"}

func contains(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

func get(files []string, v string) (int, string) {
	if i, err := strconv.Atoi(v); err == nil {
		if i >= 0 && i < len(files) {
			return i, files[i]
		}
	}

	i := rand.Intn(len(files))

	return i, files[i]
}

func typeOf(file string) int {
	ext := filepath.Ext(file)

	if contains(images, ext) {
		return imageType
	}

	return videoType
}

func main() {
	formats := append(images, videos...)
	files := []string{}

	for _, format := range formats {
		matches, _ := filepath.Glob("*" + format)
		files = append(files, matches...)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		i, f := get(files, r.URL.Query().Get("i"))
		n := i + 1

		if n >= len(files) {
			n = 0
		}

		c := strings.ReplaceAll(template, "{{ link }}", "/?i="+fmt.Sprint(n))
		v := ""

		if typeOf(f) == imageType {
			v = imageTemplate
		} else {
			v = videoTemplate
		}

		v = strings.ReplaceAll(v, "{{ source }}", "/view?i="+fmt.Sprint(i))

		fmt.Fprint(w, strings.ReplaceAll(c, "{{ content }}", v))
	})

	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		_, f := get(files, r.URL.Query().Get("i"))
		http.ServeFile(w, r, f)
	})

	http.ListenAndServe(":8080", nil)
}
