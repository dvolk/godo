package main

import (
	"fmt"
	"html/template"
)

func Icon(s string) template.HTML {
	return template.HTML(fmt.Sprintf(`<i class="fa fa-fw fa-%s"></i>`, s))
}
