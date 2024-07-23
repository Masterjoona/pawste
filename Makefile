runnoweb:
	go run .

clean:
	rm $(binary_name)

rundocker:
	UID=${UID} GID=${GID} docker compose up -d --build