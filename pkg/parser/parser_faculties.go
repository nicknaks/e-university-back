package parser

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type Fac struct {
	Name   string
	Number string
	Deps   []deps
}

type deps struct {
	Number string
	Groups []groups
}

type groups struct {
	ID     string
	Number string
	Link   string
}

func ParseFaculties(ctx context.Context) ([]Fac, error) {
	c := colly.NewCollector()

	var data []Fac

	c.OnHTML("body", func(element *colly.HTMLElement) {
		element.ForEach("div.list-group", func(main int, element *colly.HTMLElement) {
			if main >= 1 {
				return
			}
			// факультеты
			element.ForEach("a[href]", func(i int, element *colly.HTMLElement) {
				f := Fac{}
				// Номер
				element.ForEach("h4", func(i int, element *colly.HTMLElement) {
					f.Number = element.Text
				})
				// Название
				element.ForEach("p", func(i int, element *colly.HTMLElement) {
					f.Name = element.Text
				})
				if f.Name != "" {
					data = append(data, f)
				}
			})
			// Кафедры + группы
			element.ForEach("div.panel", func(mainI int, element *colly.HTMLElement) {
				element.ForEach("div.panel-body", func(i int, element *colly.HTMLElement) {
					element.ForEach("div.accordion", func(secondI int, element *colly.HTMLElement) {
						var temp deps
						// Кафедра
						element.ForEach("a.btn", func(i int, element *colly.HTMLElement) {
							element.ForEach("h4", func(i int, element *colly.HTMLElement) {
								temp.Number = element.Text
							})
						})
						data[mainI].Deps = append(data[mainI].Deps, temp)
						// Группы
						element.ForEach("div.row", func(i int, element *colly.HTMLElement) {
							element.ForEach("div.col-10", func(i int, element *colly.HTMLElement) {
								element.ForEach("div.row", func(i int, element *colly.HTMLElement) {
									element.ForEach("div.mb-1", func(i int, element *colly.HTMLElement) {
										element.ForEach("a[href]", func(i int, element *colly.HTMLElement) {
											data[mainI].Deps[secondI].Groups = append(data[mainI].Deps[secondI].Groups, groups{
												Number: strings.TrimSpace(element.Text),
												Link:   element.Attr("href"),
											})
										})
									})
								})
							})
						})
					})
				})
			})
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit("https://lks.bmstu.ru/schedule/list")
	fmt.Println("done")
	return data, err
}
