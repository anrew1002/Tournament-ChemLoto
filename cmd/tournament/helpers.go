package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func (app *App) render(w http.ResponseWriter, status int, page string, data interface{}) {
	path := filepath.Join("web", "pages", page+".html")
	// if !ok {
	// 	err := fmt.Errorf("the template %s does not exist", page)
	// 			log.Println(err.Error())
	// 	http.Error(w, "Internal Error", http.StatusInternalServerError)
	// 	return
	// }
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}

	buf := new(bytes.Buffer)

	// err := ts.ExecuteTemplate(buf, "base", data)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	err = tmpl.Execute(buf, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
	}
	w.WriteHeader(status)

	buf.WriteTo(w)
}

type envelope map[string]any

func (app *App) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
