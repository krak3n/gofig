ifndef GODOCPORT
godoc: GODOCPORT = 8080
endif
godoc:
	godoc -http=:$(GODOCPORT)
