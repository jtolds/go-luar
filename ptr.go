// Copyright 2015 JT Olds
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package luar

import (
	"fmt"
	"reflect"

	"github.com/Shopify/go-lua"
)

func checkPtr(l *lua.State, index int) reflect.Value {
	ref, ok := lua.CheckUserData(l, index, ptrName).(reflect.Value)
	if !ok || ref.Kind() != reflect.Ptr {
		lua.ArgumentError(l, 1, "ptr expected")
		panic("unreached")
	}
	return ref
}

func pushAndSetupPtrTable(l *lua.State) {
	if !lua.NewMetaTable(l, ptrName) {
		return
	}
	lua.SetFunctions(l, []lua.RegistryFunction{
		{Name: "__index", Function: func(l *lua.State) int {
			defer fixPanics()
			return getField(l, checkPtr(l, 1), 2)
		}},
		{Name: "__newindex", Function: func(l *lua.State) int {
			defer fixPanics()
			return setField(l, checkPtr(l, 1), 2, 3)
		}},
		{Name: "__tostring", Function: func(l *lua.State) int {
			defer fixPanics()
			l.PushString(fmt.Sprintf("%#v", checkPtr(l, 1).Interface()))
			return 1
		}},
		{Name: "__unm", Function: func(l *lua.State) int {
			defer fixPanics()
			ref := checkPtr(l, 1)
			pushReflectedValue(l, ref.Elem())
			return 1
		}},
		{Name: "__eq", Function: func(l *lua.State) int {
			defer fixPanics()
			l.PushBoolean(checkPtr(l, 1) == checkPtr(l, 2))
			return 1
		}},
	}, 0)
}
