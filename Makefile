make:
	go build .
	mv translator VMTranslator
	chmod +x ./VMTranslator