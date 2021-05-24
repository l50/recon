/*
Copyright © 2021 Jayson Grace <jayson.e.grace@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	utils "github.com/l50/goutils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// asnCmd represents the asn command
var asnCmd = &cobra.Command{
	Use:   "asn",
	Short: "Discover and leverage ASNs associated with your target",
	Long:  `Discover and leverage ASNs associated with your target`,
	Run: func(cmd *cobra.Command, args []string) {
		targets, _ := cmd.Flags().GetString("targets")
		targetsFile, err := fileToSlice(targets)
		if err != nil {
			log.Fatalln(err)
		}
		logrus.Info("Gathering ASNs and associated IP ranges for the following: ", targetsFile)
		amassOutput, err := amassIntel(targetsFile)
		if err != nil {
			log.Fatalln(err)
		}
		asns, ipRanges := parseInput(amassOutput)
		printOutputs(asns, ipRanges)
	},
}

func init() {
	rootCmd.AddCommand(asnCmd)
	asnCmd.Flags().StringP("targets", "t", "", "Targets File")
}

func removeTrailingEmptyStringsInStringArray(sa []string) []string {
	lastNonEmptyStringIndex := len(sa) - 1
	for i := lastNonEmptyStringIndex; i >= 0; i-- {
		if sa[i] == "" {
			lastNonEmptyStringIndex--
		} else {
			break
		}
	}
	return sa[0 : lastNonEmptyStringIndex+1]
}

func fileToSlice(fileName string) ([]string, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return removeTrailingEmptyStringsInStringArray(strings.Split(string(b), "\n")), nil
}

func removeExtn(input string) string {
	if len(input) > 0 {
		if i := strings.LastIndex(input, "."); i > 0 {
			input = input[:i]
		}
	}
	return input
}

func amassIntel(targetsFile []string) ([]string, error) {
	var amassOutput []string
	for _, t := range targetsFile {
		logrus.Debug("Running command to gather ASNs and their IP Ranges: amass intel -org ", removeExtn(t))
		rawOut, err := utils.RunCommand("amass", "intel", "-org", removeExtn(t))
		if err != nil {
			return nil, err
		}
		// Split string on whitespace
		out := strings.Fields(rawOut)
		// Append split output to amassOutput
		for _, f := range out {
			amassOutput = append(amassOutput, f)
		}
	}
	return amassOutput, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func parseInput(amassOutput []string) ([]string, []string) {
	ipRangeRegex := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+\/[0-9]{2}$`)
	asnRegex := regexp.MustCompile(`^\d{5}`)
	var ipRanges []string
	var asns []string

	for _, line := range amassOutput {
		if ipRangeRegex.MatchString(line) {
			if !stringInSlice(line, ipRanges) {
				ipRanges = append(ipRanges, line)
			}
		} else if asnRegex.MatchString(line) {
			if !stringInSlice(line, asns) {
				asns = append(asns, line)
			}
		}
	}
	return asns, ipRanges
}

func printOutputs(asns []string, ipRanges []string) {
	fmt.Println("ASNs Found")
	fmt.Println("========================================")
	for _, a := range asns {
		logrus.Info("ASN: ", a+"\n")
		fmt.Printf("%s\n", a)
	}
	fmt.Println("\nIP Ranges Found")
	fmt.Println("========================================")
	for _, i := range ipRanges {
		logrus.Info("IP Range: ", i+"\n")
		fmt.Printf("%s\n", i)
	}
}
