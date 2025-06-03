package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var strp string
var intp int
var boolp bool

var flagsCmd = &cobra.Command{
	Use:   "flags",
	Short: "Experiment with flags",
	Long:  "A simple flags experimentation command, build with Cobra.",
	Run:   flagsFunc,
}

func flagsFunc(cmd *cobra.Command, args []string) {
	fmt.Println("string:", strp)
	fmt.Println("integer:", intp)
	fmt.Println("boolean:", boolp)
	fmt.Println("args:", args)
}
