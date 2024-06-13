package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func toBashPath(path string) string {
	return filepath.ToSlash(path)
}

func toWindowsPath(path string) string {
	return filepath.FromSlash(path)
}

func bashBinary(pkgpth, file string) string {
	binfilepth := filepath.Join(pkgpth, file)
	bashpth := toBashPath(binfilepth)
	adjustedpth := strings.Replace(bashpth, "node_modules", "..", 1)
	bashContent := fmt.Sprintf(`#!/bin/sh
basedir=$(dirname "$(echo "$0" | sed -e 's,\\,/,g')")

case  "$(uname) in
    *CYGWIN*|*MINGW*|*MSYS*) basedir=$(cygpath -w "$basedir");;
esac

if [ -x "$basedir/node" ]; then
  exec "$basedir/node"  "$basedir%s" "$@"
else 
  exec node  "$basedir%s" "$@"
fi
`, adjustedpth, adjustedpth)
	return bashContent
}

func cmdBinary(pkgpth, file string) string {
	binfilepth := filepath.Join(pkgpth, file)
	bashpth := toWindowsPath(binfilepth)
	adjustedpth := strings.Replace(bashpth, "node_modules", "..", 1)
	bashContent := `@ECHO off
GOTO start
:find_dp0
SET dp0=%~dp0
EXIT /b
:start
SETLOCAL
CALL :find_dp0

IF EXIST "%dp0%\node.exe" (
  SET "_prog=%dp0%\node.exe"
) ELSE (
  SET "_prog=node"
  SET PATHEXT=%PATHEXT:;.JS;=;%
)

endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & "%_prog%"  "%dp0%` + adjustedpth + `" %*`
	return bashContent
}

func pwsBinary(pkgpth, file string) string {
	binfilepth := filepath.Join(pkgpth, file)
	bashpth := toBashPath(binfilepth)
	adjustedpth := strings.Replace(bashpth, "node_modules", "..", 1)
	bashContent := fmt.Sprintf(`#!/usr/bin/env pwsh
$basedir=Split-Path $MyInvocation.MyCommand.Definition -Parent

$exe=""
if ($PSVersionTable.PSVersion -lt "6.0" -or $IsWindows) {
  # Fix case when both the Windows and Linux builds of Node
  # are installed in the same directory
  $exe=".exe"
}
$ret=0
if (Test-Path "$basedir/node$exe") {
  # Support pipeline input
  if ($MyInvocation.ExpectingInput) {
    $input | & "$basedir/node$exe"  "$basedir%s" $args
  } else {
    & "$basedir/node$exe"  "$basedir%s" $args
  }
  $ret=$LASTEXITCODE
} else {
  # Support pipeline input
  if ($MyInvocation.ExpectingInput) {
    $input | & "node$exe"  "$basedir%s" $args
  } else {
    & "node$exe"  "$basedir%s" $args
  }
  $ret=$LASTEXITCODE
}
exit $ret

`, adjustedpth, adjustedpth, adjustedpth, adjustedpth)
	return bashContent
}

func createFileAndWrite(fname, content string) error {
	file, err := os.Create(fname)
	if err != nil {
		// fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		// fmt.Println("Error writing to file:", err)
		return err
	}
	return nil
}

func createBinaries(pkgpth, file, fname string) {
	bashf := bashBinary(pkgpth, file)
	// fmt.Println(bashf)
	cmdf := cmdBinary(pkgpth, file)
	// fmt.Println(cmdf)
	pwsf := pwsBinary(pkgpth, file)
	// fmt.Println(pwsf)
	os.WriteFile(
		filepath.Join(MODULESBIN, fname),
		[]byte(bashf),
		os.ModePerm,
	)
	// fmt.Println(err)
	// fmt.Println(err)
	os.WriteFile(
		filepath.Join(MODULESBIN, fname+".cmd"),
		// fmt.Sprintf("%s/%s.cmd", MODULESBIN, fname),
		[]byte(cmdf),
		os.ModePerm,
	)
	// err := createFileAndWrite(fmt.Sprintf("./node_modules/.bin/%s.sh", fname), bashf)
	// fmt.Println(err)
	// createFileAndWrite(filepath.Join(MODULESBIN, fname+".cmd"), cmdf)
	// createFileAndWrite(filepath.Join(MODULESBIN, fname+".ps1"), pwsf)
	os.WriteFile(
		filepath.Join(MODULESBIN, fname+".ps1"),
		// fmt.Sprintf("%s/%s.ps1", MODULESBIN, fname),
		[]byte(pwsf),
		os.ModePerm,
	)
}
