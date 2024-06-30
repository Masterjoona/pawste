build:
	pnpm dlx golte && go build -o pawste "-ldflags=-s -w"

run:
	pnpm dlx golte dev && go run .

runnoweb:
	go run .

clean:
	rm $(binary_name)

newdb:
	pnpm dlx golte dev && PAWSTE_I_UNDERSTAND_THE_RISKS="true" go run .
