package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type MarcRecord struct {
	XMLName    xml.Name `xml:"record"`
	Leader     string   `xml:"leader"`
	CtrlFields []cField `xml:"controlfield"`
	DataFields []dField `xml:"datafield"`
}

type cField struct {
	Tag   string `xml:"tag,attr"`
	Field string `xml:",chadata"`
}

type dField struct {
	Tag       string     `xml:"tag,attr"`
	Ind1      string     `xml:"ind1,attr"`
	Ind2      string     `xml:"ind2,attr"`
	SubFields []subField `xml:"subfield"`
}

type subField struct {
	Code  string `xml:"code,attr"`
	Value string `xml:",chardata"`
}

// Default status for records missing a status for leader position 5
const DefaultStatus = "c"

func Parse(lines []string) (MarcRecord, error) {
	size := 24 // size of leader is 24 chars
	r := MarcRecord{}
	var ctrl cField
	var status string // for leader, position 5
	for _, line := range lines {
		if !strings.HasPrefix(line, "*") {
			println("cannot parse this line", line)
			continue
		}

		// increment size, needed for leader
		size = size + len(line)

		if strings.HasPrefix(line, "*00") {
			ctrl = parseCtrlField(line)
			if ctrl.Tag == "000" {
				status = ctrl.Field
			}
			r.CtrlFields = append(r.CtrlFields, ctrl)
		} else {
			tag := line[1:4]
			ind1 := line[4:5]
			ind2 := line[5:6]
			subFields := []subField{}
			for _, sf := range strings.Split(line[6:len(line)], "$") {
				if len(sf) == 0 {
					//println(tag, line)
					continue
				}
				subFields = append(subFields, subField{sf[0:1], sf[1:len(sf)]})
			}
			r.DataFields = append(r.DataFields, dField{tag, ind1, ind2, subFields})
		}
	}
	if status == "" {
		status = DefaultStatus
	}
	r.Leader = fmt.Sprintf("%05d%sam a22     1  4500", size, status)
	return r, nil
}

func parseCtrlField(line string) cField {
	tag := line[1:4]
	var v string
	switch tag {
	case "000":
		// make sure we have a character in position 5, if not
		// set the default status
		if len(line) < 10 || line[9:10] == " " {
			v = DefaultStatus
		} else {
			v = line[9:10]
		}
	case "008":
		// make sure field 008 has excactly 40 characters
		v = fmt.Sprintf("%-40s", line[4:len(line)])
	default:
		v = line[4:len(line)]
	}

	return cField{tag, v}
}

func main() {

	inFile := flag.String("i", "", "input file (line-MARC)")
	outFile := flag.String("o", "", "output file (MARCXML)")
	maxRecords := flag.Int("n", -1, "limit number of records to process (optional)")
	flag.Parse()

	if *inFile == "" || *outFile == "" {
		fmt.Println("Missing parameters:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	f, err := os.Open(*inFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	out, err := os.Create(*outFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	w := bufio.NewWriter(out)

	scanner := bufio.NewScanner(f)
	var record []string
	count := 0
	start := time.Now()
	for scanner.Scan() {
		line := scanner.Text()
		if line == "^" {
			if len(record) > 0 {
				mRec, _ := Parse(record)
				xmlRecord, err := xml.MarshalIndent(mRec, "", "  ")
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				_, err = w.Write(xmlRecord)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				_, err = w.WriteString("\n")
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				count = count + 1
				fmt.Printf("%d records processed\r", count)
				if count == *maxRecords {
					break
				}
			}
			record = []string{}
			continue
		}
		record = append(record, line)
	}

	w.Flush()
	diff := time.Now().Sub(start)
	fmt.Printf("Done, processed %d MARC records in %v\n", count, diff)

}
