/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokeAPI"
	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokedex"
	"github.com/spf13/cobra"
)

// catchCmd represents the catch command
var catchCmd = &cobra.Command{
	Use:   "catch [Pokémon's name]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateCatchArgs(args); err != nil { // Custom validation
			fmt.Fprintln(os.Stderr, err) // Print error to stderr
			cmd.Usage()                  // Show usage help
		}
		fmt.Println("Throwing a pokeball at ", args[0], "...")
		isCaught, err := catch(args[0])
		if err != nil {
			log.Println(err)
			return
		}

		if isCaught {
			fmt.Printf("%s was caught!\n", args[0])
		} else {
			fmt.Printf("%s excaped!\n", args[0])
		}
	},
}

func validateCatchArgs(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must enter a pokemon name to catch")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(catchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func catch(pokemonName string) (bool, error) {
	p := pokeAPI.GetPokemon()
	err := p.GetPokemon(pokemonName)
	if err != nil {
		return false, err
	}

	if c := isCaught(p.BaseExperience); c {
		pokedex.AddPokemon(*p)
		return true, nil
	} 
	
	return false, nil
}

/*
isCaught determines whether a wild Pokémon is caught by the player based
on its base experience points. The higher the base experience points,
the lower the probability of catching the Pokémon.

@param baseXP The base experience points of the wild Pokémon.
@return true if the Pokémon is caught, false if not.
*/
func isCaught(baseXP int) bool {
	roll := rand.Float64()
	thresh := 1 - math.Pow(0.99, float64(baseXP-100)/508)
	return roll >= thresh
}
