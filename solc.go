package solc

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/hashicorp/go-version"
	"rogchap.com/v8go"
)

// Build info & wasm binary can be found here:
// https://github.com/ethereum/solc-bin/tree/gh-pages/wasm

type Compiler struct {
	isolate *v8go.Isolate
	ctx     *v8go.Context

	mux      *sync.Mutex // protect underlying v8 context from concurrent access
	compiler *v8go.Value
	ver      *version.Version
}

func NewFromFile(file, ver string) (*Compiler, error) {
	v, err := version.NewVersion(ver)
	if err != nil {
		return nil, err
	}

	var (
		isolate = v8go.NewIsolate()
		ctx     = v8go.NewContext(isolate)
		c       = &Compiler{
			isolate: isolate,
			ctx:     ctx,
			ver:     v,
			mux:     &sync.Mutex{},
		}
	)

	soljson, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = c.init(string(soljson)); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Compiler) init(wasmScript string) error {
	var err error
	if _, err = c.ctx.RunScript(wasmScript, "main.js"); err != nil {
		return err
	}

	var (
		ver6, _ = version.NewVersion("0.6.0")
		ver5, _ = version.NewVersion("0.5.0")
	)
	if c.ver.LessThan(ver5) {
		c.compiler, err = c.ctx.RunScript("Module.cwrap('compileStandard', 'string', ['string', 'number'])",
			"wrap_compile.js")
	} else if c.ver.GreaterThanOrEqual(ver5) && c.ver.LessThan(ver6) {
		c.compiler, err = c.ctx.RunScript("Module.cwrap('solidity_compile', 'string', ['string', 'number'])",
			"wrap_compile.js")
	} else {
		c.compiler, err = c.ctx.RunScript("Module.cwrap('solidity_compile', 'string', ['string', 'number', 'number'])",
			"wrap_compile.js")
	}

	return err
}

func (c *Compiler) Compile(input *Input) (*Output, error) {
	fn, err := c.compiler.AsFunction()
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	value, err := v8go.NewValue(c.isolate, string(b))
	if err != nil {
		return nil, err
	}
	result, err := fn.Call(c.ctx.Global(), value)
	if err != nil {
		return nil, err
	}
	var output *Output
	if err = json.Unmarshal([]byte(result.String()), &output); err != nil {
		return nil, err
	}

	return output, nil
}

func (c *Compiler) Close() {
	c.mux.Lock()
	defer c.mux.Lock()
	c.ctx.Close()
	c.isolate.Dispose()
}
