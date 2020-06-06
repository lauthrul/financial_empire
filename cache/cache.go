package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

// cache iem structure
type Item struct {
	Names        []string   `json:"names"`
	Combinations [][]string `json:"combinations"`
}

// cache structure
type Cache struct {
	Items map[string]Item `json:"items"` // map key is md5Sum of (sorted) Item.Names
}

// default cache file
const cacheFile = ".cache"

// load cache from file
func LoadCache(cache *Cache) error {
	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, cache)
}

// get md5 sum with given name slice
func md5Sum(names []string) string {
	sort.Slice(names, func(i, j int) bool {
		return names[i] > names[j]
	})
	str := strings.Join(names, ",")
	sum := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", sum)
}

// check if exist the given names in cache
func (this *Cache) Exist(names []string) bool {
	sum := md5Sum(names)
	_, ok := this.Items[sum]
	return ok
}

// read data from cache
func (this *Cache) Read(names []string, result *[][]string) error {
	if !this.Exist(names) {
		return fmt.Errorf("%v not exist", names)
	}
	sum := md5Sum(names)
	*result = this.Items[sum].Combinations
	return nil
}

// save data to cache
func (this *Cache) Save(names []string, result [][]string) error {
	sum := md5Sum(names)
	if this.Items == nil {
		this.Items = make(map[string]Item)
	}
	this.Items[sum] = Item{
		Names:        names,
		Combinations: result,
	}
	data, err := json.MarshalIndent(this, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cacheFile, data, os.FileMode(644))
}
