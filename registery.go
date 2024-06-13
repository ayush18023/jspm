package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func getSinglePackage(pkg_name string, version string) PackageResponse {
	resp, err := http.Get(REGISTRY + "/" + pkg_name + "/" + version)
	// req.Header = http.Header{
	// 	"Accept": {"application/vnd.npm.install-v1+json; q=1.0, application/json; q=0.8, */*"},
	// }
	// resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println(err)
	}
	var result PackageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		// fmt.Println(err)
	}
	return result
}

func getEntirePackage(pkg_name string) AllPackageResponse {
	resp, err := http.Get(REGISTRY + "/" + pkg_name + "/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println(err)
	}
	var result AllPackageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		// fmt.Println("Can not unmarshal JSON")
		// fmt.Println(err)
	}
	return result
}
