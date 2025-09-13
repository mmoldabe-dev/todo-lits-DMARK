start:
	docker-compose up --build -d
	@echo "Starting Wails in dev mode..."
	@wails dev &
	@wait
