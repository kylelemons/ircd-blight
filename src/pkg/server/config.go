package server

import (
	"bytes"
	"io/ioutil"
	"os"
	"xml"
)

type Password struct {
	Type     string "attr"
	Password string "chardata"
}

type Oper struct {
	Name     string "attr"
	Password *Password
	Host     []string
	Flag     []string
}

type Class struct {
	Name string "attr"
	Host []string
	Flag []string
}

type Link struct {
	Name string "attr"
	Host []string
	Flag []string
}

type Network struct {
	Name        string "attr"
	Description string
	Link        []*Link
}

type Configuration struct {
	Name     string "attr"
	Admin    string
	Network  *Network
	Prefix   string
	Class    []*Class
	Operator []*Oper
}

var DefaultXML = `
<server name="blight.local">
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

func LoadConfigString(confxml string) os.Error {
	conf, err := parseXMLConfig([]byte(confxml))
	if err != nil {
		return err
	}
	Config = conf
	return nil
}

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
