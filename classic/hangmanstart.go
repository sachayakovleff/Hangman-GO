package classic

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type Variables struct {
	Tested      [10]string
	Realword    string
	Hiddenword  string
	Index       int
	Hp          int
	Win         bool
	Lose        bool
	Currentword string
	Hangman     [8]string
}

func Start() Variables { //initialisation des variables de la structure
	a := Variables{}
	for i := range a.Tested {
		a.Tested[i] = ""
	}
	a.Currentword = ""
	a.Win = false
	a.Lose = false
	a.Index = 0
	if len(os.Args) != 2 { //vérifie si la syntaxe de base est correcte
		fmt.Println("Error syntax : (go run main.go <file.txt>)")
		os.Exit(0)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error: File not found")
		os.Exit(0)
	}
	b, _ := ioutil.ReadAll(file)                  //récupérer l'ensemble des mots
	a.Realword = randomword(b)                    //choisir un mot aléatoirement
	nletter := len(a.Realword)/2 - 1              //permet de définir le nombre de lettre à afficher pour le mot caché
	a.Hiddenword = newstring(a.Realword, nletter) //permet de générer le mot caché
	a.Hp = 10                                     //définition du nombre de vies
	return a                                      //return de la structure
}

func randomword(s []byte) string {
	count := 0
	for _, v := range s {
		if v == '\n' {
			count++
		}
	}
	rand.Seed(time.Now().UnixNano())
	index := 0
	tabs := make([]string, count)
	for _, v := range s {
		if v == '\n' {
			index++
		}
		if v != '\n' && v != 13 && index != len(tabs) {
			tabs[index] += string(v)
		}
	}
	return tabs[rand.Intn(len(tabs))]
}

func changestring(s string, news string, nb int) string {
	nnews := ""
	for i := 0; i <= len(s)-1; i++ {
		if i != nb {
			nnews += string(news[i])
		}
		if i == nb {
			nnews += string(s[i])
		}
	}
	return nnews
}

func newstring(s string, nb int) string {
	var newmsg string
	for range s {
		newmsg += "_"
	}
	if nb != 0 {
		for i := 0; i < nb; i++ {
			indexts := rand.Intn(len(s) - 1)
			newmsg = changestring(s, newmsg, indexts)
		}
	}

	return newmsg
}

func Hangman(hp int) [8]string {
	f, _ := os.Open("hangman.txt")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var hangman [8]string
	var line int
	if hp < 0 {
		hp = 0
	}
	i := 0
	a := 71 - hp*8
	b := 78 - hp*8
	for scanner.Scan() {
		if line >= a && line <= b {
			hangman[i] = scanner.Text()
			i++
		}
		line++
	}
	return hangman
}

func Check(stats Variables) Variables {
	if stats.Currentword == stats.Realword { //le joueur a rentré le mot recherché
		stats.Win = true
		return stats
	}
	if len(stats.Currentword) > 1 && stats.Currentword != stats.Realword { //le joueur a rentré un mot qui n'est pas le mot recherché
		stats.Hp -= 2
		stats.Hangman = Hangman(stats.Hp)
		if stats.Hp <= 0 {
			stats.Lose = true
		}
		return stats
	}
	if len(stats.Currentword) == 1 {
		if istested(stats) {
			stats.Currentword = "lettre déjà testée"
		} else if verif(stats) {
			stats.Tested[stats.Index] = stats.Currentword
			stats.Index++
			stats.Hiddenword = goodletter(stats)
			if stats.Hiddenword == stats.Realword {
				stats.Win = true
				return stats
			}
		} else {
			stats.Tested[stats.Index] = stats.Currentword
			stats.Index++
			if stats.Hp > 1 {
				stats.Hp--
				stats.Hangman = Hangman(stats.Hp)
				return stats
			} else {
				stats.Hp--
				stats.Hangman = Hangman(stats.Hp)
				stats.Lose = true
			}
		}
	}
	return stats
}

func istested(stats Variables) bool {
	for _, v := range stats.Tested {
		if v == stats.Currentword {
			return true
		}
	}
	return false
}

func verif(stats Variables) bool {
	for _, v := range stats.Realword {
		if stats.Currentword == string(v) {
			return true
		}
	}
	return false
}

func goodletter(stats Variables) string {
	newhword := ""
	for i, v := range stats.Realword {
		if stats.Currentword == string(v) {
			newhword += string(v)
		} else {
			newhword += string(stats.Hiddenword[i])
		}
	}
	return newhword
}
