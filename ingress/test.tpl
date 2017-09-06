{{range $host,$value := .}}
server {
    server_name {{$host}}
    listen 80;
    {{range $value}}
    location {{.Path}} {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_pass http://{{.Vip}}:{{.Port}}
    }
    {{end}}
}
{{end}}
