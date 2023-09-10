// Copyright 2020 The casbin_ex Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package casbin

import (
	"bytes"
	"encoding/json"
)

func CasbinJsGetPermissionForUser(e IEnforcer, user string) (string, error) {
	model := e.GetModel()
	m := map[string]interface{}{}

	m["m"] = model.ToText()

	pRules := [][]string{}
	for ptype := range model["p"] {
		policies := model.GetPolicy("p", ptype)
		for _, rules := range policies {
			pRules = append(pRules, append([]string{ptype}, rules...))
		}
	}
	m["p"] = pRules

	gRules := [][]string{}
	for ptype := range model["g"] {
		policies := model.GetPolicy("g", ptype)
		for _, rules := range policies {
			gRules = append(gRules, append([]string{ptype}, rules...))
		}
	}
	m["g"] = gRules

	result := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(result)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(m)
	return result.String(), err
}
