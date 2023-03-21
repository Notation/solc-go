package solc

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Compile(t *testing.T) {
	tcs := []struct {
		Version    string
		File       string
		SolcBinary string
	}{
		{
			Version:    "0.4.25",
			File:       "./testdata/0.4.25.sol",
			SolcBinary: "./solc_bin/soljson-v0.4.25+commit.59dbf8f1.js",
		},
		{
			Version:    "0.5.0",
			File:       "./testdata/0.5.0.sol",
			SolcBinary: "./solc_bin/soljson-v0.5.0+commit.1d4f565a.js",
		},
		{
			Version:    "0.6.2",
			File:       "./testdata/0.6.2.sol",
			SolcBinary: "./solc_bin/soljson-v0.6.2+commit.bacdbe57.js",
		},
	}

	for _, tc := range tcs {
		compiler, err := NewFromFile(tc.SolcBinary, tc.Version)
		assert.Nil(t, err)
		if err != nil {
			continue
		}
		fileData, err := os.ReadFile(tc.File)
		assert.Nil(t, err)

		input := &Input{
			Language: "Solidity",
			Sources: map[string]SourceIn{
				tc.File: {Content: string(fileData)},
			},
			Settings: Settings{
				Optimizer: Optimizer{
					Enabled: false,
					Runs:    200,
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
		_, err = compiler.Compile(input)
		assert.Nil(t, err)
	}
}
