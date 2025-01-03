package main

import "gophkeeper/pkg/staticlint"

func main() {
	lint := staticlint.New()
	lint.Run()
}
