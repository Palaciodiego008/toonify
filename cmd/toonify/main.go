package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Palaciodiego008/toonify"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "toonify",
	Short: "Convert between JSON and TOON formats",
	Long:  "A CLI tool to convert JSON to TOON format and vice versa",
}

var encodeCmd = &cobra.Command{
	Use:   "encode [file]",
	Short: "Convert JSON to TOON format",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input io.Reader = os.Stdin
		if len(args) > 0 {
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer file.Close()
			input = file
		}

		data, err := io.ReadAll(input)
		if err != nil {
			return err
		}

		result, err := encodeJSON(data)
		if err != nil {
			return err
		}

		fmt.Print(string(result))
		return nil
	},
}

var decodeCmd = &cobra.Command{
	Use:   "decode [file]",
	Short: "Convert TOON to JSON format",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input io.Reader = os.Stdin
		if len(args) > 0 {
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer file.Close()
			input = file
		}

		data, err := io.ReadAll(input)
		if err != nil {
			return err
		}

		result, err := decodeToon(data)
		if err != nil {
			return err
		}

		fmt.Print(string(result))
		return nil
	},
}

func encodeJSON(jsonData []byte) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	return toonify.EncodeBytes(data)
}

func decodeToon(toonData []byte) ([]byte, error) {
	var data interface{}
	if err := toonify.DecodeBytes(toonData, &data); err != nil {
		return nil, fmt.Errorf("invalid TOON: %v", err)
	}

	return json.MarshalIndent(data, "", "  ")
}

func main() {
	rootCmd.AddCommand(encodeCmd)
	rootCmd.AddCommand(decodeCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
