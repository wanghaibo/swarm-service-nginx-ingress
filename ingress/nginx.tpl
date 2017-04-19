{{range .}}
server {
    server_name {{.VirtualHost}}
    listen 80;
    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_pass http://{{.Vip}}:{{.Port}}
    }
    
}
{{end}}

