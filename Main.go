package main

import (
	"flag"
	"fmt"
	"github.com/newm4n/go-searchreplace/globber"
	"io/ioutil"
	"os"
	"regexp"
)

func main() {
	regexPtr := flag.String("pattern", "", "Regex search pattern")
	replacementPtr := flag.String("replacement", "", "Replacement text")
	folderPtr := flag.String("folder", "", "The base folder to look for")
	filterPtr := flag.String("filter", "/**/*", "File filter")

	if *regexPtr == "" || *replacementPtr == "" || *folderPtr == "" || *filterPtr == "" {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR : One of the required argument missing. fol=%s,fil=%s,reg=%s,rep=%s\n", *folderPtr, *filterPtr, *regexPtr, *replacementPtr)
		exitAndShowUsage()
	} else {
		regex := regexp.MustCompile(*regexPtr)
		fmt.Printf("Looking in        : %s\n", *folderPtr)
		fmt.Printf("File pattern      : %s\n", *filterPtr)
		fmt.Printf("Pattern to search : %s\n", *regexPtr)
		fmt.Printf("Replacement text  : %s\n", *replacementPtr)
		process(*folderPtr, *filterPtr, regex, *replacementPtr)
	}
}

func exitAndShowUsage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage : %s -pattern <string> -replacement <replacement> -folder <string> [-filter <string>]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func process(folder, filter string, regex *regexp.Regexp, replace string) {
	files, err := findPaths(folder, filter)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error : %s\n", err)
		exitAndShowUsage()
	} else {
		success := 0
		fail := 0
		for _, file := range files {
			err := replaceIt(file, regex, replace)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s ... error : %s\n", file, err)
				fail++
			} else {
				fmt.Printf("%s ... done\n", file)
				success++
			}
		}
		fmt.Printf("Processed %d files. %d processed successfuly. %d failed.", success+fail, success, fail)
	}
}

func replaceIt(target string, regex *regexp.Regexp, replace string) error {
	temp := fmt.Sprintf("%s.temp", target)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return err
	} else {
		data, err := ioutil.ReadFile(target)
		if err != nil {
			return err
		} else {
			newData := regex.ReplaceAllLiteralString(string(data), replace)
			ftemp, err := os.Create(temp)
			if err != nil {
				return err
			} else {
				_, err := ftemp.Write([]byte(newData))
				if err != nil {
					_ = ftemp.Close()
					nerr := os.Remove(temp)
					if nerr != nil {
						return nerr
					} else {
						return nil
					}
				} else {
					_ = ftemp.Close()
					err = os.Remove(target)
					if err != nil {
						return err
					}
					err = os.Rename(temp, target)
					if err != nil {
						return err
					} else {
						return nil
					}
				}
			}
		}
	}
}

func findPaths(dir, filter string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, file := range files {
		fullpath := fmt.Sprintf("%s/%s", dir, file.Name())
		if file.IsDir() {
			rfile, err := findPaths(fullpath, filter)
			if err != nil || rfile == nil {
				fmt.Printf("Returning from %s with empty result or error\n", fullpath)
			} else {
				ret = append(ret, rfile...)
			}
		} else {
			match, err := globber.IsPathMatch(filter, fullpath)
			if err != nil {
				fmt.Printf("Error matching filter %s, to file %s. Got %v", filter, fullpath, err)
			} else if match {
				ret = append(ret, fullpath)
			}
		}
	}
	return ret, nil
}
