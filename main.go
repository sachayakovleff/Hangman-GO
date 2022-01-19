package main

import (
	"hangman-web/classic"
	"html/template"
	"net/http"
)

var temp *template.Template
var isstarted bool
var stats classic.Variables

func index(w http.ResponseWriter, r *http.Request) {
	if stats.Win {
		http.Redirect(w, r, "/win", http.StatusSeeOther)
	}
	if stats.Lose {
		http.Redirect(w, r, "/lost", http.StatusSeeOther)
	}
	temp = template.Must(template.ParseGlob("index.gohtml")) //instancie le fichier index.html
	if !isstarted {
		stats = classic.Start() //initialise les valeures des statistiques de la partie
		isstarted = true        //permet d'initialiser qu'une fois la partie
	}
	//stats variables: Realword, Hiddenword, Tested, Index, Hp, Win
	if r.Method == "POST" {
		stats.Currentword = r.FormValue("Hword") //get le mot sur le site                      //ceci est un test dans le terminal pour suivre la partie en cours
		stats = classic.Check(stats)             //teste le mot rentré par le joueur
		data := struct {                         //création d'un objet contenant une structure avec des variables
			Hword      string
			Wordtested [10]string
			Hp         int
			Hangman    [8]string
		}{
			Hword:      stats.Hiddenword, //Word ==> string  word ==> contenue du
			Wordtested: stats.Tested,
			Hp:         stats.Hp,
			Hangman:    stats.Hangman,
		}
		if stats.Lose {
			http.Redirect(w, r, "/lost", http.StatusSeeOther) //redirrection avant affichage
		}
		if stats.Win {
			http.Redirect(w, r, "/win", http.StatusSeeOther)
		}
		temp.ExecuteTemplate(w, "index.gohtml", data)
	}
	if r.Method == "GET" {
		data := struct { //création d'un objet contenant une structure avec des variables
			Hword      string
			Wordtested [10]string
			Hp         int
			Hangman    [8]string
		}{
			Hword:      stats.Hiddenword, //Word ==> string  word ==> contenue du string
			Wordtested: stats.Tested,
			Hp:         stats.Hp,
			Hangman:    stats.Hangman,
		}
		temp.ExecuteTemplate(w, "index.gohtml", data)
	}
}

func lost(w http.ResponseWriter, r *http.Request) {
	if !stats.Lose {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	temp = template.Must(template.ParseGlob("lost.gohtml"))
	data := struct { //création d'un objet contenant une structure avec des variables
		Word       string
		Hword      string
		Wordtested [10]string
		Hp         int
		Hangman    [8]string
	}{
		Word:       stats.Realword,
		Hword:      stats.Hiddenword, //Word ==> string  word ==> contenue du string
		Wordtested: stats.Tested,
		Hp:         stats.Hp,
		Hangman:    stats.Hangman,
	}
	if r.Method == "POST" {
		stats = classic.Start()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	temp.ExecuteTemplate(w, "lost.gohtml", data)
}

func win(w http.ResponseWriter, r *http.Request) {
	if !stats.Win {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	temp = template.Must(template.ParseGlob("win.gohtml"))
	data := struct { //création d'un objet contenant une structure avec des variables
		Word       string
		Hword      string
		Wordtested [10]string
		Hp         int
		Hangman    [8]string
	}{
		Word:       stats.Realword,
		Hword:      stats.Hiddenword, //Word ==> string  word ==> contenue du string
		Wordtested: stats.Tested,
		Hp:         stats.Hp,
		Hangman:    stats.Hangman,
	}
	if r.Method == "POST" {
		stats = classic.Start()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	temp.ExecuteTemplate(w, "win.gohtml", data)
}

func main() {
	http.HandleFunc("/", index) //préparer la page internet localhost:7777/
	http.HandleFunc("/lost", lost)
	http.HandleFunc("/win", win)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.ListenAndServe(":7777", nil) //définie le port 7777
}
