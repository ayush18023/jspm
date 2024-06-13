package main

import "encoding/json"

func (pj *PackageJson) UnmarshalJSON(bytes []byte) error {
	pj.raw = make(map[string]json.RawMessage)
	if err := json.Unmarshal(bytes, &pj.raw); err != nil {
		return err
	}
	if dep, ok := pj.raw["dependencies"]; ok {
		if err := json.Unmarshal(dep, &pj.Dependencies); err != nil {
			return err
		}
	}
	return nil
}

func (pj *PackageJson) MarshalJSON() ([]byte, error) {
	pj.raw["dependencies"], _ = json.Marshal(pj.Dependencies)
	return json.MarshalIndent(pj.raw, "", "\t")
}

func (pj *LockJson) UnmarshalJSON(bytes []byte) error {
	pj.raw = make(map[string]json.RawMessage)
	if err := json.Unmarshal(bytes, &pj.raw); err != nil {
		return err
	}
	if dep, ok := pj.raw["packages"]; ok {
		if err := json.Unmarshal(dep, &pj.Packages); err != nil {
			return err
		}
	}
	return nil
}

func (pj *LockJson) MarshalJSON() ([]byte, error) {
	pj.raw["packages"], _ = json.Marshal(pj.Packages)
	return json.MarshalIndent(pj.raw, "", "\t")
}
