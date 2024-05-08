package kafka

import (
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

func TestListeners(t *testing.T) {
	stringB := strings.Builder{}

	names := []string{"alice", "bob", "charlie"}
	addresses := []string{"localhost1", "localhost2", "localhost3"}

	for i := 1; i <= 3; i++ {
		authM := "none"
		if i%2 == 0 {
			authM = "sasl"
		}

		ls := Listener{
			Name:                 "LISTENER_" + strings.ToUpper(names[i-1]),
			Address:              addresses[i-1],
			Port:                 i,
			AuthenticationMethod: authM,
		}

		stringB.WriteString(ls.String())
		if i < 3 {
			stringB.WriteString(",")
		}
	}

	expected := "LISTENER_ALICE://localhost1:1,LISTENER_BOB://localhost2:2,LISTENER_CHARLIE://localhost3:3"

	if stringB.String() != expected {
		t.Errorf("Expected %s, but got %s", expected, stringB.String())
	}
}

func TestListener_Parse(t *testing.T) {
	ls := Listener{
		Name:                 "LISTENER_WRONG_PORT",
		Address:              "0.0.0.0",
		Port:                 99999,
		AuthenticationMethod: "none",
	}

	if err := ls.Parse(); err == nil {
		t.Error("Expected an error, but got nil")
	}
}

func TestRegisterListeners(t *testing.T) {
	req := &testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Networks: []string{"network1", "network2"},
			NetworkAliases: map[string][]string{
				"network1": {"alias1", "alias2"},
				"network2": {"alias3", "alias4"},
			},
			Env: map[string]string{},
		},
	}

	settings := options{
		Listeners: []Listener{
			NewListener("LISTENER_ALICE", "localhost1", 49091),
			NewListener("LISTENER_BOB", "localhost2", 49092),
			NewListener("LISTENER_CHARLIE", "localhost3", 49093),
		},
	}

	err := registerListeners(settings, req)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// the expected advertised listeners environment variable is set
	expectedAdvertisedListeners := "LISTENER_ALICE://localhost1:49091,LISTENER_BOB://localhost2:49092,LISTENER_CHARLIE://localhost3:49093"
	if req.Env["KAFKA_ADVERTISED_LISTENERS"] != expectedAdvertisedListeners {
		t.Errorf("Expected %s, but got %s", expectedAdvertisedListeners, req.Env["KAFKA_ADVERTISED_LISTENERS"])
	}

	// the advertised listeners are set as network aliases
	expectedNetworkAliases := map[string][]string{
		"network1": {"alias1", "alias2", "localhost1", "localhost2", "localhost3"},
		"network2": {"alias3", "alias4", "localhost1", "localhost2", "localhost3"},
	}

	for network, aliases := range expectedNetworkAliases {
		if len(req.NetworkAliases[network]) != len(aliases) {
			t.Errorf("Expected %v, but got %v", aliases, req.NetworkAliases[network])
		}

		for i, alias := range aliases {
			if req.NetworkAliases[network][i] != alias {
				t.Errorf("Expected %s, but got %s", alias, req.NetworkAliases[network][i])
			}
		}
	}
}
