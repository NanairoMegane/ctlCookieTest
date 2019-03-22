package main

import (
	"net/http"
)

func main() {

	/* testpageをサーブする */
	http.Handle("/", &templateHandler{filename: "/testpage.html"})

	/* 入力から文字列を取り出し、クッキー経由でページに反映する */
	http.HandleFunc("/show", showHandler)

	/* クッキーを削除する */
	http.HandleFunc("/delete", deleteHandler)

	/* 8080番ポートでlistenする */
	http.ListenAndServe(":8080", nil)
}
