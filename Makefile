runnoweb:
	go run .

clean:
	rm $(binary_name)

rundocker:
	docker compose up -d --build