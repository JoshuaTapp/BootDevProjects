package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokeAPI"
	"github.com/spf13/cobra"
)

// exploreCmd represents the explore command
var exploreCmd = &cobra.Command{
	Use:   "explore [location]",
	Short: "explore the provided area's pokemon",
	Long:  `Something about how to use this cmd`,
	// Removed the Args validation to use the custom validation below
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateExploreArgs(args); err != nil { // Custom validation
			fmt.Fprintln(os.Stderr, err) // Print error to stderr
			cmd.Usage()                  // Show usage help
		}
		fmt.Println("Exploring ", args[0], "...")
		explore(args[0])
	},
}

// function to validate the arguments
func validateExploreArgs(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide exactly one location name as an argument")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(exploreCmd)
}

func explore(locationName string) error {
	pokemon, err := pokeAPI.GetLocationAreaDetail().GetLocationPokemon(locationName)
	if err != nil {
		log.Println("error when getting Pokémon at: ", locationName)
		return err
	}
	if len(pokemon) < 1 {
		fmt.Println("No Pokémon found in ", locationName)
	} else {
		fmt.Println("Found Pokémon: ")
		for _, p := range pokemon {
			fmt.Println("  - ", p)
		}
	}
	return nil
}
