/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokeAPI"
	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokedex"
	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect [Pokemon]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateInspectArgs(args); err != nil { // Custom validation
			fmt.Fprintln(os.Stderr, err) // Print error to stderr
			cmd.Usage()                  // Show usage help
		}
		inspectPokemon(args[0])
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inspectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func validateInspectArgs(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must enter a pokemon name to catch")
	}
	return nil
}

func inspectPokemon(name string) {
	if p, ok := pokedex.GetPokemon(name); ok {
		printPokemonDetails(&p)
	} else {
		fmt.Println("you have not caught that pokemon")
	}
}

func printPokemonDetails(p *pokeAPI.Pokemon) {
	name, height, weight := p.Name, p.Height, p.Weight
	stats := p.Stats
	types := p.Types

	fmt.Println("Name: ", name)
	fmt.Println("Height: ", height)
	fmt.Println("Weight: ", weight)

	fmt.Println("Stats:")
	if len(stats) < 1 {
		fmt.Println("N/A")
	} else {
		for _, s := range stats {
			fmt.Printf("\t- %v: %v\n", s.Stat.Name, s.BaseStat)
		}
	}

	fmt.Println("Types:")
	if len(types) < 1 {
		fmt.Println("N/A")
	} else {
		for _, t := range types {
			fmt.Printf("\t- %v\n", t.Type.Name)
		}
	}
}
