.PHONY: dev vite templ air build run clean

dev:
	@make -j3 vite templ air

vite:
	@cd web/app && pnpm dev

templ:
	@templ generate -watch -path ./web/template -proxy="http://localhost:3011" -proxyport 3000

air:
	@air

build:
	@templ generate -path ./web/template
	@cd web/app && pnpm build
	@go build -o ./bin/app ./cmd/app

run:
	@ENV=prod ./bin/app

clean:
	@rm -rf ./tmp ./bin ./web/static/* .reload-trigger

