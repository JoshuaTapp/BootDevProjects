/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokeAPI"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pokedexcli",
	Short: "A Pokedex in your Command Line!",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pokedexcli.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	pokeAPI.GetLocations()
	cmdLoop()
}

func cmdLoop() {
	reader := bufio.NewReader(os.Stdin)
	errWriter := bufio.NewWriter(os.Stderr)
	//writer := bufio.NewWriter(os.Stdout)

	for {
		fmt.Print("pokedex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			errWriter.WriteString(err.Error())
			os.Exit(1)
		}

		input = strings.TrimSpace(input)

		cmd, args, err := rootCmd.Find(strings.Fields(input))
		//log.Default().Printf("CMD: %v, args: %v, error: %v\n", cmd.Short, args, err)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Available Commands:")
			for _, cmd := range rootCmd.Commands() {
				fmt.Println("  -", cmd.Name(), ": ", cmd.Short) // Print command names
			}
			continue
		}

		cmd.Run(cmd, args)

	}
}
