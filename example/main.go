package main

import (
	"fmt"
	"os"

	"github.com/Notation/solc-go"

	"github.com/pkg/errors"
)

func main() {
	var (
		version    = "0.4.25"
		file       = "./testdata/0.4.25.sol"
		solcBinary = "./solc_bin/soljson-v0.4.25+commit.59dbf8f1.js"
	)
	compiler, err := solc.NewFromFile(solcBinary, version)
	if err != nil {
		panic(errors.Wrap(err, "NewFromFile"))
	}
	fileData, err := os.ReadFile(file)
	if err != nil {
		panic(errors.Wrap(err, "ReadFile"))
	}

	input := &solc.Input{
		Language: "Solidity",
		Sources: map[string]solc.SourceIn{
			file: {Content: string(fileData)},
		},
		Settings: solc.Settings{
			Optimizer: solc.Optimizer{
				Enabled: false,
			},
			OutputSelection: map[string]map[string][]string{
				"*": {
					"*": []string{
						"metadata",
						"evm.bytecode",
						"evm.deployedBytecode",
						"evm.methodIdentifiers",
					},
					"": []string{
						"ast",
					},
				},
			},
		},
	}
	out, err := compiler.Compile(input)
	if err != nil {
		panic(errors.Wrap(err, "Compile"))
	}

	fmt.Println(out)
}
