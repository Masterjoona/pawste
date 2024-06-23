binary_name="pawste"

build:
	pnpm dlx golte dev && go build -o $(binary_name) "-ldflags=-s -w"

run:
	pnpm dlx golte dev && go run .

clean:
	rm $(binary_name)

newdb:
	PAWSTE_I_UNDERSTAND_THE_RISKS="true" pnpm dlx golte dev && go run .
