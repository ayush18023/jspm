package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Version Semantics
// Patch releases: 1.0 or 1.0.x or ~1.0.4 (>= 1.0.4)
// Minor releases: 1 or 1.x or ^1.0.4(anything ie 1.0.5 or 1.2.0 but less than 2.0.0)
// Major releases: * or x any version

// test cases
// "foo": "1.0.0 - 2.9999.9999",
// "bar": ">=1.0.2 <2.1.2",
// "baz": ">1.0.2 <=2.3.4",
// "boo": "2.0.1",
// "qux": "<1.0.0 || >=2.3.1 <2.4.5 || >=2.5.2 <3.0.0",
// "asd": "http://asdf.com/asdf.tar.gz",
// "til": "~1.2",
// "elf": "~1.2.3",
// "two": "2.x",
// "thr": "3.3.x",
// "lat": "latest",
// "dyl": "file:../dyl"

func resolveCmd(stmnt string) {
	tokens := GetTokens(stmnt)
	// fmt.Println(tokens)
	parser(tokens)
}

func parser(tokens []Token) {
	// if tokens[0].Value == "purge" {
	// 	currPth := ""
	// 	newPth := filepath.Join(currPth, "node_modules")
	// 	purge(newPth)
	// } else {
	// 	parseInstall(tokens)
	// }
	if len(tokens) > 0 && tokens[0].Value == "purge" {
		currPth := ""
		newPth := filepath.Join(currPth, "node_modules")
		purge(newPth)
	} else if len(tokens) > 0 && tokens[0].Value == "install" {
		parseInstall(tokens)
	} else {
		fmt.Println("Welcome to jspm, available features:\n ")
		fmt.Println("\tinstall: specify the package to install pkg@version \n ")
		fmt.Println("\tpurge  : clears the node_modules directory \n ")
	}
}

func parseInstall(tokens []Token) {
	//vermapper keeps all package input data
	//1. Populate with package json then
	//2. Populate with the user input command
	//3. Then Install
	//4. Update package json then locked json
	fmt.Println("Preparing to install packages")
	vermapper := make(map[string][]Token)
	var mu = &sync.RWMutex{}
	var wg = &sync.WaitGroup{}
	var preprocess = &sync.WaitGroup{}
	depm := &DepManager{
		mu:        mu,
		wg:        wg,
		depmap:    make(map[string]PackageResponse),
		cachelock: make(map[string][]PackageResponse),
	}
	// var lockFile LockJson

	preprocess.Add(1)
	go func() {
		loadLockFile(depm)
		preprocess.Done()
	}()
	preprocess.Add(1)
	go func() {
		purge("./node_modules")
		_, err := os.Stat(MODULESBIN)
		if os.IsNotExist(err) {
			if err := os.Mkdir(MODULESBIN, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
		preprocess.Done()
	}()

	loadingChars := []string{"⠈⠁", "⠈⠑", "⠈⠱", "⠈⡱", "⢀⡱", "⢄⡱", "⢄⡱", "⢆⡱", "⢎⡱", "⢎⡰", "⢎⡠", "⢎⡀", "⢎⠁", "⠎⠁", "⠊⠁"}
	green := color.New(color.FgGreen).SprintFunc()
	l1 := loader{
		design: loadingChars,
		speed:  100 * time.Millisecond,
		preMsg: "  Discovering Packages\t ",
		posMsg: "",
		endMsg: fmt.Sprintf("  Discovering Packages\t %s  ", green("\u2713")),
	}
	start := time.Now()
	l1.Start()
	loadJsonPackages(vermapper)
	// fmt.Println("Loaded package file")
	loadCommandPackages(vermapper, tokens)
	// // fmt.Println(vermapper)
	preprocess.Wait()
	l1.Stop()
	l2 := loader{
		design: loadingChars,
		speed:  100 * time.Millisecond,
		preMsg: "  Fetching Packages\t ",
		posMsg: "",
		endMsg: fmt.Sprintf("  Fetching Packages\t %s  ", green("\u2713")),
	}
	l2.Start()
	getPackages(vermapper, depm)
	l2.Stop()
	l3 := loader{
		design: loadingChars,
		speed:  100 * time.Millisecond,
		preMsg: "  Installing Packages\t ",
		posMsg: "",
		endMsg: fmt.Sprintf("  Installing Packages\t %s  ", green("\u2713")),
	}
	l3.Start()
	installPackages(depm)
	updateJsons(depm, vermapper)
	// fmt.Println(err)
	l3.Stop()
	// updatePackageJson(vermapper)
	// updateLockJson(depm)
	yellow := color.New(color.FgYellow).SprintFunc()
	elapsed := time.Since(start)
	seconds := elapsed.Seconds()
	fmt.Printf("Installed in %s Packages installed %s", yellow(fmt.Sprintf("%.2fs", seconds)), yellow(len(vermapper)))
}

func loadLockFile(depm *DepManager) error {
	var lockFile *LockJson
	// err := readLockJson(lockFile)
	jsonFile, err := os.Open("locked.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	bytes, _ := io.ReadAll(jsonFile)
	if err := json.Unmarshal(bytes, &lockFile); err != nil {
		return err
	}
	for key, value := range lockFile.Packages {
		splits := strings.Split(key, "/")
		packageName := splits[len(splits)-1]
		depm.cachelock[packageName] = append(depm.cachelock[packageName], value)
	}
	return nil
}

func readLockJson(result *LockJson) error {
	jsonFile, err := os.Open("locked.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	return nil
}

func loadJsonPackages(vermapper map[string][]Token) error {
	var result PackageJson
	jsonFile, err := os.Open("package.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	for key, value := range result.Dependencies {
		vermapper[key] = GetTokens(value)
	}
	return nil
}

func loadCommandPackages(vermapper map[string][]Token, tokens []Token) {
	i := 1
	state := 0
	pkgname := ""
	for i < len(tokens) {
		token := tokens[i]
		switch state {
		case 0:
			if token.Kind != Identifier {
				state = 99
			} else {
				pkgname = token.Value
				state = 1
				if i+1 < len(tokens) && tokens[i+1].Kind == At {
					state = 1
				} else {
					vermapper[pkgname] = []Token{{
						Value: "latest",
						Kind:  Identifier,
					}}
					state = 0
				}
			}
		case 1:
			state = 2
		case 2:
			var start int
			start = i
			for i < len(tokens) && tokens[i].Kind != Identifier {
				i += 1
			}
			vermapper[pkgname] = tokens[start:i]
			i -= 1
			state = 0
		default:
			log.Fatal("Oops!! something's not Right")
		}
		i += 1
	}
}

func getPackages(vermapper map[string][]Token, depm *DepManager) {
	// // fmt.Println("Vermapper:", vermapper)
	for key, tokens := range vermapper {
		depm.wg.Add(1)
		go packageResolver(key, tokens, "", "/node_modules/", depm)
	}
	depm.wg.Wait()
}

func installPackages(depm *DepManager) {
	for key, value := range depm.depmap {
		depm.wg.Add(1)
		// fmt.Println("Adding: ", key)
		go downloadAndUntar(value.Dist.Tarball, ".", key, depm.wg)
	}
	depm.wg.Wait()
}

func packageResolver(pkgName string, vertokens []Token, parentname, parentpth string, depm *DepManager) {
	var resolvedPkg PackageResponse
	var isIncluded location
	resolvedPkg, isIncluded = parseSemanticVer(vertokens, pkgName, parentpth, depm)

	if resolvedPkg.Name != "" {
		if isIncluded != INCLUDED {
			var finalpath string = ""
			if isIncluded == NOT_INCLUDED {
				finalpath = parentpth
			} else if isIncluded == INCLUDED_WRONG_VERSION {
				finalpath = parentpth + parentname + "/node_modules/"
			}
			// fmt.Println("Adding: ", finalpath+pkgName)
			installpth := finalpath + pkgName
			depm.Add(installpth, resolvedPkg)
			for key, value := range resolvedPkg.Dependencies {
				depm.wg.Add(1)
				go packageResolver(key, GetTokens(value), pkgName, finalpath, depm)
			}
			for key, value := range resolvedPkg.OptionalDependencies {
				depm.wg.Add(1)
				go packageResolver(key, GetTokens(value), pkgName, finalpath, depm)
			}
			for key, value := range resolvedPkg.Bin {
				depm.wg.Add(1)
				go func(key, value string) {
					defer depm.wg.Done()
					createBinaries(installpth, value, key)
				}(key, value)
			}
		}
	} else {
		panic("No suitable version found for " + pkgName)
	}
	depm.wg.Done()
}

func get_version_from_tokens(tokens []Token) string {
	var result string
	for _, token := range tokens {
		result += token.Value
	}
	return result
}

func updatePackageJson(vermapper map[string][]Token) {
	jsonFile, err := os.ReadFile("package.json")
	if err != nil {
		log.Fatal(err)
	}
	var result PackageJson
	if err := result.UnmarshalJSON(jsonFile); err != nil {
		log.Fatal(err)
	}
	for key, value := range vermapper {
		result.Dependencies[key] = get_version_from_tokens(value)
	}
	updatedData, err := result.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	// Write updated JSON back to file
	if err := os.WriteFile("package.json", updatedData, 0644); err != nil {
		log.Fatal(err)
	}
}

func updateLockJson(depm *DepManager) {
	var result LockJson
	jsonFile, err := os.ReadFile("locked.json")
	if err == os.ErrNotExist {
		f, err := os.Create("locked.json")
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
		defaultLock := fmt.Sprintf(`{
			"packages":{}
		}`)
		_, err = f.Write([]byte(defaultLock))
		if err != nil {
			log.Fatal(err)
		}
		// jsonFile = []byte{
		// 	"name":{},
		// }
	} else {
		if err := result.UnmarshalJSON(jsonFile); err != nil {
			log.Fatal(err)
		}
	}
	for key, value := range depm.depmap {
		result.Packages[key] = value
	}
	updatedData, err := result.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("locked.json", updatedData, 0644); err != nil {
		log.Fatal(err)
	}
}

func createJsonFile(defaultPackage string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(defaultPackage))
	if err != nil {
		return err
	}
	return nil
}

func GetPackageJson() (*PackageJson, error) {
	jsonFile, err := os.ReadFile("package.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// fmt.Println("here it is ")
			defaultPackage := `{
	"dependencies":{}
}`
			err := createJsonFile(defaultPackage, "package.json")
			if err != nil {
				return nil, err
			}
			jsonFile = []byte(defaultPackage)
		} else {
			return nil, err
		}
	}
	var pacakgeFile PackageJson
	if err := pacakgeFile.UnmarshalJSON(jsonFile); err != nil {
		return nil, err
	}
	return &pacakgeFile, nil
}

func GetLockJson() (*LockJson, error) {
	jsonFile, err := os.ReadFile("locked.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			defaultPackage := `{
	"packages":{}
}`
			err := createJsonFile(defaultPackage, "locked.json")
			if err != nil {
				return nil, err
			}
			jsonFile = []byte(defaultPackage)
		} else {
			return nil, err
		}
	}
	var pacakgeFile LockJson
	if err := pacakgeFile.UnmarshalJSON(jsonFile); err != nil {
		return nil, err
	}
	return &pacakgeFile, nil
}

func updateJsons(depm *DepManager, vermapper map[string][]Token) error {
	packageJson, err := GetPackageJson()
	if err != nil {
		return err
	}
	// fmt.Println("Package json:", packageJson)
	lockJson, err := GetLockJson()
	if err != nil {
		return err
	}
	// fmt.Println("Lock json:", lockJson)
	for key, value := range depm.depmap {
		lockJson.Packages[key] = value
		// // fmt.Println(vermapper[key])
		if val, ok := vermapper[value.Name]; ok {
			if val[0].Value == "latest" {
				packageJson.Dependencies[value.Name] = "^" + value.Version
			} else {
				packageJson.Dependencies[value.Name] = get_version_from_tokens(val)
			}
		}
	}
	updatedData, err := lockJson.MarshalJSON()
	if err != nil {
		return err
	}
	if err := os.WriteFile("locked.json", updatedData, 0644); err != nil {
		return err
	}
	updatedData, err = packageJson.MarshalJSON()
	// fmt.Println(updatedData)
	if err != nil {
		return err
	}
	if err := os.WriteFile("package.json", updatedData, 0644); err != nil {
		return err
	}
	return nil
}

func purge(pathname string) error {
	d, err := os.Open(pathname)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(pathname, name))
		if err != nil {
			return err
		}
	}
	return nil
}
