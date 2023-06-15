package main

import "regexp"

func main() {

	input := "/bzz/rrrr/56E2E4D5-96FE-4815-BE96-25DCCECE982F/"

	re := regexp.MustCompile(`/([^/]*)/`)

	matches := re.FindStringSubmatch(input)

	if len(matches) > 1 {
		println(matches[1])
	}

	println(input[len(matches[0])-1:])

}
