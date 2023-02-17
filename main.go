package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Game struct {
	Number	int
	User    string
	Bet	int
	Error   string
	Success string
	Status  int
	Win     bool
	Lose     bool
}

// Stocke la structure de jeux dans `t`
var t Game = Game{}

func manageForm(f url.Values){
	for key, value := range f {
		if key == "bet" {
			if i, err := strconv.Atoi(value[0]); err == nil {
				t.Bet = i
			} else {
				t.Error = "Votre choix doit être compris entre 0 et 5"
			}
			if t.Bet > 5 {
				t.Error = "Mise supérieure à 5"
			} else if t.Bet < 0 {
				t.Error = "Mise inférieure à 0"
			}
		}
	}
}

// Gestion du l'url /game
//
// C'est le coeur du jeu
func gameHandler(w http.ResponseWriter, r *http.Request) {
	t.Error = ""
	r.ParseForm()
	manageForm(r.Form)
	if (t.User == "") && (r.Form.Get("user") != "") {
		t.Win = false
		t.Lose = false
		t.Status = 1
		t.Number = rand.Intn(5) + 1
		t.User = r.Form.Get("user")
		log.Println("New User")
		log.Printf("Number %v\n", t.Number)
	} else if t.Bet == 0 {
	} else {
		// Recupération de la mise
		log.Printf("Bet %v - Number %v\n", t.Bet, t.Number)

		if t.Bet == t.Number {
			t.Win = true
		} else if t.Error != "" {
		} else {
			t.Lose = true
		}

	}
	var index = template.Must(template.ParseFiles("templates/base.html", "templates/game.tmpl"))
	err := index.ExecuteTemplate(w, "base.html", t)
	if err != nil {
		log.Fatalf("Template error : %s", err.Error())
	}
}

// Gestion de l'accueil du jeu
func indexHandler(w http.ResponseWriter, r *http.Request) {
	var index = template.Must(template.ParseFiles("templates/base.html", "templates/index.tmpl"))
	err := index.ExecuteTemplate(w, "base.html", t)
	if err != nil {
		log.Fatalf("Template error : %s", err.Error())
	}
}

// Gestion de l'authentification
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var login = template.Must(template.ParseFiles("templates/base.html", "templates/login.tmpl"))
	err := login.ExecuteTemplate(w, "base.html", t)
	if err != nil {
		log.Fatalf("Template error : %s", err.Error())
	}

}

// Gestion de la deconnexion
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	t.User = ""
	http.Redirect(w, r, "/login", http.StatusFound)
}

// Gestion de rejeux
func replayHandler(w http.ResponseWriter, r *http.Request) {
	t.Win = false
	t.Lose = false
	t.Bet = 0
	t.Number = rand.Intn(5) + 1
	http.Redirect(w, r, "/game", http.StatusFound)
	log.Printf("Number %v\n", t.Number)
}


// Fournit le jeu et attend que les joueurs arrivent
func main() {
	img := http.FileServer(http.Dir("./static/img/"))
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/img/", http.StripPrefix("/static/img/", img))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/replay", replayHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
