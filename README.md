# go-additional-json

## Example
```go
package main

import (
	"fmt"
	"github.com/leberKleber/go-additional-json"
	"os"
)

func main() {
	s := struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description struct {
			Short string `json:"short"`
			Long  string `json:"long"`
		} `json:"description"`
		Other map[string]string      `json:"-" aj:"other"`
	}{}

	j := `{
		"id": 54321,
		"name": "NaMe",
		"description": {
			"short": "short",
			"long": "lloonngg"
		},
		"other1": "ootthheerr11",
		"other2": "ootthheerr22"
	}`

	err := additionaljson.DefaultUnmarshaler.Unmarshal([]byte(j), &s)
	if err != nil {
		fmt.Printf("a error ocurred: %s", err)
		os.Exit(1)
	}

	fmt.Printf("unmarshaled:\n%v", s)
}
```

Output: 
```text
unmarshaled:
{54321 NaMe {short lloonngg} map[other1:ootthheerr11 other2:ootthheerr22]}
```