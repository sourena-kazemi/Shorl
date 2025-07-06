templ:
	@templ generate --watch
tailwind:
	@tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
air:
	@air