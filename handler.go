package main

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

/*
HTMLテンプレートをサーブするためのハンドラ
*/
type templateHandler struct {
	once     sync.Once          //HTMLテンプレートを１度だけコンパイルするための指定
	filename string             //テンプレートとしてHTMLファイル名を指定
	tmpl     *template.Template //テンプレート
}

/* templateHandlerをhttp.Handleに適合させるため、ServeHttpを実装する */
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// テンプレートディレクトリを指定する
	path, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	// 指定された名称のテンプレートファイルを一度だけコンパイルする
	t.once.Do(
		func() {
			t.tmpl = template.Must(template.ParseFiles(path + t.filename))
		})

	data := map[string]interface{}{
		"Host": r.Host,
	}

	// クッキーに入力文字が保存されていれば、ページに反映させる。
	msg, err := r.Cookie("msg")
	if msg != nil {
		data["msg"] = msg.Value
	}

	// テンプレートをパースする際、テンプレートに渡すデータも指定する。
	t.tmpl.Execute(w, data)
}

/*
formから文字列を読み取り、クッキーに詰めてから元のページに
リダイレクトさせるハンドラ。
*/
func showHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	input := r.FormValue("msg")
	data := &http.Cookie{
		Name:  "msg",
		Value: input,
		Path:  "/",
	}

	http.SetCookie(w, data)
	w.Header()["Location"] = []string{"/"}
	w.WriteHeader(http.StatusTemporaryRedirect)

}

/*
クッキーの内容を削除して、元のページにリダイレクトさせるハンドラ。
*/
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "msg",
		Value: "",
		Path:  "",
	})
	w.Header()["Location"] = []string{"/"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}
