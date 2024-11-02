.PHONY: tailwind-watch
tailwind-watch:
	tailwindcss -c tailwind.config.js -i ./static/css/input.css -o ./static/css/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	tailwindcss -c tailwind.config.js -i ./static/css/input.css -o ./static/css/style.css --minify

.PHONY: templ-generate
templ-generate:
	~/go/bin/templ generate

.PHONY: build
build:
	make templ-generate
	make tailwind-build
	@go build -o tmp/main main.go

# Live Reload
dev:
	air -c .air.toml
