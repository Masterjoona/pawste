binary_name="pawste"

build:
	go build -o $(binary_name) "-ldflags=-s -w" *.go

run:
	go run *.go

clean:
	rm $(binary_name)

newdb:
	PAWSTE_I_UNDERSTAND_THE_RISKS="true" go run *.go
