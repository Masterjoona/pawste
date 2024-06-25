binary_name="pawste"

build:
	npx golte dev --use-pnpm && go build -o $(binary_name) "-ldflags=-s -w"

run:
	npx golte dev --use-pnpm && go run .

runnoweb:
	go run .

clean:
	rm $(binary_name)

newdb:
	npx golte dev --use-pnpm && PAWSTE_I_UNDERSTAND_THE_RISKS="true" go run .
