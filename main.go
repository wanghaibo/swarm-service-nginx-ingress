package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"log"

	"flag"

	"github.com/docker/docker/client"
	"github.com/wanghaibo/swarm-service-nginx-ingress/ingress"
)

func main() {
	var tpl, dst string
	flag.StringVar(&tpl, "tpl", "", "tpl path")
	flag.StringVar(&dst, "dst", "", "dst path")
	flag.Parse()

	if tpl == "" || dst == "" {
		log.Fatal("tpl and dst is required")
	}

	tmpDst, err := ioutil.TempFile(filepath.Dir(dst), "nginx-ingress")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		tmpDst.Close()
		os.Remove(tmpDst.Name())
	}()

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	ing := ingress.NewIngress(cli)
	err = ing.Render(tpl, tmpDst)
	if err != nil {
		log.Fatal(err)
	}

	os.Rename(tmpDst.Name(), dst)
}
