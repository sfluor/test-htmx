#!env sh

templ generate
npx tailwindcss -i ./static/input.css -o ./static/tailwind_out.css
go build -o ./test-htmx .
