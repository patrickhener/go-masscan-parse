package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Nmaprun struct {
	XMLName          xml.Name `xml:"nmaprun"`
	Scanner          string   `xml:"scanner,attr"`
	Start            string   `xml:"start,attr"`
	Version          string   `xml:"version,attr"`
	XMLOutputVersion string   `xml:"xmloutputversion,attr"`
	Scaninfo         Scaninfo `xml:"scaninfo"`
	Hosts            []Host   `xml:"host"`
	Runstats         Runstats `xml:"runstats"`
}

type Scaninfo struct {
	XMLName  xml.Name `xml:"scaninfo"`
	Type     string   `xml:"type,attr"`
	Protocol string   `xml:"protocol,attr"`
}

type Host struct {
	XMLName xml.Name `xml:"host"`
	Endtime string   `xml:"endtime,attr"`
	Address Address  `xml:"address"`
	Ports   Ports    `xml:"ports"`
}

type Address struct {
	XMLName  xml.Name `xml:"address"`
	Addr     string   `xml:"addr,attr"`
	AddrType string   `xml:"addrtype,attr"`
}

type Ports struct {
	XMLName xml.Name `xml:"ports"`
	Port    Port     `xml:"port"`
}

type Port struct {
	XMLName  xml.Name `xml:"port"`
	Protocol string   `xml:"protocol,attr"`
	PortID   string   `xml:"portid,attr"`
	State    State    `xml:"state"`
}

type State struct {
	XMLName xml.Name `xml:"state"`
	State   string   `xml:"state:attr"`
	Reason  string   `xml:"reason,attr"`
	TTL     string   `xml:"reason_ttl,attr"`
}

type Runstats struct {
	XMLName  xml.Name `xml:"runstats"`
	Finished Finished `xml:"finished"`
	Hosts    Hosts    `xml:"hosts"`
}

type Finished struct {
	XMLName xml.Name `xml:"finished"`
	Time    string   `xml:"time,attr"`
	TimeStr string   `xml:"timestr,attr"`
	Elapsed string   `xml:"elapsed,attr"`
}

type Hosts struct {
	XMLName xml.Name `xml:"hosts"`
	Up      string   `xml:"up,attr"`
	Down    string   `xml:"down,attr"`
	Total   string   `xml:"total,attr"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <input-file.xml>\n", os.Args[0])
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var masscan Nmaprun

	if err := xml.Unmarshal(data, &masscan); err != nil {
		panic(err)
	}

	f, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, i := range masscan.Hosts {
		line := fmt.Sprintf("%s:%s\n", i.Address.Addr, i.Ports.Port.PortID)
		if _, err := w.WriteString(line); err != nil {
			panic(err)
		}
	}

	w.Flush()
	fmt.Printf(`
Masscan %s xmloutputformat %s
Scantype: %s
Protocol: %s
Duration: %s seconds

Total: %s
Up: %s
Down: %s

'output.txt' has been written successfully - Bye

`, masscan.Version, masscan.XMLOutputVersion, masscan.Scaninfo.Type, masscan.Scaninfo.Protocol, masscan.Runstats.Finished.Elapsed, masscan.Runstats.Hosts.Total, masscan.Runstats.Hosts.Up, masscan.Runstats.Hosts.Down)
}
