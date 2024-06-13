package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Algo part
func splitVersion(version string) []int {
	onlyver := strings.Split(version, "-")[0]
	components := strings.Split(onlyver, ".")
	intComponents := make([]int, len(components))
	for i, component := range components {
		intVal, err := strconv.Atoi(component)
		if err != nil {
			panic(fmt.Sprintf("Invalid version component: %s", component))
		}
		intComponents[i] = intVal
	}
	return intComponents
}

func CompareVer(v1, v2 string) int {
	if v1 == v2 {
		return 0
	}
	vi1 := splitVersion(v1)
	vi2 := splitVersion(v2)
	// // fmt.Println(v1, vi1)
	// // fmt.Println(v2, vi2)
	for k := 0; k < len(vi1) && k < len(vi2); k++ {
		if vi1[k] > vi2[k] {
			return 1
		} else if vi1[k] < vi2[k] {
			return -1
		}
	}
	return 0
}

// func isGreaterThan(check, pivot string) bool {
// 	if CompareVer(check, pivot) == 1 {
// 		return true
// 	}
// 	return false
// }

// func isLesserThan(check, pivot string) bool {
// 	if CompareVer(check, pivot) == -1 {
// 		return true
// 	}
// 	return false
// }

// func isEqual(check, pivot string) bool {
// 	if CompareVer(check, pivot) == 0 {
// 		return true
// 	}
// 	return false
// }

func SortVersions(versions []string) {
	sort.Slice(versions, func(i, j int) bool {
		v1 := splitVersion(versions[i])
		v2 := splitVersion(versions[j])
		for k := 0; k < len(v1) && k < len(v2); k++ {
			if v1[k] > v2[k] {
				return false
			} else if v1[k] < v2[k] {
				return true
			}
		}
		return len(v2) > len(v1)
	})
}

func binarySearch(cur_version string, versions []string) int {
	low := 0
	high := len(versions)
	found := false
	selectedi := 0
	// fmt.Println("Version recieved", cur_version)
	for !found && low <= high && low > 0 {
		mid := low + (high-low)/2
		cmp := CompareVer(cur_version, versions[mid])
		// fmt.Println(cmp, versions[mid])
		// cmp := strings.Compare(cur_version, versions[mid])
		if cmp == 0 {
			found = true
			selectedi = mid
		} else if cmp < 0 {
			// if mid == 0{
			// 	found = true
			// 	selectedi = mid
			// } else
			if CompareVer(cur_version, versions[mid-1]) > 0 {
				found = true
				selectedi = mid - 1
			} else {
				high = mid
			}
		} else if cmp > 0 {
			// if mid == len(versions){
			// 	found = true
			// 	selectedi = mid
			// } else
			if CompareVer(cur_version, versions[mid+1]) < 0 {
				found = true
				selectedi = mid
			} else {
				low = mid
			}
		}
	}
	return selectedi
}

func (ap *AllPackageResponse) getVersions() []string {
	keys := make([]string, 0, len(ap.Versions))
	for k := range ap.Versions {
		keys = append(keys, k)
	}
	SortVersions(keys)
	// fmt.Println(keys)
	return keys
}

// Dependency manager
func (d *DepManager) Consists(pkgpath string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, ok := d.depmap[pkgpath]
	return ok
}

func (d *DepManager) Get(pkgpath string) PackageResponse {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.depmap[pkgpath]
}

func (d *DepManager) Add(key string, value PackageResponse) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.depmap[key] = value
}

func (d *DepManager) CheckCache(key string) ([]PackageResponse, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	pkg, ok := d.cachelock[key]
	return pkg, ok
}
