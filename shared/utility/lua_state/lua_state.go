package lua_state

import (
	"bufio"
	"os"
	"shared/utility/errors"
	"sync"

	"github.com/yuin/gopher-lua/parse"

	lua "github.com/yuin/gopher-lua"
)

type LuaMatcher struct {
	sync.Mutex
	*lua.LState
}

func NewLuaMatcher() *LuaMatcher {
	return &LuaMatcher{
		LState: lua.NewState(lua.Options{
			IncludeGoStackTrace: true,
		}),
	}
}

func (L *LuaMatcher) Init() error {
	lua.MaxArrayIndex = 2000
	proto, err := CompileFile("./Data/CombatPowerUtil.lua")
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = L.DoCompiledFile(proto)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (L *LuaMatcher) DoCompiledFile(proto *lua.FunctionProto) error {
	lfunc := L.NewFunctionFromProto(proto)
	L.Push(lfunc)
	return L.PCall(0, lua.MultRet, nil)
}

// 把lua脚本编译成字节码
func CompileFile(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

func (L *LuaMatcher) CalCharacterCombatPower(characInfo lua.LValue) (int32, error) {

	getCombatPower := L.GetGlobal("GetCharacterCombatPower")

	err := L.CallByParam(lua.P{
		Fn:      getCombatPower,
		NRet:    1,
		Protect: true,
		Handler: &lua.LFunction{},
	}, characInfo)

	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	ret := L.Get(-1)
	// 从堆栈中扔掉返回结果
	L.Pop(1)

	// 如果结果为数字
	res, ok := ret.(lua.LNumber)
	if !ok {
		return 0, errors.New("ErrRetNotLNumber")
	}

	return int32(res), nil
}

func (L *LuaMatcher) CalHeroCombatPower(heroInfo lua.LValue) (int32, error) {

	getHeroCombatPower := L.GetGlobal("GetHeroCombatPower")

	err := L.CallByParam(lua.P{
		Fn:      getHeroCombatPower,
		NRet:    1,
		Protect: true,
		Handler: &lua.LFunction{},
	}, heroInfo)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	ret := L.Get(-1)
	// 从堆栈中扔掉返回结果
	L.Pop(1)

	// 如果结果为数字
	res, ok := ret.(lua.LNumber)
	if !ok {
		return 0, errors.New("ErrRetNotLNumber")
	}

	return int32(res), nil
}

//==============Simple LuaPool?===============
type LStatePoolWithProto struct {
	m     sync.Mutex
	saved []*LuaMatcher
	Proto *lua.FunctionProto
}

func NewLuaPoolWithProto() *LStatePoolWithProto {
	lua.MaxArrayIndex = 2000
	proto, err := CompileFile("E:/code/lua/CombatPowerUtil_v9.lua")
	if err != nil {
		return nil
	}

	return &LStatePoolWithProto{
		saved: make([]*LuaMatcher, 0, 10),
		Proto: proto,
	}
}

func (pl *LStatePoolWithProto) New() *LuaMatcher {
	L := NewLuaMatcher()
	L.DoCompiledFile(pl.Proto)
	return L
}

func (pl *LStatePoolWithProto) Get() *LuaMatcher {
	pl.m.Lock()
	defer pl.m.Unlock()

	n := len(pl.saved)
	if n == 0 {
		return pl.New()
	}
	x := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]
	return x
}

func (pl *LStatePoolWithProto) Put(L *LuaMatcher) {
	pl.m.Lock()
	defer pl.m.Unlock()

	pl.saved = append(pl.saved, L)
}

// 关闭所有的虚拟机
func (pl *LStatePoolWithProto) Shutdown() {
	for _, L := range pl.saved {
		L.Close()
	}
}
