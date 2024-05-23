package pokedex

import (
	"log"

	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokeAPI"
)

var (
	caughtPokemon map[string]pokeAPI.Pokemon
)

func init() {
	caughtPokemon = make(map[string]pokeAPI.Pokemon)
	log.Println("pokedex is online!")
}

func IsCaught(name string) bool {
	if _, ok := caughtPokemon[name]; ok {
		return true
	}
	return false
}

func AddPokemon(p pokeAPI.Pokemon) error {
	name := p.Name
	if IsCaught(name) {
		log.Println(name, " already caught!")
		return nil
	}

	caughtPokemon[name] = p
	log.Println("Adding ", name, " to pokedex!")
	return nil
}

func GetPokemon(name string) (p pokeAPI.Pokemon, b bool) {
	if IsCaught(name) {
		b, p = true, caughtPokemon[name]
	} 
	return
}

func GetPokedex() (*map[string]pokeAPI.Pokemon) {
	return &caughtPokemon
}
