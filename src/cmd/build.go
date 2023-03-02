/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bookmark-manager/actions"

	"github.com/spf13/cobra"
)

var Input, Output string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Long:  `Configure a bunch of links and a folder structure and this will build out the bookmarks for importing`,
	Run: func(cmd *cobra.Command, args []string) {
		actions.BuildStructure(Input, Output)

	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&Input, "inputDirectory", "i", ".", "file or directory where input yaml is.")
	buildCmd.Flags().StringVarP(&Output, "outputDirectory", "o", "output.html", "root of output directories.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
