package graph

import (
	"back/graph/model"
	"back/internal/store"
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/vfaronov/httpheader"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

// Handler — serving static files from cloud storage
type Handler struct {
	Storage store.Storage
}

func (h Handler) fillFile(ctx context.Context, file *os.File, subjectID string, filter string) (string, error) {
	subjects, err := h.Storage.ListSubjects(ctx, &model.SubjectsFilter{
		ID: []string{subjectID},
	})

	subject := subjects[0]

	classes, err := h.Storage.ListClasses(ctx, &model.ClassesFilter{
		SubjectID: lo.ToPtr(subjectID),
	})
	if err != nil {
		return "", err
	}

	students, err := h.Storage.ListStudents(ctx, &model.StudentsFilter{
		GroupID: lo.ToPtr(classes[0].GroupID),
	})

	var classessIDs []string
	for _, class := range classes {
		classessIDs = append(classessIDs, class.ID)
	}

	classesProgresses, err := h.Storage.ListClassesProgresses(ctx, &model.ClassesProgressFilter{
		ClassIDIn: classessIDs,
	})

	subjectsResults, err := h.Storage.ListSubjectsResults(ctx, &model.SubjectResultsFilter{
		SubjectID: lo.ToPtr(subjectID),
	})

	studentsMap := map[string]int{}
	classesMap := map[string]int{}
	moduleCols := make([]int, 6)
	f := excelize.NewFile()

	moduleStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"588DBD"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		panic(err)
	}

	moduleMainStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"588DBD"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{{
			Type:  "bottom",
			Color: "000000",
			Style: 2,
		}, {
			Type:  "right",
			Color: "000000",
			Style: 2,
		}, {
			Type:  "left",
			Color: "000000",
			Style: 2,
		}},
	})
	if err != nil {
		panic(err)
	}

	firstBorderStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{{
			Type:  "bottom",
			Color: "000000",
			Style: 2,
		}, {
			Type:  "right",
			Color: "000000",
			Style: 2,
		}},
	})

	titleBorderStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{{
			Type:  "bottom",
			Color: "000000",
			Style: 2,
		}},
	})
	if err != nil {
		panic(err)
	}

	titleBorderStyle2, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{{
			Type:  "right",
			Color: "000000",
			Style: 2,
		}},
	})
	if err != nil {
		panic(err)
	}

	//classStyle, err := f.NewStyle(&excelize.Style{
	//	Alignment: &excelize.Alignment{
	//		Horizontal: "center",
	//		Vertical:   "center",
	//		WrapText:   true,
	//	},
	//	Border: []excelize.Border{{
	//		Type:  "bottom",
	//		Color: "000000",
	//		Style: 2,
	//	}},
	//})

	err = f.SetRowStyle("Sheet1", 1, 1, titleBorderStyle)
	if err != nil {
		panic(err)
	}

	err = f.SetColStyle("Sheet1", "A", titleBorderStyle2)
	if err != nil {
		panic(err)
	}

	// заполнение фио
	firstCell := "A1"
	err = f.SetCellStr("Sheet1", firstCell, "ФИО")
	f.SetCellStyle("Sheet1", firstCell, firstCell, firstBorderStyle)

	err = f.SetColWidth("Sheet1", "A", "A", 20)
	if err != nil {
		panic(err)
	}

	err = f.SetRowHeight("Sheet1", 1, 45)
	if err != nil {
		panic(err)
	}

	// заполнение студентов
	for i, student := range students {
		cell, err := excelize.CoordinatesToCellName(1, i+2)
		if err != nil {
			panic(err)
		}
		err = f.SetCellStr("Sheet1", cell, fmt.Sprintf("%s", student.Name))
		if err != nil {
			panic(err)
		}
		studentsMap[student.ID] = i + 2
	}

	lastClassIndex := 2
	// заполнение классов 1 модуля
	for _, class := range classes {
		if class.Module != 1 {
			continue
		}
		cell, err := excelize.CoordinatesToCellName(lastClassIndex, 1)
		if err != nil {
			panic(err)
		}
		err = f.SetCellStr("Sheet1", cell, fmt.Sprintf("%s\n%s\n%s", parseClassType(class.Type), class.Name.String, class.Day.Time.Format("02.01")))
		if err != nil {
			panic(err)
		}
		classesMap[class.ID] = lastClassIndex
		lastClassIndex++
	}

	// модуль 1
	cell, err := excelize.CoordinatesToCellName(lastClassIndex, 1)
	if err != nil {
		panic(err)
	}
	moduleCols[0] = lastClassIndex
	err = f.SetCellStr("Sheet1", cell, fmt.Sprintf("Модуль 1"))
	if err != nil {
		panic(err)
	}

	err = f.SetCellStyle("Sheet1", cell, cell, moduleMainStyle)
	if err != nil {
		panic(err)
	}
	for _, result := range subjectsResults {
		fmt.Println(result)
		cell, err := excelize.CoordinatesToCellName(lastClassIndex, studentsMap[result.StudentID])
		if err != nil {
			panic(err)
		}
		fmt.Println(result.FirstModuleMark)
		err = f.SetCellFloat("Sheet1", cell, result.FirstModuleMark, 2, 64)
		if err != nil {
			panic(err)
		}

		err = f.SetCellStyle("Sheet1", cell, cell, moduleStyle)
		if err != nil {
			panic(err)
		}

		fmt.Println(err)
	}

	lastClassIndex++
	// заполнение классов 2 модуля
	for _, class := range classes {
		if class.Module != 2 {
			continue
		}
		cell, err := excelize.CoordinatesToCellName(lastClassIndex, 1)
		if err != nil {
			panic(err)
		}
		err = f.SetCellStr("Sheet1", cell, fmt.Sprintf("%s\n%s\n%s", parseClassType(class.Type), class.Name.String, class.Day.Time.Format("02.01")))
		if err != nil {
			panic(err)
		}
		classesMap[class.ID] = lastClassIndex
		lastClassIndex++
	}

	// модуль 2
	cell, err = excelize.CoordinatesToCellName(lastClassIndex, 1)
	if err != nil {
		panic(err)
	}
	moduleCols[1] = lastClassIndex
	err = f.SetCellStr("Sheet1", cell, fmt.Sprintf("Модуль 2"))
	if err != nil {
		panic(err)
	}

	err = f.SetCellStyle("Sheet1", cell, cell, moduleMainStyle)
	if err != nil {
		panic(err)
	}

	for _, result := range subjectsResults {
		fmt.Println(result)
		cell, err := excelize.CoordinatesToCellName(lastClassIndex, studentsMap[result.StudentID])
		if err != nil {
			panic(err)
		}
		fmt.Println(cell)
		err = f.SetCellStyle("Sheet1", cell, cell, moduleStyle)
		if err != nil {
			panic(err)
		}

		err = f.SetCellFloat("Sheet1", cell, result.SecondModuleMark, 2, 64)
		if err != nil {
			panic(err)
		}
	}
	lastClassIndex++

	// заполнение классов 3 модуля
	for _, class := range classes {
		if class.Module != 3 {
			continue
		}
		cell, err := excelize.CoordinatesToCellName(lastClassIndex, 1)
		if err != nil {
			panic(err)
		}
		err = f.SetCellStr("Sheet1", cell, fmt.Sprintf("%s\n%s\n%s", parseClassType(class.Type), class.Name.String, class.Day.Time.Format("02.01")))
		if err != nil {
			panic(err)
		}
		classesMap[class.ID] = lastClassIndex
		lastClassIndex++
	}

	// модуль 3
	cell, err = excelize.CoordinatesToCellName(lastClassIndex, 1)
	if err != nil {
		panic(err)
	}
	moduleCols[2] = lastClassIndex
	err = f.SetCellStr("Sheet1", cell, fmt.Sprintf("Модуль 3"))
	if err != nil {
		panic(err)
	}

	err = f.SetCellStyle("Sheet1", cell, cell, moduleMainStyle)
	if err != nil {
		panic(err)
	}

	for _, result := range subjectsResults {
		fmt.Println(result)
		cell, err := excelize.CoordinatesToCellName(lastClassIndex, studentsMap[result.StudentID])
		if err != nil {
			panic(err)
		}
		fmt.Println(cell)
		err = f.SetCellStyle("Sheet1", cell, cell, moduleStyle)
		if err != nil {
			panic(err)
		}

		err = f.SetCellFloat("Sheet1", cell, result.ThirdModuleMark, 2, 64)
		if err != nil {
			panic(err)
		}

	}
	lastClassIndex++

	if subject.Type == 2 {
		// экзамен
		moduleCols[3] = lastClassIndex
		cell, err = excelize.CoordinatesToCellName(lastClassIndex, 1)
		if err != nil {
			panic(err)
		}
		err = f.SetCellStr("Sheet1", cell, fmt.Sprintf("Экзамен"))
		if err != nil {
			panic(err)
		}

		for _, result := range subjectsResults {
			cell, err := excelize.CoordinatesToCellName(lastClassIndex, studentsMap[result.StudentID])
			if err != nil {
				panic(err)
			}
			err = f.SetCellInt("Sheet1", cell, result.ExamResult)
			if err != nil {
				panic(err)
			}
		}
		lastClassIndex++
	}

	// результат
	cell, err = excelize.CoordinatesToCellName(lastClassIndex, 1)
	if err != nil {
		panic(err)
	}
	err = f.SetCellStr("Sheet1", cell, fmt.Sprintf("Результат"))
	if err != nil {
		panic(err)
	}

	for _, result := range subjectsResults {
		fmt.Println(result)
		cell, err := excelize.CoordinatesToCellName(lastClassIndex, studentsMap[result.StudentID])
		if err != nil {
			panic(err)
		}
		fmt.Println(cell)
		err = f.SetCellFloat("Sheet1", cell, math.Floor(result.Mark), 2, 64)
		if err != nil {
			panic(err)
		}
	}
	lastClassIndex++

	// Сумма баллов за семинары
	cell, err = excelize.CoordinatesToCellName(lastClassIndex, 1)
	if err != nil {
		panic(err)
	}
	err = f.SetCellStr("Sheet1", cell, fmt.Sprintf("Сумма баллов за семинары"))
	if err != nil {
		panic(err)
	}

	for _, student := range students {
		var count float64
		for _, class := range classes {
			if class.Type != 3 {
				continue
			}
			for _, progress := range classesProgresses {
				if progress.ClassID == class.ID && progress.StudentID == student.ID {
					count += progress.Mark
				}
			}
		}
		cell, err = excelize.CoordinatesToCellName(lastClassIndex, studentsMap[student.ID])
		err = f.SetCellFloat("Sheet1", cell, count, 2, 64)
	}
	lastClassIndex++

	// Сумма баллов за ЛР
	cell, err = excelize.CoordinatesToCellName(lastClassIndex, 1)
	if err != nil {
		panic(err)
	}
	f.SetCellStr("Sheet1", cell, fmt.Sprintf("Сумма баллов за ЛР"))
	for _, student := range students {
		var count float64
		for _, class := range classes {
			if class.Type != 1 {
				continue
			}
			for _, progress := range classesProgresses {
				if progress.ClassID == class.ID && progress.StudentID == student.ID {
					count += progress.Mark
				}
			}
		}
		cell, err = excelize.CoordinatesToCellName(lastClassIndex, studentsMap[student.ID])
		err = f.SetCellFloat("Sheet1", cell, count, 2, 64)
	}
	lastClassIndex++

	// Количество пропусков лекций
	cell, err = excelize.CoordinatesToCellName(lastClassIndex, 1)
	moduleCols[5] = lastClassIndex
	if err != nil {
		panic(err)
	}
	f.SetCellStr("Sheet1", cell, fmt.Sprintf("Кол-во пропусков лекций"))
	for _, student := range students {
		var count int
		for _, class := range classes {
			if class.Type != 2 {
				continue
			}
			for _, progress := range classesProgresses {
				if progress.ClassID == class.ID && progress.StudentID == student.ID {
					if progress.IsAbsent == true {
						count++
					}
				}
			}
		}
		cell, err = excelize.CoordinatesToCellName(lastClassIndex, studentsMap[student.ID])
		err = f.SetCellInt("Sheet1", cell, count)
	}
	lastClassIndex++

	// % пропусков на сегодня
	moduleCols[4] = lastClassIndex
	cell, err = excelize.CoordinatesToCellName(lastClassIndex, 1)
	if err != nil {
		panic(err)
	}
	f.SetCellStr("Sheet1", cell, "% пропусков на сегодня")
	today := time.Now()
	for _, student := range students {
		var count float64
		var totalCount float64
		for _, class := range classes {
			// если занятие после сегодня
			if class.Day.Time.After(today) {
				continue
			}
			totalCount++
			for _, progress := range classesProgresses {
				if progress.ClassID == class.ID && progress.StudentID == student.ID {
					if progress.IsAbsent == true {
						count++
					}
				}
			}
		}
		cell, err = excelize.CoordinatesToCellName(lastClassIndex, studentsMap[student.ID])
		err = f.SetCellInt("Sheet1", cell, int(math.Floor(count/totalCount*100)))
	}
	lastClassIndex++

	// Итог
	cell, err = excelize.CoordinatesToCellName(lastClassIndex, 1)
	if err != nil {
		panic(err)
	}
	f.SetCellStr("Sheet1", cell, "Итог")
	switch subject.Type {
	case 1: // зачет
		for _, result := range subjectsResults {
			mark := "Зачет"
			if math.Floor(result.FirstModuleMark) < 18 {
				mark = "Незачет"
			}
			if math.Floor(result.SecondModuleMark) < 18 {
				mark = "Незачет"
			}
			if math.Floor(result.ThirdModuleMark) < 18 {
				mark = "Незачет"
			}
			if math.Floor(result.Mark) < 60 {
				mark = "Незачет"
			}

			cell, err := excelize.CoordinatesToCellName(lastClassIndex, studentsMap[result.StudentID])
			if err != nil {
				panic(err)
			}
			err = f.SetCellStr("Sheet1", cell, mark)
		}
	case 2: // экз
		for _, result := range subjectsResults {
			mark := 2
			if math.Floor(result.FirstModuleMark) < 18 {
				mark = 2
			}
			if math.Floor(result.SecondModuleMark) < 18 {
				mark = 2
			}
			if math.Floor(result.ThirdModuleMark) < 18 {
				mark = 2
			}
			if result.ExamResult < 18 {
				mark = 2
			}

			if math.Floor(result.Mark) < 60 {
				mark = 2
			} else if math.Floor(result.Mark) < 71 {
				mark = 3
			} else if math.Floor(result.Mark) < 85 {
				mark = 4
			} else {
				mark = 5
			}

			cell, err := excelize.CoordinatesToCellName(lastClassIndex, studentsMap[result.StudentID])
			if err != nil {
				panic(err)
			}
			err = f.SetCellInt("Sheet1", cell, mark)
		}
	}
	lastClassIndex++

	// заполнение прогресса
	for _, progress := range classesProgresses {
		cell, err := excelize.CoordinatesToCellName(classesMap[progress.ClassID], studentsMap[progress.StudentID])
		if err != nil {
			panic(err)
		}
		err = f.SetCellFloat("Sheet1", cell, progress.Mark, 2, 64)
		if progress.IsAbsent {
			style, err := f.NewStyle(&excelize.Style{
				Fill: excelize.Fill{Type: "pattern", Color: []string{"BCBCBC"}, Pattern: 1},
			})
			if err != nil {
				fmt.Println(err)
			}

			f.SetCellStyle("Sheet1", cell, cell, style)
			err = f.SetCellFloat("Sheet1", cell, progress.Mark, 2, 64)
		}
		if err != nil {
			panic(err)
		}
	}

	// Создание доп таблиц
	index1, err := f.NewSheet("Модуль 1")
	index2, err := f.NewSheet("Модуль 2")
	index3, err := f.NewSheet("Модуль 3")
	index4, err := f.NewSheet("Итог")
	mainIndex, err := f.GetSheetIndex("Sheet1")
	err = f.CopySheet(mainIndex, index1)
	if err != nil {
		panic(err)
	}
	err = f.CopySheet(mainIndex, index2)
	if err != nil {
		panic(err)
	}
	err = f.CopySheet(mainIndex, index3)
	if err != nil {
		panic(err)
	}
	err = f.CopySheet(mainIndex, index4)
	if err != nil {
		panic(err)
	}

	module1Last, err := excelize.ColumnNumberToName(moduleCols[0] - 1)
	if err != nil {
		panic(err)
	}
	module2Start, err := excelize.ColumnNumberToName(moduleCols[0] + 1)
	if err != nil {
		panic(err)
	}
	module2Last, err := excelize.ColumnNumberToName(moduleCols[1] - 1)
	if err != nil {
		panic(err)
	}
	module3Start, err := excelize.ColumnNumberToName(moduleCols[1] + 1)
	if err != nil {
		panic(err)
	}
	module3Last, err := excelize.ColumnNumberToName(moduleCols[2] - 1)
	if err != nil {
		panic(err)
	}
	err = f.SetColVisible("Модуль 1", fmt.Sprintf("%s:%s", module2Start, module2Last), false)
	if err != nil {
		panic(err)
	}

	err = f.SetColVisible("Модуль 1", fmt.Sprintf("%s:%s", module3Start, module3Last), false)
	if err != nil {
		panic(err)
	}

	err = f.SetColVisible("Модуль 2", fmt.Sprintf("%s:%s", "B", module1Last), false)
	if err != nil {
		panic(err)
	}

	err = f.SetColVisible("Модуль 2", fmt.Sprintf("%s:%s", module3Start, module3Last), false)
	if err != nil {
		panic(err)
	}

	err = f.SetColVisible("Модуль 3", fmt.Sprintf("%s:%s", "B", module1Last), false)
	if err != nil {
		panic(err)
	}

	err = f.SetColVisible("Модуль 3", fmt.Sprintf("%s:%s", module2Start, module2Last), false)
	if err != nil {
		panic(err)
	}

	err = f.SetColVisible("Итог", fmt.Sprintf("%s:%s", "B", module1Last), false)
	if err != nil {
		panic(err)
	}

	err = f.SetColVisible("Итог", fmt.Sprintf("%s:%s", module2Start, module2Last), false)
	if err != nil {
		panic(err)
	}

	err = f.SetColVisible("Итог", fmt.Sprintf("%s:%s", module3Start, module3Last), false)
	if err != nil {
		panic(err)
	}

	str1, err := excelize.CoordinatesToCellName(1, studentsMap[students[0].ID], true)
	str2, err := excelize.CoordinatesToCellName(1, studentsMap[students[len(students)-1].ID], true)
	categoriesStudent := fmt.Sprintf("Sheet1!%s:%s", str1, str2)

	firstModuleCol, err := excelize.ColumnNumberToName(moduleCols[0])
	secondModuleCol, err := excelize.ColumnNumberToName(moduleCols[1])
	thirdModuleCol, err := excelize.ColumnNumberToName(moduleCols[2])
	examModuleCol, err := excelize.ColumnNumberToName(moduleCols[3])
	apsentCol, err := excelize.ColumnNumberToName(moduleCols[4])
	apsentLecCol, err := excelize.ColumnNumberToName(moduleCols[5])

	firstValue := fmt.Sprintf("Sheet1!$%s$%d:$%s$%d", firstModuleCol, studentsMap[students[0].ID], firstModuleCol, studentsMap[students[len(students)-1].ID])
	secondValue := fmt.Sprintf("Sheet1!$%s$%d:$%s$%d", secondModuleCol, studentsMap[students[0].ID], secondModuleCol, studentsMap[students[len(students)-1].ID])
	thirdValue := fmt.Sprintf("Sheet1!$%s$%d:$%s$%d", thirdModuleCol, studentsMap[students[0].ID], thirdModuleCol, studentsMap[students[len(students)-1].ID])
	fourValue := fmt.Sprintf("Sheet1!$%s$%d:$%s$%d", examModuleCol, studentsMap[students[0].ID], examModuleCol, studentsMap[students[len(students)-1].ID])
	fiveValue := fmt.Sprintf("Sheet1!$%s$%d:$%s$%d", apsentCol, studentsMap[students[0].ID], apsentCol, studentsMap[students[len(students)-1].ID])
	sixValue := fmt.Sprintf("Sheet1!$%s$%d:$%s$%d", apsentLecCol, studentsMap[students[0].ID], apsentLecCol, studentsMap[students[len(students)-1].ID])

	series := []excelize.ChartSeries{
		{
			Name:       fmt.Sprintf("Sheet1!$%s$1", firstModuleCol),
			Categories: categoriesStudent,
			Values:     firstValue,
		},
		{
			Name:       fmt.Sprintf("Sheet1!$%s$1", secondModuleCol),
			Categories: categoriesStudent,
			Values:     secondValue,
		},
		{
			Name:       fmt.Sprintf("Sheet1!$%s$1", thirdModuleCol),
			Categories: categoriesStudent,
			Values:     thirdValue,
		},
	}

	if subject.Type == 2 {
		series = append(series, excelize.ChartSeries{
			Name:       fmt.Sprintf("Sheet1!$%s$1", examModuleCol),
			Categories: categoriesStudent,
			Values:     fourValue,
		})
	}

	err = f.AddChart("Sheet1", "B20", &excelize.Chart{
		Type:   excelize.ColStacked,
		Series: series,
		Format: excelize.GraphicOptions{
			OffsetX: 15,
			OffsetY: 10,
		},
		Legend: excelize.ChartLegend{
			Position: "left",
		},
		Title: excelize.ChartTitle{
			Name: "Результаты студентов",
		},
		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     true,
			ShowVal:         true,
		},
		Dimension: excelize.ChartDimension{
			Width:  20 * 54,
			Height: 15 * 15,
		},
		ShowBlanksAs: "gap",
	})
	if err != nil {
		panic(err)
	}

	err = f.AddChart("Sheet1", "B40", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			{
				Name:       fmt.Sprintf("Sheet1!$%s$1", apsentCol),
				Categories: categoriesStudent,
				Values:     fiveValue,
			},
			{
				Name:       fmt.Sprintf("Sheet1!$%s$1", apsentLecCol),
				Categories: categoriesStudent,
				Values:     sixValue,
			},
		},
		Format: excelize.GraphicOptions{
			OffsetX: 15,
			OffsetY: 10,
		},
		Dimension: excelize.ChartDimension{
			Width:  20 * 54,
			Height: 15 * 15,
		},
		Legend: excelize.ChartLegend{
			Position: "left",
		},
		Title: excelize.ChartTitle{
			Name: "Количество пропусков",
		},
		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     true,
			ShowVal:         true,
		},
		ShowBlanksAs: "gap",
	})
	if err != nil {
		panic(err)
	}

	err = f.Write(file)
	if err != nil {
		panic(err)
	}

	return subjectID + "\n" + filter, nil
}

func parseClassType(clasType int) string {
	switch clasType {
	case 1:
		return "лаб"
	case 2:
		return "лек"
	case 3:
		return "сем"
	}
	return ""
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	f, err := os.CreateTemp("", "example.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())

	filename, err := h.fillFile(req.Context(), f, req.URL.Query().Get("id"), req.URL.Query().Get("filter"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	err = serve(req, w, f, filename)
	if err != nil {
		log.Fatal(err)
	}

	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
}

const (
	kib        = 1 << 10
	chunkSize  = kib * 200
	inline     = "inline"
	attachment = "attachment"
)

func disposition(req *http.Request) string {
	switch req.URL.Query().Get("inline") {
	case "1", "true":
		return inline
	default:
		return attachment
	}
}

func serve(req *http.Request, w http.ResponseWriter, file *os.File, fileName string) error {
	flusher := w.(http.Flusher)

	w.Header().Set("access-control-allow-origin", "*")
	httpheader.SetContentType(w.Header(), "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", map[string]string{})
	httpheader.SetContentDisposition(w.Header(), disposition(req), fileName, map[string]string{})

	for {
		buf := make([]byte, chunkSize)
		n, err := file.Read(buf)
		buf = buf[:n]

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		_, _ = w.Write(buf)
		flusher.Flush()
	}

	return nil
}
