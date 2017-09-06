package ingress

import (
	"testing"

	"os"

	"github.com/docker/docker/client"
)

func TestRender(t *testing.T) {
	cli, err := client.NewEnvClient()
	if err != nil {
		t.Error(err)
	}

	ingress := NewIngress(cli)
	err = ingress.Render("./test.tpl", os.Stdout)
	if err != nil {
		t.Error(err)
	}
}
