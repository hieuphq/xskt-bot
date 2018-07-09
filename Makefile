.PHONY: up build run deploy

up:
	docker-compose up -d

build:
	cd cmd/server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o api .

run:
	cd cmd/server && go run main.go

deploy:
	cd ansible && ANSIBLE_HOST_KEY_CHECKING=false ansible-playbook playbook.yml -e ansible_ssh_port=2222 