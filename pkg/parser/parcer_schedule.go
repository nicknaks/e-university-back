package parser

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

type Lesson struct {
	GroupID       string
	Day           string
	IsDenominator bool
	IsNumerator   bool
	Couple        string
	Name          string
	Type          string
	Teacher       string
	Place         string
	TeacherID     string
	SubjectID     string
}

// нужно сохранить в БД учителей, предметы и успеваемост

func parseScheduleForGroup(c *colly.Collector, groupID string, link string) ([]Lesson, error) {
	var lessons []Lesson

	c.OnHTML("div.row", func(element *colly.HTMLElement) {
		element.ForEach("div.col-lg-6", func(i int, element *colly.HTMLElement) {
			element.ForEach("table.table", func(i int, element *colly.HTMLElement) {
				element.ForEach("tbody", func(i int, element *colly.HTMLElement) {
					// день недели
					var day string
					element.ForEach("td[width]", func(i int, element *colly.HTMLElement) {
						element.ForEach("strong", func(i int, element *colly.HTMLElement) {
							day = element.Text
						})
					})

					element.ForEach("tr", func(i int, element *colly.HTMLElement) {
						var tmp Lesson
						tmp.GroupID = groupID
						tmp.Day = day
						// время
						element.ForEach("td.text-nowrap", func(i int, element *colly.HTMLElement) {
							tmp.Couple = element.Text
						})

						var lastElement string
						var lastCElement *colly.HTMLElement
						// предмет ЧС
						element.ForEach("td.text-info-bold", func(i int, element *colly.HTMLElement) {
							lastElement = strings.TrimSpace(element.Text)
							lastCElement = element
						})
						if lastElement != "" {
							lessonType, name, place, teachers := parseLesson(lastCElement)
							tmp.Type = lessonType
							tmp.Name = name
							tmp.Place = place
							tmp.Teacher = teachers
							tmp.IsNumerator = true
							lessons = append(lessons, tmp)
						}

						tmp.GroupID = groupID

						// предмет ЗН
						element.ForEach("td.text-primary", func(i int, element *colly.HTMLElement) {
							lastElement = strings.TrimSpace(element.Text)
							lastCElement = element
						})
						if lastElement == "" {
							tmp.Name = ""
							tmp.Type = ""
							tmp.Name = ""
							tmp.Place = ""
							tmp.Teacher = ""
							tmp.IsNumerator = false
						} else {
							lessonType, name, place, teachers := parseLesson(lastCElement)
							tmp.Name = lastElement
							tmp.Type = lessonType
							tmp.Name = name
							tmp.Place = place
							tmp.Teacher = teachers
							tmp.IsDenominator = true
							tmp.IsNumerator = false
						}
						// предмет общий
						element.ForEach("td[colspan]", func(i int, element *colly.HTMLElement) {
							lessonType, name, place, teachers := parseLesson(element)
							tmp.Type = lessonType
							tmp.Name = name
							tmp.Place = place
							tmp.Teacher = teachers
						})
						lessons = append(lessons, tmp)
					})
				})
			})
		})
	})

	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(fmt.Sprintf("%s%s", "https://lks.bmstu.ru", link))
	if err != nil {
		return nil, err
	}
	fmt.Println("done")

	var res []Lesson

	for _, lesson := range lessons {
		if lesson.Name != "" {
			lesson.GroupID = groupID
			res = append(res, lesson)
		}
	}

	for _, re := range res {
		fmt.Printf("день: %s время:%s название:%s_%s_%s_%s знаменатель:%v числитель:%v \n", re.Day, re.Couple, re.Type, re.Name, re.Place, re.Teacher, re.IsDenominator, re.IsNumerator)
	}

	return res, nil
}

// тип, предмет, кабинет, преподаватель
func parseLesson(element *colly.HTMLElement) (string, string, string, string) {
	var (
		lessonType string
		name       string
		place      string
		teachers   string
	)

	name = element.ChildText("span")
	element.ForEach("i", func(i int, element *colly.HTMLElement) {
		switch i {
		case 0:
			lessonType = element.Text
		case 1:
			place = element.Text
		case 2:
			teachers = element.Text
		default:
			panic(1)
		}
	})

	//if place == "" {
	//	fmt.Printf("%s_%s_%s_%s\n", lessonType, name, place, teachers)
	//}

	return lessonType, name, place, teachers
}

func ParseSchedule(ctx context.Context, data []Fac) ([]Lesson, error) {
	c := colly.NewCollector()

	var lessons []Lesson

	for _, fac := range data {
		for _, dep := range fac.Deps {
			for _, group := range dep.Groups {
				//if group.Number == "ИСОТ2-21А" {
				//	return lessons, nil
				//}

				fmt.Println("try for " + group.Number)
				fmt.Println(group.ID)
				tmp, err := parseScheduleForGroup(c, group.ID, group.Link)
				if err != nil {
					fmt.Printf("bad for %s with err %v\n", group.Number, err)
				} else {
					fmt.Println("good for " + group.Number)
				}
				lessons = append(lessons, tmp...)
			}
		}
	}

	fmt.Println(len(lessons))

	return lessons, nil
}
