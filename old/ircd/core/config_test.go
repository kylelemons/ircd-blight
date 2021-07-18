package core

import (
	"reflect"
	"testing"
)

var testDefaultConfig = &Configuration{
	Name:  "blight.local",
	SID:   "8LI",
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
	Ports: []*Ports{
		&Ports{
			PortString: "6666-6669",
		},
		&Ports{
			PortString: "6696-6699,9999",
			SSL:        "true",
		},
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
		t.Fatalf("ErrorMessage: %s", err)
	}
	if want := testDefaultConfig; !reflect.DeepEqual(got, want) {
		// TODO(kevlar): this is broken?
		t.Logf("config = %#v, want %#v", got, want)
	}
}

var portRangeTest = []struct {
	Range       string
	ExpectPorts []int
}{
	{
		Range:       "6667",
		ExpectPorts: []int{6667},
	},
	{
		Range:       "6666-6669",
		ExpectPorts: []int{6666, 6667, 6668, 6669},
	},
	{
		Range:       "6666-6669,6697",
		ExpectPorts: []int{6666, 6667, 6668, 6669, 6697},
	},
}

func TestPortRanges(t *testing.T) {
	for idx, test := range portRangeTest {
		directive := &Ports{
			PortString: test.Range,
		}
		ports, err := directive.GetPortList()
		if err != nil {
			t.Errorf("#%d: ErrorMessage: %s", idx, err)
			continue
		}
		if got, want := len(ports), len(test.ExpectPorts); got != want {
			t.Errorf("#%d: len(ports) = %d, want %d", idx, got, want)
			continue
		}
		for i := 0; i < len(ports); i++ {
			if got, want := ports[i], test.ExpectPorts[i]; got != want {
				t.Errorf("#%d: ports[%d] = %d, want %d", idx, i, got, want)
			}
		}
	}
}
