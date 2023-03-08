make:
	go build .
	mv vmtranslator VMTranslator
	chmod +x ./VMTranslator