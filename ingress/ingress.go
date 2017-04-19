package ingress

import (
	"context"
	"strings"
	"text/template"

	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var virtualHostEnvKey = "VIRTUAL_HOST"

type Ingress struct {
	client *client.Client
}

type VirtualHostService struct {
	VirtualHost string
	Vip         string
	Port        string
}

func NewIngress(client *client.Client) *Ingress {
	return &Ingress{
		client: client,
	}
}

func (i *Ingress) Render(tplPath string, wr io.Writer) error {
	var virtualHostServices = make(map[string]*VirtualHostService)
	services, _ := i.client.ServiceList(context.Background(), types.ServiceListOptions{})

	for _, service := range services {
		env := service.Spec.TaskTemplate.ContainerSpec.Env
		vips := service.Endpoint.VirtualIPs
		envMap := i.parseEnv(env)

		if v, ok := envMap[virtualHostEnvKey]; ok && (len(vips) > 0) {
			//remove repeat host
			for _, host := range v {
				virtualHostServices[host] = &VirtualHostService{
					VirtualHost: host,
					Vip:         strings.Split(service.Endpoint.VirtualIPs[0].Addr, "/")[0],
					Port:        "80",
				}
			}
		}
	}

	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return err
	}
	err = tpl.Execute(wr, virtualHostServices)
	if err != nil {
		return err
	}
	return nil
}

func (i *Ingress) parseEnv(env []string) map[string][]string {
	envMap := map[string][]string{}
	for _, e := range env {
		vs := strings.Split(e, "=")
		if len(vs) == 2 {
			envMap[vs[0]] = strings.Split(vs[1], ",")
		}
	}
	return envMap
}
