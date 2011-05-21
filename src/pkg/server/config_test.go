package server

import (
	"testing"
	"reflect"
)

var testDefaultConfig = &Configuration{
	Name:  "blight.local",
	Admin: "Foo Bar [foo@bar.com]",
	Network: &Network{
		Name:        "IRCD-Blight",
		Description: "An unconfigured IRC network.",
		Link: []*Link{&Link{
			Name: "blight2.local",
			Host: []string{
				"blight2.localdomain.local",
				"127.0.0.1",
			},
			Flag: []string{
				"leaf",
			},
		}},
	},
	Class: []*Class{&Class{
		Name: "users",
		Host: []string{
			"*",
		},
		Flag: []string{
			"noident",
		},
	}},
	Operator: []*Oper{&Oper{
		Name: "god",
		Password: &Password{
			Type:     "plain",
			Password: "blight",
		},
		Host: []string{
			"127.0.0.1",
			"*.google.com",
		},
		Flag: []string{
			"admin",
			"oper",
		},
	}},
}

func TestDefaultConfig(t *testing.T) {
	got, err := parseXMLConfig([]byte(DefaultXML))
	if err != nil {
		t.Fatalf("Error: %s")
	}
	if want := testDefaultConfig; !reflect.DeepEqual(got, want) {
		t.Errorf("config = %#v, got %#v", got, want)
	}
}
