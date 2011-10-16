package core

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"xml"
)

// a Password stores Passwords for Oper and User directives.
type Password struct {
	Type     string `xml:"attr"`
	Password string `xml:"chardata"`
}

// An Oper is an operator configuration directive.
type Oper struct {
	Name     string `xml:"attr"`
	Password *Password
	Host     []string
	Flag     []string
}

// A Class is a user/server connection class directive.
type Class struct {
	Name string `xml:"attr"`
	Host []string
	Flag []string
}

// A Link represents the configuration information for a remote
// server link.
type Link struct {
	Name string `xml:"attr"`
	Host []string
	Flag []string
}

// A Ports direcive stores a port range and whether or not it is an SSL port.
type Ports struct {
	SSL        string `xml:"attr"`
	PortString string `xml:"chardata"`
}

// GetPortList gets the port list specified by the range(s) in this ports directive.
// The following port range formats are understood:
//   6667           // A single port
//   6666-6669      // A port range
//   6666-6669,6697 // Comma-separated ranges
func (p *Ports) GetPortList() (ports []int, err os.Error) {
	ranges := strings.Split(p.PortString, ",")
	for _, rng := range ranges {
		extremes := strings.Split(strings.TrimSpace(rng), "-")
		if len(extremes) > 2 {
			return nil, os.NewError("Invalid port range: " + rng)
		}
		low, err := strconv.Atoi(extremes[0])
		if err != nil {
			return nil, err
		}
		if len(extremes) == 1 {
			ports = append(ports, low)
			continue
		}
		high, err := strconv.Atoi(extremes[1])
		if err != nil {
			return nil, err
		}
		if low > high {
			return nil, os.NewError("Inverted range: " + rng)
		}
		for port := low; port <= high; port++ {
			ports = append(ports, port)
		}
	}
	return
}

// AreSSL returns true if the ports represented by the Ports directive
// are intended to be SSL-enabled ports.
func (p *Ports) AreSSL() bool {
	switch strings.ToLower(p.SSL) {
	case "1", "true", "yes", "on", "ssl", "enabled":
		return true
	}
	return false
}

// A Network represents the configuration data for the network on which
// this server is running.
type Network struct {
	Name        string "attr"
	Description string
	Link        []*Link
}

// A Configuration stores the configuration information for this server.
type Configuration struct {
	Name     string "attr"
	SID      string "attr"
	Admin    string
	Network  *Network
	Ports    []*Ports
	Class    []*Class
	Operator []*Oper
}

// A suitable default XML configuration file on which an admin should
// base his config.xml.
var DefaultXML = `` +
	`<server name="blight.local" sid="8LI">
	<ports>6666-6669</ports>
	<ports ssl="true">6696-6699,9999</ports>
	<network name="IRCD-Blight">
		<description>An unconfigured IRC network.</description>
		<link name="blight2.local">
			<host>blight2.localdomain.local</host>
			<host>127.0.0.1</host>
			<flag>leaf</flag>
		</link>
	</network>
	<admin>Foo Bar [foo@bar.com]</admin>
	<class name="users">
		<host>*</host>
		<flag>noident</flag>
	</class>
	<operator name="god">
		<password type="plain">blight</password>
		<host>127.0.0.1</host>
		<host>*.google.com</host>
		<flag>admin</flag>
		<flag>oper</flag>
	</operator>
</server>
`

var Config *Configuration

// LoadConfigFile loads an XML configuration string of the format shown in DefaultXML
// as the configuration for the server.
func LoadConfigString(confxml string) os.Error {
	conf, err := parseXMLConfig([]byte(confxml))
	if err != nil {
		return err
	}
	Config = conf
	return nil
}

// LoadConfigFile loads an XML configuration file of the format shown in DefaultXML
// as the configuration for the server.
func LoadConfigFile(filename string) os.Error {
	confxml, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	conf, err := parseXMLConfig([]byte(confxml))
	if err != nil {
		return err
	}
	Config = conf
	return nil
}

func parseXMLConfig(confxml []byte) (conf *Configuration, err os.Error) {
	conf = &Configuration{}
	buf := bytes.NewBuffer(confxml)
	err = xml.Unmarshal(buf, conf)
	return
}
