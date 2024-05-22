/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokeAPI"
	"github.com/spf13/cobra"
)

// mapbCmd represents the mapb command
var mapbCmd = &cobra.Command{
	Use:   "mapb",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		mapPrev()
	},
}

func init() {
	rootCmd.AddCommand(mapbCmd)
	pokeAPI.GetLocations()

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mapbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mapbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func mapPrev() {
	locPtr := pokeAPI.GetLocations()

	err := locPtr.GetPrevious()
	if err != nil {
		fmt.Println(err)
		return
	}
	locPtr.PrintLocations()
}
