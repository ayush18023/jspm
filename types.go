package main

import (
	"encoding/json"
	"sync"
)

type kind uint8
type location uint8

type Token struct {
	Kind  kind
	Value string
}

type TokenizedPackage struct {
	PkgName string
	Tokens  []Token
}

const (
	Identifier kind = 0
	Number     kind = 1
	Star       kind = 2 // for *
	Any        kind = 3 // for x
	Tilde      kind = 4 // for ~
	Raise      kind = 5 // for ^
	Dot        kind = 6 // for .
	Nil        kind = 7
	Great      kind = 8
	Less       kind = 9
	Equal      kind = 10
	Hyphen     kind = 11
	Pipe       kind = 12
	At         kind = 13
	GreatEqual kind = 14
	LessEqual  kind = 15
	Eov        kind = 16
)

const (
	INCLUDED               location = 0
	INCLUDED_WRONG_VERSION location = 1
	NOT_INCLUDED           location = 2
)

type PackageResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Dist    struct {
		Tarball string `json:"tarball"`
	} `json:"dist"`
	Dependencies         map[string]string `json:"dependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
	Bin                  map[string]string `json:"bin"`
}

type AllPackageResponse struct {
	Name     string                     `json:"name"`
	Versions map[string]PackageResponse `json:"versions"`
}

type PackageJson struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
	raw          map[string]json.RawMessage
}

type LockJson struct {
	Packages map[string]PackageResponse `json:"packages"`
	raw      map[string]json.RawMessage
}

type DepManager struct {
	name string
	// installedDep []string
	mu *sync.RWMutex
	// deproot      *Package
	wg *sync.WaitGroup
	// slowroot     *Package
	depmap    map[string]PackageResponse
	cachelock map[string][]PackageResponse
	// pkgq         chan Package
}

type DefaultLockFile struct {
	Name string ``
}
