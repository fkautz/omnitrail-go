package omnitrail

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestEmpty(t *testing.T) {
	name := "empty"
	testAdd(t, name)
}

func TestOneFiles(t *testing.T) {
	name := "one-file"
	testAdd(t, name)
}

func TestTwoFiles(t *testing.T) {
	name := "two-files"
	testAdd(t, name)
}

func TestDeepStructure(t *testing.T) {
	name := "deep"
	testAdd(t, name)
}

func testAdd(t *testing.T, name string) {
	mapping := NewTrail()

	err := mapping.Add("./test/" + name)
	assert.NoError(t, err)

	expectedBytes, err := os.ReadFile("./test/" + name + ".json")
	assert.NoError(t, err)

	var expectedEnvelope Envelope
	err = json.Unmarshal(expectedBytes, &expectedEnvelope)
	assert.NoError(t, err)

	shortestExpectedKey := getShortestKey(&expectedEnvelope)
	shortestActualKey := getShortestKey(mapping.Envelope())

	for oldKey, val := range expectedEnvelope.Mapping {
		newKey := strings.Replace(oldKey, shortestExpectedKey, shortestActualKey, 1)
		//fmt.Printf("%s\n%s\n%s\n%s\n\n", shortestExpectedKey, shortestActualKey, oldKey, newKey)
		delete(expectedEnvelope.Mapping, oldKey)
		expectedEnvelope.Mapping[newKey] = val
	}

	assert.Equal(t, &expectedEnvelope, mapping.Envelope())

	res := formatADGString(mapping)

	expectedBytes, err = os.ReadFile("./test/" + name + ".adg")
	assert.NoError(t, err)
	assert.Equal(t, string(expectedBytes), res)
	if string(expectedBytes) != res {
		err = os.WriteFile("./"+name+".adg", []byte(res), 0644)
		assert.NoError(t, err)
	}
}

func getShortestKey(expectedEnvelope *Envelope) string {
	// get map keys
	keys := make([]string, 0, len(expectedEnvelope.Mapping))
	for key := range expectedEnvelope.Mapping {
		keys = append(keys, key)
	}

	// sort keys from shortest to longest
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) < len(keys[j])
	})

	shortestKey := keys[0]
	return shortestKey
}

func formatADGString(mapping Factory) string {
	res := ""
	sha1adgs := mapping.Sha1ADGs()
	// create a list of all keys sorted in lexical order
	keys := make([]string, 0, len(sha1adgs))
	for k := range sha1adgs {
		keys = append(keys, k)
	}
	// sort the keys
	sort.Strings(keys)

	for _, k := range keys {
		v := sha1adgs[k]
		if v != "" {
			res += fmt.Sprintln(k)
			res += fmt.Sprintln(v)
			res += fmt.Sprintln("--")
		}
	}
	res += fmt.Sprintln("----")

	keys = make([]string, 0, len(sha1adgs))
	sha2adgs := mapping.Sha256ADGs()
	keys = make([]string, 0, len(sha2adgs))
	for k := range sha2adgs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := sha2adgs[k]
		if v != "" {
			res += fmt.Sprintln(k)
			res += fmt.Sprintln(v)
			res += fmt.Sprintln("--")
		}
	}
	return res
}
