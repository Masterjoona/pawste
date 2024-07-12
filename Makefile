build:
	pnpm dlx golte && go build -o pawste "-ldflags=-s -w"

run:
	pnpm dlx golte dev && go run .

runnoweb:
	go run .

clean:
	rm $(binary_name)

rundocker:
	docker compose up -d --build