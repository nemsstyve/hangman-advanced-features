package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func Selection_Mot(file string) string {
	//- 1.1 : Mise en place de l'aléatoire
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	//- 1.2 : Fetch du fichier words
	f, _ := os.Open(file)
	b1 := make([]byte, 9999)
	n1, _ := f.Read(b1)
	var words [][]byte
	index := 0
	for indice, lettre := range string(b1[:n1]) {
		if lettre == 10 {
			words = append(words, b1[index:indice])
			index = indice + 1
		}
	}
	//- 1.3 : Selection aléatoire
	mot_au_hasard := string(words[r1.Intn(len(words))])
	defer f.Close()
	return mot_au_hasard
}

func Affiche_Mot_Pendu(str string, valid []byte) {
	for i, char := range str {
		if i != 0 {
			fmt.Print(" ")
		}
		ecrit := 0
		for _, in_lv := range valid {
			if byte(char) == in_lv && ecrit == 0 {
				ecrit = 1
				fmt.Print(string(char - 'a' + 'A'))
			}
		}
		if ecrit == 0 {
			fmt.Print("_")
		}
	}
}

func Is_This_In_The_Word(lettre string, mot string) bool {
	for _, char := range mot {
		if string(char) == lettre {
			return true
		}
	}
	return false
}

func Toutes_Les_Lettres_Sont_Dans_Le_Mot(lettres []byte, mot string) bool {
	for _, char1 := range mot {
		ecrit := 0
		for _, char2 := range lettres {
			if byte(char1) == byte(char2) {
				ecrit = 1
			}
		}
		if ecrit == 0 {
			return false
		}
	}
	return true
}

func Mot_Exact(letters []byte, mot string) bool {
	for i := range mot {
		if letters[i] != []byte(mot)[i] {
			return false
		}
	}
	return true
}

func Ending(message string, position_pendu int, mot_au_hasard string, lettres_validées []byte) {
	hangman, _ := os.Open("positions/" + strconv.Itoa(position_pendu) + ".txt")
	hangman_tab := make([]byte, 9999)
	hangman_read, _ := hangman.Read(hangman_tab)
	fmt.Println(string(hangman_tab[:hangman_read]))
	Affiche_Mot_Pendu(mot_au_hasard, lettres_validées)
	fmt.Println("\n" + message)
}

func main() {
	press_on := 1
	for press_on == 1 {
		mot_au_hasard := Selection_Mot("words.txt")
		started := 0
		position_pendu := 0
		reussi := 0
		valid_letters := []byte{mot_au_hasard[len(mot_au_hasard)/2-1]}
		var letters_tried []byte
		var mots_essayées []string
		var error string
		for position_pendu < 10 {
			fmt.Print("\n[Attempt]\n")
			if started == 0 {
				fmt.Print("Good Luck, you have 10 attempts." + "\n")
				started = 1
			}
			if position_pendu > 0 && reussi == 0 {
				fmt.Print("Not present in the word, " + strconv.Itoa(10-position_pendu) + " attempts remaining\n\n")
			}
			if !Toutes_Les_Lettres_Sont_Dans_Le_Mot(valid_letters, mot_au_hasard) {
				hangman, _ := os.Open("positions/" + strconv.Itoa(position_pendu) + ".txt")
				hangman_tab := make([]byte, 9999)
				hangman_read, _ := hangman.Read(hangman_tab)
				fmt.Print(string(hangman_tab[:hangman_read]) + "\n")
				Affiche_Mot_Pendu(mot_au_hasard, valid_letters)
				fmt.Print("\n")
				fmt.Print("Already tried :")
				for _, mots := range mots_essayées {
					fmt.Print(" " + string(mots))
				}
				for _, ltr := range letters_tried {
					fmt.Print(" " + string(ltr))
				}
				fmt.Print("\n" + error)
				error = ""
				fmt.Print("\nProposition : ")
				var proposition string
				fmt.Scanln(&proposition)
				if len(proposition) == 1 {
					if Is_This_In_The_Word(proposition, string(letters_tried)) || Is_This_In_The_Word(string(proposition[0]-'A'+'a'), string(letters_tried)) || Is_This_In_The_Word(string(proposition[0]-'a'+'A'), string(letters_tried)) {
						error = "You already tried " + proposition
					} else {
						if Is_This_In_The_Word(proposition, mot_au_hasard) || Is_This_In_The_Word(string(proposition[0]-'A'+'a'), mot_au_hasard) {
							if Is_This_In_The_Word(proposition, mot_au_hasard) {
								valid_letters = append(valid_letters, proposition[0])
							}
							if Is_This_In_The_Word(string(proposition[0]-'A'+'a'), mot_au_hasard) {
								valid_letters = append(valid_letters, proposition[0]-'A'+'a')
							}
							reussi = 1
						} else if !(Is_This_In_The_Word(proposition, string(letters_tried))) && !(Is_This_In_The_Word(string(proposition[0]-'A'+'a'), string(letters_tried))) && !Is_This_In_The_Word(string(proposition[0]-'a'+'A'), string(letters_tried)) {
							letters_tried = append(letters_tried, proposition[0])
							position_pendu++
							reussi = 0
						}
					}
				}
				if len(proposition) > 1 {
					if Mot_Exact([]byte(proposition), mot_au_hasard) {
						fmt.Print("\n[Attempt]\n")
						valid_letters = []byte(proposition)
						Ending("Congrats !", position_pendu, mot_au_hasard, valid_letters)
						break
					} else {
						ecrit := 0
						for _, mot_essayé := range mots_essayées {
							if proposition == mot_essayé {
								ecrit = 1
							}
						}
						if ecrit == 1 {
							error = "You already tried " + proposition
						}
						if ecrit == 0 {
							mots_essayées = append(mots_essayées, proposition)
							position_pendu = position_pendu + 2
							reussi = 0
						}
					}
				}
			} else {
				fmt.Print("\n[End]\n")
				Ending("Congrats !", position_pendu, mot_au_hasard, valid_letters)
				break
			}
			fmt.Print("\n")
		}
		if position_pendu == 10 {
			fmt.Print("\n[End]\n")
			hangman, _ := os.Open("positions/10.txt")
			hangman_tab := make([]byte, 9999)
			hangman_read, _ := hangman.Read(hangman_tab)
			fmt.Print(string(hangman_tab[:hangman_read]) + "\n")
			Ending("You lose !", position_pendu, mot_au_hasard, []byte(mot_au_hasard))
		}
		var start_again string
		for start_again != "n" && start_again != "N" && start_again != "o" && start_again != "O" {
			fmt.Print("\nstart_again ? (o/n) : ")
			fmt.Scanln(&start_again)
			if start_again == "n" || start_again == "N" {
				press_on = 0
			}
			if start_again == "o" || start_again == "O" {
				press_on = 1
			}
		}
	}
}
