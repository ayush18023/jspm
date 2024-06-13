package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func parseGreatThan(pkgName, lowerbound, pkgpth string, depm *DepManager) (PackageResponse, location) {
	// for condition ">"
	var isIncluded location = NOT_INCLUDED
	if depm.Consists(pkgpth + pkgName) {
		pkg := depm.Get(pkgpth + pkgName)
		if CompareVer(pkg.Version, lowerbound) == 1 {
			return pkg, INCLUDED
		} else {
			isIncluded = INCLUDED_WRONG_VERSION
		}
	}
	packages, ok := depm.CheckCache(pkgName)
	if ok && len(packages) != 0 {
		for _, cpkg := range packages {
			if CompareVer(cpkg.Version, lowerbound) == 1 {
				return cpkg, isIncluded
			}
		}
	}
	unpacked := getEntirePackage(pkgName)
	vv := unpacked.getVersions()
	i := binarySearch(lowerbound, vv)
	if vv[i] == lowerbound {
		if i+1 < len(vv) {
			return unpacked.Versions[vv[i+1]], isIncluded
		} else {
			return PackageResponse{}, isIncluded
		}
	}
	return unpacked.Versions[vv[i]], isIncluded
}

func parseGreatEqual(pkgName, lowerbound, pkgpth string, depm *DepManager) (PackageResponse, location) {
	// for condition ">="
	var isIncluded location = NOT_INCLUDED
	if depm.Consists(pkgpth + pkgName) {
		pkg := depm.Get(pkgpth + pkgName)
		if CompareVer(pkg.Version, lowerbound) >= 0 {
			return pkg, INCLUDED
		} else {
			isIncluded = INCLUDED_WRONG_VERSION
		}
	}
	packages, ok := depm.CheckCache(pkgName)
	if ok && len(packages) != 0 {
		for _, cpkg := range packages {
			if CompareVer(cpkg.Version, lowerbound) >= 0 {
				return cpkg, isIncluded
			}
		}
	}
	unpacked := getEntirePackage(pkgName)
	vv := unpacked.getVersions()
	i := binarySearch(lowerbound, vv)
	return unpacked.Versions[vv[i]], isIncluded
}

func parseLessThan(pkgName, lowerbound, pkgpth string, depm *DepManager) (PackageResponse, location) {
	// for condition "<"
	var isIncluded location = NOT_INCLUDED
	if depm.Consists(pkgpth + pkgName) {
		pkg := depm.Get(pkgpth + pkgName)
		if CompareVer(pkg.Version, lowerbound) == -1 {
			return pkg, INCLUDED
		} else {
			isIncluded = INCLUDED_WRONG_VERSION
		}
	}
	packages, ok := depm.CheckCache(pkgName)
	if ok && len(packages) != 0 {
		for _, cpkg := range packages {
			if CompareVer(cpkg.Version, lowerbound) == -1 {
				return cpkg, isIncluded
			}
		}
	}
	unpacked := getEntirePackage(pkgName)
	vv := unpacked.getVersions()
	i := binarySearch(lowerbound, vv)
	if vv[i] == lowerbound {
		if i-1 < len(vv) {
			return unpacked.Versions[vv[i-1]], isIncluded
		} else {
			return PackageResponse{}, isIncluded
		}
	}
	return unpacked.Versions[vv[i]], isIncluded
}

func parseLessEqual(pkgName, lowerbound, pkgpth string, depm *DepManager) (PackageResponse, location) {
	// for condition "<="
	var isIncluded location = NOT_INCLUDED
	if depm.Consists(pkgpth + pkgName) {
		pkg := depm.Get(pkgpth + pkgName)
		if CompareVer(pkg.Version, lowerbound) <= -1 {
			return pkg, INCLUDED
		} else {
			isIncluded = INCLUDED_WRONG_VERSION
		}
	}
	packages, ok := depm.CheckCache(pkgName)
	if ok && len(packages) != 0 {
		for _, cpkg := range packages {
			if CompareVer(cpkg.Version, lowerbound) <= -1 {
				return cpkg, isIncluded
			}
		}
	}
	unpacked := getEntirePackage(pkgName)
	vv := unpacked.getVersions()
	i := binarySearch(lowerbound, vv)
	return unpacked.Versions[vv[i]], isIncluded
}

func parseGreatLess(pkgName, lowerbound, upperbound, pkgpth string, depm *DepManager) (PackageResponse, location) {
	// for condition "> && <"
	var isIncluded location = NOT_INCLUDED
	if depm.Consists(pkgpth + pkgName) {
		pkg := depm.Get(pkgpth + pkgName)
		if CompareVer(pkg.Version, lowerbound) == 1 && CompareVer(pkg.Version, upperbound) == -1 {
			return pkg, INCLUDED
		} else {
			isIncluded = INCLUDED_WRONG_VERSION
		}
	}
	packages, ok := depm.CheckCache(pkgName)
	if ok && len(packages) != 0 {
		for _, cpkg := range packages {
			if CompareVer(cpkg.Version, lowerbound) == 1 && CompareVer(cpkg.Version, upperbound) == -1 {
				return cpkg, isIncluded
			}
		}
	}
	unpacked := getEntirePackage(pkgName)
	vv := unpacked.getVersions()
	i := binarySearch(lowerbound, vv)
	if i+1 < len(vv) && CompareVer(vv[i+1], upperbound) < 0 {
		return unpacked.Versions[vv[i+1]], isIncluded
	} else {
		return PackageResponse{}, isIncluded
	}
}

func parseGreatLessEqual(pkgName, lowerbound, upperbound, pkgpth string, depm *DepManager) (PackageResponse, location) {
	// for condition "> && <="
	var isIncluded location = NOT_INCLUDED
	if depm.Consists(pkgpth + pkgName) {
		pkg := depm.Get(pkgpth + pkgName)
		if CompareVer(pkg.Version, lowerbound) == 1 && CompareVer(pkg.Version, upperbound) <= 0 {
			return pkg, INCLUDED
		} else {
			isIncluded = INCLUDED_WRONG_VERSION
		}
	}
	packages, ok := depm.CheckCache(pkgName)
	if ok && len(packages) != 0 {
		for _, cpkg := range packages {
			if CompareVer(cpkg.Version, lowerbound) == 1 && CompareVer(cpkg.Version, upperbound) <= 0 {
				return cpkg, isIncluded
			}
		}
	}
	unpacked := getEntirePackage(pkgName)
	vv := unpacked.getVersions()
	i := binarySearch(lowerbound, vv)
	if i+1 < len(vv) && CompareVer(vv[i+1], upperbound) <= 0 {
		return unpacked.Versions[vv[i+1]], isIncluded
	} else {
		return PackageResponse{}, isIncluded
	}
}

func parseGreatEqualLess(pkgName, lowerbound, upperbound, pkgpth string, depm *DepManager) (PackageResponse, location) {
	// for condition ">= && <"
	var isIncluded location = NOT_INCLUDED
	if depm.Consists(pkgpth + pkgName) {
		pkg := depm.Get(pkgpth + pkgName)
		if CompareVer(pkg.Version, lowerbound) >= 0 && CompareVer(pkg.Version, upperbound) == -1 {
			return pkg, INCLUDED
		} else {
			isIncluded = INCLUDED_WRONG_VERSION
		}
	}
	packages, ok := depm.CheckCache(pkgName)
	if ok && len(packages) != 0 {
		for _, cpkg := range packages {
			if CompareVer(cpkg.Version, lowerbound) >= 0 && CompareVer(cpkg.Version, upperbound) == -1 {
				return cpkg, isIncluded
			}
		}
	}
	unpacked := getEntirePackage(pkgName)
	vv := unpacked.getVersions()
	i := binarySearch(lowerbound, vv)
	// // fmt.Println(vv, i)
	// fmt.Println("Collected package:", vv[i])
	if CompareVer(vv[i], upperbound) < 0 {
		return unpacked.Versions[vv[i]], isIncluded
	} else {
		return PackageResponse{}, isIncluded
	}
}

func parseGreatEqualLessEqual(pkgName, lowerbound, upperbound, pkgpth string, depm *DepManager) (PackageResponse, location) {
	// for condition ">= && <="
	var isIncluded location = NOT_INCLUDED
	if depm.Consists(pkgpth + pkgName) {
		pkg := depm.Get(pkgpth + pkgName)
		if CompareVer(pkg.Version, lowerbound) >= 0 && CompareVer(pkg.Version, upperbound) <= 0 {
			return pkg, INCLUDED
		} else {
			isIncluded = INCLUDED_WRONG_VERSION
		}
	}
	packages, ok := depm.CheckCache(pkgName)
	if ok && len(packages) != 0 {
		for _, cpkg := range packages {
			if CompareVer(cpkg.Version, lowerbound) >= 0 && CompareVer(cpkg.Version, upperbound) <= 0 {
				return cpkg, isIncluded
			}
		}
	}
	unpacked := getEntirePackage(pkgName)
	vv := unpacked.getVersions()
	i := binarySearch(lowerbound, vv)
	if CompareVer(vv[i], upperbound) <= 0 {
		return unpacked.Versions[vv[i]], isIncluded
	} else {
		return PackageResponse{}, isIncluded
	}
}

func parseExact(pkgName, lowerbound, pkgpth string, depm *DepManager) (PackageResponse, location) {
	// for condition "=="
	var isIncluded location = NOT_INCLUDED
	if depm.Consists(pkgpth + pkgName) {
		pkg := depm.Get(pkgpth + pkgName)
		if pkg.Version == lowerbound {
			return pkg, INCLUDED
		} else {
			isIncluded = INCLUDED_WRONG_VERSION
		}
	}
	packages, ok := depm.CheckCache(pkgName)
	if ok && len(packages) != 0 {
		for _, cpkg := range packages {
			if cpkg.Version == lowerbound {
				return cpkg, isIncluded
			}
		}
	}
	latestpkg := getSinglePackage(pkgName, lowerbound)
	return latestpkg, isIncluded
}

func parseAny(pkgName, pkgpth string, depm *DepManager) (PackageResponse, location) {
	//  checks if already collected else returns latest
	if depm.Consists(pkgpth + pkgName) {
		return depm.Get(pkgpth + pkgName), INCLUDED
	}
	return getSinglePackage(pkgName, "latest"), NOT_INCLUDED
}

func parseSemanticVer(
	vertokens []Token,
	pkgName, parentpth string,
	depm *DepManager,
) (PackageResponse, location) {
	var resolvedPkg PackageResponse
	var isIncluded location
	state := 0
	i := 0
	lowerbound := ""
	upperbound := ""
	n := len(vertokens)
	for i < len(vertokens) {
		token := vertokens[i]
		// fmt.Println("State: ", state, token)
		switch state {
		case 0:
			switch token.Kind {
			case Number:
				if strings.Contains(token.Value, "x") || len(token.Value) < 5 {
					resv := strings.ReplaceAll(token.Value, ".x", "")
					splitv := strings.Split(resv, ".")
					k := 0
					for k < 3-len(splitv) {
						resv += ".0"
						k += 1
					}
					lowerbound = resv
					if len(splitv) == 1 {
						num, err := strconv.Atoi(splitv[0])
						if err != nil {
							panic("Error in conversion to int")
						}
						upperbound = fmt.Sprintf("%d.0.0", num+1)
					} else {
						num, err := strconv.Atoi(splitv[1])
						if err != nil {
							panic("Error in conversion to int")
						}
						upperbound = fmt.Sprintf("%s.%d.0", splitv[0], num+1)
						// upperbound = resv[:2] + string(resv[2]+1) + ".0"
					}
					// greater than equal and less than
					// fmt.Println(lowerbound, upperbound)
					resolvedPkg, isIncluded = parseGreatEqualLess(pkgName, lowerbound, upperbound, parentpth, depm)
					i = n
				} else {
					// fmt.Println("is it here?")
					lowerbound = token.Value
					if n == 1 {
						//exact
						resolvedPkg, isIncluded = parseExact(pkgName, lowerbound, parentpth, depm)
						if resolvedPkg.Name != "" {
							i = n
						} else {
							state = 0
						}
						// // fmt.Println(resolvedPkg)
						i = n
					} else {
						// between
						state = 71
					}
				}
			case Pipe:
				lowerbound = ""
				upperbound = ""
				state = 0
			case Great:
				state = 11
				// >
			case Less:
				state = 21
			case GreatEqual:
				state = 31
			case LessEqual:
				state = 41
			case Tilde:
				state = 51
			case Raise:
				state = 61
			case Identifier:
				if token.Value == "latest" {
					resolvedPkg, isIncluded = parseExact(pkgName, "latest", parentpth, depm)
					if resolvedPkg.Name != "" { // Any greater than
						i = n
					} else {
						state = 0
					}
				} else if token.Value == "x" {
					resolvedPkg, isIncluded = parseAny(pkgName, parentpth, depm)
					if resolvedPkg.Name != "" {
						i = n
					} else {
						state = 0
					}
				}
			case Star:
				resolvedPkg, isIncluded = parseAny(pkgName, parentpth, depm)
				if resolvedPkg.Name != "" {
					i = n
				} else {
					state = 0
				}
			}
		case 11:
			if token.Kind != Number {
				state = 99
			}
			lowerbound = token.Value
			// > 1.2.3 || <= 2.0.5
			// > 1.2.3 <= 2.0.5
			// > 1.2.3 < 2.0.5
			if i+1 == n {
				resolvedPkg, isIncluded = parseGreatThan(pkgName, lowerbound, parentpth, depm)
				i = n
			} else if i+1 < n && vertokens[i+1].Kind == Pipe {
				resolvedPkg, isIncluded = parseGreatThan(pkgName, lowerbound, parentpth, depm)
				if resolvedPkg.Name != "" { // Any greater than
					i = n
				} else {
					state = 0
				}
			} else {
				state = 12
			}
		case 12:
			switch token.Kind {
			case Less:
				state = 13
			case LessEqual:
				state = 14
			}
		case 13: // great than and less than
			if token.Kind != Number {
				state = 99
			}
			upperbound = token.Value
			resolvedPkg, isIncluded = parseGreatLess(pkgName, lowerbound, upperbound, parentpth, depm)
			if i+1 < len(vertokens) && vertokens[i+1].Kind == Pipe && resolvedPkg.Name == "" {
				state = 0
			} else {
				i = n
			}
		case 14: // great than and less than Equal
			if token.Kind != Number {
				state = 99
			}
			upperbound = token.Value
			resolvedPkg, isIncluded = parseGreatLessEqual(pkgName, lowerbound, upperbound, parentpth, depm)
			if i+1 < len(vertokens) && vertokens[i+1].Kind == Pipe && resolvedPkg.Name == "" {
				state = 0
			} else {
				i = n
			}
		case 21:
			if token.Kind != Number {
				state = 99
			}
			lowerbound = token.Value
			if i+1 == n {
				resolvedPkg, isIncluded = parseLessThan(pkgName, lowerbound, parentpth, depm)
				i = n
			} else if i+1 < n && vertokens[i+1].Kind == Pipe {
				resolvedPkg, isIncluded = parseLessThan(pkgName, lowerbound, parentpth, depm)
				if resolvedPkg.Name != "" {
					i = n
				} else {
					state = 0
				}
			} else {
				state = 22
			}
		case 22:
			switch token.Kind {
			case Great:
				state = 23
			case GreatEqual:
				state = 24
			}
		case 23: // great than and less than
			if token.Kind != Number {
				state = 99
			}
			upperbound = token.Value
			resolvedPkg, isIncluded = parseGreatLess(pkgName, lowerbound, upperbound, parentpth, depm)
			if i+1 < len(vertokens) && vertokens[i+1].Kind == Pipe && resolvedPkg.Name == "" {
				state = 0
			} else {
				i = n
			}
		case 24: // greater than Equal and less than
			if token.Kind != Number {
				state = 99
			}
			upperbound = token.Value
			resolvedPkg, isIncluded = parseGreatEqualLess(pkgName, lowerbound, upperbound, parentpth, depm)
			if i+1 < len(vertokens) && vertokens[i+1].Kind == Pipe && resolvedPkg.Name != "" {
				state = 0
			} else {
				i = n
			}
		case 31:
			if token.Kind != Number {
				state = 99
			}
			lowerbound = token.Value
			if i+1 == n {
				resolvedPkg, isIncluded = parseGreatEqual(pkgName, lowerbound, parentpth, depm)
				i = n
			} else if i+1 < n && vertokens[i+1].Kind == Pipe {
				resolvedPkg, isIncluded = parseGreatEqual(pkgName, lowerbound, parentpth, depm)
				if resolvedPkg.Name != "" {
					i = n
				} else {
					state = 0
				}
			} else {
				state = 32
			}
		case 32:
			switch token.Kind {
			case Less:
				state = 33
			case LessEqual:
				state = 34
			}
		case 33: // Great than Equal and less than
			if token.Kind != Number {
				state = 99
			}
			upperbound = token.Value
			resolvedPkg, isIncluded = parseGreatEqualLess(pkgName, lowerbound, upperbound, parentpth, depm)
			if i+1 < len(vertokens) && vertokens[i+1].Kind == Pipe && resolvedPkg.Name == "" {
				state = 0
			} else {
				i = n
			}
		case 34: // Great than Equal and lesser than Equal
			if token.Kind != Number {
				state = 99
			}
			upperbound = token.Value
			resolvedPkg, isIncluded = parseGreatEqualLessEqual(pkgName, lowerbound, upperbound, parentpth, depm)
			if i+1 < len(vertokens) && vertokens[i+1].Kind == Pipe && resolvedPkg.Name == "" {
				state = 0
			} else {
				i = n
			}
		case 41:
			if token.Kind != Number {
				state = 99
			}
			lowerbound = token.Value
			if i+1 == n {
				resolvedPkg, isIncluded = parseLessEqual(pkgName, lowerbound, parentpth, depm)
				i = n
			} else if i+1 < n && vertokens[i+1].Kind == Pipe {
				resolvedPkg, isIncluded = parseLessEqual(pkgName, lowerbound, parentpth, depm)
				if resolvedPkg.Name != "" { // Any greater than
					i = n
				} else {
					state = 0
				}
			} else {
				state = 42
			}
		case 42:
			switch token.Kind {
			case Great:
				state = 43
			case GreatEqual:
				state = 44
			}
		case 43: // Less than Equal and Great than
			if token.Kind != Number {
				state = 99
			}
			upperbound = token.Value
			resolvedPkg, isIncluded = parseGreatLessEqual(pkgName, lowerbound, upperbound, parentpth, depm)
			if i+1 < len(vertokens) && vertokens[i+1].Kind == Pipe && resolvedPkg.Name == "" {
				state = 0
			} else {
				i = n
			}
		case 44: // Great than Equal and Less than Equal
			if token.Kind != Number {
				state = 99
			}
			upperbound = token.Value
			resolvedPkg, isIncluded = parseGreatEqualLessEqual(pkgName, lowerbound, upperbound, parentpth, depm)
			if i+1 < len(vertokens) && vertokens[i+1].Kind == Pipe && resolvedPkg.Name == "" {
				state = 0
			} else {
				i = n
			}
		case 51:
			resv := strings.ReplaceAll(token.Value, ".x", "")
			splitv := strings.Split(resv, ".")
			k := 0
			for k < 3-len(splitv) {
				resv += ".0"
				k += 1
			}
			lowerbound = resv
			if len(splitv) == 1 {
				num, err := strconv.Atoi(splitv[0])
				if err != nil {
					panic("Error in conversion to int")
				}
				upperbound = fmt.Sprintf("%d.0.0", num)
			} else {
				num, err := strconv.Atoi(splitv[1])
				if err != nil {
					panic("Error in conversion to int")
				}
				upperbound = fmt.Sprintf("%s.%d.0", splitv[0], num+1)
				// upperbound = resv[:2] + string(resv[2]+1) + ".0"
			}
			// greater than equal and less than
			resolvedPkg, isIncluded = parseGreatEqualLess(pkgName, lowerbound, upperbound, parentpth, depm)
			i = n
		case 61:
			resv := strings.ReplaceAll(token.Value, ".x", "")
			splitv := strings.Split(resv, ".")
			k := 0
			for k < 3-len(splitv) {
				resv += ".0"
				k += 1
			}
			lowerbound = resv
			num, err := strconv.Atoi(splitv[0])
			if err != nil {
				panic("Error in conversion to int")
			}
			upperbound = fmt.Sprintf("%d.0.0", num+1)
			// greater than equal and less than
			// fmt.Println(lowerbound, upperbound)
			resolvedPkg, isIncluded = parseGreatEqualLess(pkgName, lowerbound, upperbound, parentpth, depm)
			// fmt.Println(resolvedPkg, isIncluded)
			i = n
		case 71:
			if token.Kind != Hyphen {
				state = 99
			}
			state = 72
		case 72:
			if token.Kind != Number {
				state = 99
			}
			upperbound = token.Value
			// greater than equal and less than equal
			resolvedPkg, isIncluded = parseGreatEqualLessEqual(pkgName, lowerbound, upperbound, parentpth, depm)
			i = n
		default:
			log.Fatal("Oops!! Something went wrong")
		}
		i += 1
	}
	return resolvedPkg, isIncluded
}
