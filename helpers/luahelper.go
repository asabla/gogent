package helpers

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// TODO: refactor this flow into:
// 1. Always(?) generate lua stubs file based on available functions
//		* Register tool
//		* Making web requests
//		* Calling chat completion with/without history
//		* etc.
// 2. Read lua scripts which registers tools
// 3. Register tools based on these lua files
// 4. Run the main loop

type LuaFunc struct {
	Name     string
	Desc     string
	Args     []string
	Function func(L *lua.LState) int // WARN: not sure about this one
}

// Is used for settuping up lua embedded scripts. It generates a functions.lua file
// which contains all exposed functions for lua scripts in the same directory.
func SetupLua(scriptDir string, luaFuncs []LuaFunc) (*lua.LState, error) {

	// Check if scriptDir exists
	if _, err := os.Stat(scriptDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("Script directory does not exist: %s", scriptDir)
	}

	// Generate lua stubs
	stubs := ""
	for _, f := range luaFuncs {
		stubs += fmt.Sprintf("%s\n", generateLuaStub(f.Name, f.Desc, f.Args))
	}

	file, err := os.Create(fmt.Sprintf("%s/functions.lua", scriptDir))
	if err != nil {
		return nil, err
	}

	defer file.Close()

	_, err = file.WriteString(stubs)
	if err != nil {
		return nil, err
	}

	// TODO: maybe move this to a new function
	// since this logic should happened everytime a prompt is being
	// sent and processed.
	// Is it vyable to expose LState when implementing each function
	// which should be exposeable to lua scripts?

	// Create lua global state for defined functions
	// lstate := lua.NewState()
	// defer lstate.Close()

	// TOOD: create a global lua state for each registered function
	// and then invoke

	return nil, nil
}

func DoLuaStuff() {

	stubs := generateLuaStub("get_current_time", "Gets current time", []string{"city"})

	file, err := os.Create("gofunctions.lua")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteString(stubs)
	if err != nil {
		panic(err)
	}

	fmt.Println("Lua stubs generated successfully!")

	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("get_current_time", L.NewFunction(func(L *lua.LState) int {
		city := L.ToString(1)
		L.Push(lua.LString(time.Now().String()))
		fmt.Println("Called from Lua with arg:", city)

		return 1
	}))

	if err := L.DoFile("test.lua"); err != nil {
		log.Fatal(err)
	}
}

func generateLuaStub(funcName string, funcDesc string, parameters []string) string {
	params := ""

	if len(parameters) > 0 {
		params = parameters[0]
		for _, p := range parameters[1:] {
			params += ", " + p
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("-- %s\n", funcDesc))
	sb.WriteString(fmt.Sprintf("function %s(%s) end\n", funcName, params))

	return sb.String()
}
