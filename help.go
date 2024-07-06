package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func help() {
	titleStyle := lipgloss.NewStyle().Bold(true).Underline(true)

	fmt.Println("A JSON blob analyser, breaking down by field to visualise blob makeup.")
	fmt.Println()
	fmt.Println(titleStyle.Render("Usage:"))
	fmt.Println()
	fmt.Println("  blobby [OPTIONS] [FILE]")
	fmt.Println("  blobby [OPTIONS] (input via stdin)")
	fmt.Println()

	fmt.Println(titleStyle.Render("Options:"))
	fmt.Println()
	fmt.Println("  -d, the maximum depth to parse to (default 3)")
	fmt.Println("  -r, root of the blob to analyse. '.' seperated keys e.g. 'sprites.versions'")
}
