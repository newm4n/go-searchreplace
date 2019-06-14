package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

var (
	dummyContent = `
SubexpNames {placeholder} returns the names of the parenthesized 
subexpressions in this Regexp. The name for the first 
sub-expression is names[1], {placeholder}so that if m is a match slice, 
the name for m[i] is SubexpNames()[i]. 

{placeholder}

Since the Regexp as a whole cannot be named, 
names[0] is always the empty string. The slice should not 
be modified.{placeholder}
`
	newContent = `
SubexpNames replaced returns the names of the parenthesized 
subexpressions in this Regexp. The name for the first 
sub-expression is names[1], replacedso that if m is a match slice, 
the name for m[i] is SubexpNames()[i]. 

replaced

Since the Regexp as a whole cannot be named, 
names[0] is always the empty string. The slice should not 
be modified.replaced
`

	testFile = "/tmp/target.txt"
)

func TestReplaceIt(t *testing.T) {
	if _, err := os.Stat(testFile); os.IsExist(err) {
		if os.Remove(testFile) != nil {
			t.FailNow()
		}
	}
	if f, err := os.Create(testFile); err != nil {
		t.FailNow()
	} else {
		_, err := f.Write([]byte(dummyContent))
		if err != nil {
			t.FailNow()
		}
		if f.Close() != nil {
			t.FailNow()
		}

		regex := regexp.MustCompile("{placeholder}")

		if replaceIt(testFile, regex, "replaced") != nil {
			t.FailNow()
		} else {
			nudata, err := ioutil.ReadFile(testFile)
			if err != nil {
				t.FailNow()
			}
			if newContent != string(nudata) {
				t.Fail()
			}
			if os.Remove(testFile) != nil {
				t.Fail()
			}
		}
	}
}
