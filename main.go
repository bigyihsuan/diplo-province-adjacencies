package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	generateProvinces()
}

func generateProvinces() {
	provinces := []string{}
	err := filepath.WalkDir("./sample-orders", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		orderFile, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer orderFile.Close()
		provinces = append(provinces, parseOrders(orderFile)...)

		return nil
	})
	if err != nil {
		panic(err)
	}
	provinces = slices.DeleteFunc(provinces, func(s string) bool { return s == "" })
	sort.Strings(provinces)
	provinces = slices.Compact(provinces)
	// pretty.Println(orders)

	// generate for copy-paste into a const block
	for province := range slices.Values(provinces) {
		varName := province
		varName = strings.ReplaceAll(varName, "-", " ")
		varName = strings.ReplaceAll(varName, ".", " ")
		varName = strings.ReplaceAll(varName, "'", "")
		varName = cases.Title(language.AmericanEnglish).String(varName)
		varName = strings.ReplaceAll(varName, " ", "")
		fmt.Printf("%s = %q\n", varName, province)
	}
}

func parseOrders(orderFile io.Reader) (provinces []string) {
	textBytes, err := io.ReadAll(orderFile)
	if err != nil {
		panic(err)
	}
	text := string(textBytes)
	text = strings.ToLower(text)
	text = regexp.MustCompile(`_____\n*([a-zA-Z-]+\n*)?`).ReplaceAllString(text, "") // remove seperators
	text = regexp.MustCompile(`submitted orders:\n*`).ReplaceAllString(text, "")     // remove `submitted orders:`
	lines := strings.Split(text, "\n")
	for s := range slices.Values(lines) {
		provinces = append(provinces, cleanOrder(s)...)
	}

	return
}

func cleanOrder(order string) (provinces []string) {
	// strip fleet/army prefix
	order = strings.TrimPrefix(order, "f ")
	order = strings.TrimPrefix(order, "a ")
	// remove cores, holds
	order = strings.TrimSuffix(order, " cores")
	order = strings.TrimSuffix(order, " core")
	order = strings.TrimSuffix(order, " holds")
	order = strings.TrimSuffix(order, " hold")
	// split supports and convoys
	for _, r := range []string{` supports `, ` s `, ` convoys `} {
		parts := regexp.MustCompile(r).Split(order, -1)
		if len(parts) > 1 { // more than 1 element means it was a support/convoy order
			provinces = append(provinces, parts[0]) // append the supporting province
			order = parts[1]                        // there will always have at least 1 more order/province
		}
	}
	// split A - B
	splitOrder := strings.Split(order, " - ")
	if len(splitOrder) == 1 { // 1 element means no more orders to split
		provinces = append(provinces, order)
	} else {
		// otherwise, append the split orders
		provinces = append(provinces, splitOrder...)
	}
	// remove things like "coast" and "ec/nc/wc/sc" and "coast #1" etc
	for i, province := range slices.All(provinces) {
		province = regexp.MustCompile(` ([nsew]c)|(coast( #[0-9]+)?)`).ReplaceAllString(province, "")
		province = strings.ReplaceAll(province, "_", " ")
		province = strings.TrimSpace(province)
		provinces[i] = province
	}
	return
}
