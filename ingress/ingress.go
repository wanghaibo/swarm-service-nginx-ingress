package ingress

import (
	"context"
	"strconv"
	"strings"
	"text/template"

	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	virtualHostEnvKey = "VIRTUAL_HOST"
	virtualPortEnvKey = "VIRTUAL_PORT"
)

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

		if hosts, ok := envMap[virtualHostEnvKey]; ok && len(vips) > 0 && len(hosts) > 0 {
			var port string
			if envPort, ok := envMap[virtualPortEnvKey]; ok {
				port = envPort
			} else if len(service.Endpoint.Ports) == 1 {
				port = strconv.Itoa(int(service.Endpoint.Ports[0].TargetPort))
			} else {
				port = "80"
			}

			virtualHostServices[hosts] = &VirtualHostService{
				VirtualHost: hosts,
				Vip:         strings.Split(service.Endpoint.VirtualIPs[0].Addr, "/")[0],
				Port:        port,
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

func (i *Ingress) parseEnv(env []string) map[string]string {
	envMap := map[string]string{}
	for _, e := range env {
		vs := strings.Split(e, "=")
		if len(vs) == 2 {
			if vs[0] == virtualHostEnvKey {
				envMap[vs[0]] = strings.Replace(vs[1], ",", " ", -1)
			} else {
				envMap[vs[0]] = vs[1]
			}
		}
	}
	return envMap
}
