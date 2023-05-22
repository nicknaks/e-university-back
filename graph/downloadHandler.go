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
	"net/http"
	"os"
)

// Handler — serving static files from cloud storage
type Handler struct {
	Storage store.Storage
}

func (h Handler) fillFile(ctx context.Context, file *os.File, subjectID string, filter string) (string, error) {
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

	studentsMap := map[string]int{}
	classesMap := map[string]int{}

	//
	f := excelize.NewFile()
	firstCell := "A1"
	f.SetCellStr("Sheet1", firstCell, "ФИО")
	for i, class := range classes {
		cell, err := excelize.CoordinatesToCellName(i+2, 1)
		if err != nil {
			panic(err)
		}
		err = f.SetCellStr("Sheet1", cell, fmt.Sprintf("%s %s %d", class.Day.Time.Format("02.01"), class.Name.String, class.Type))
		if err != nil {
			panic(err)
		}
		classesMap[class.ID] = i + 2
	}

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

	for _, progress := range classesProgresses {
		cell, err := excelize.CoordinatesToCellName(classesMap[progress.ClassID], studentsMap[progress.StudentID])
		if err != nil {
			panic(err)
		}
		err = f.SetCellFloat("Sheet1", cell, progress.Mark, 2, 64)
		if err != nil {
			panic(err)
		}
	}

	err = f.Write(file)
	if err != nil {
		panic(err)
	}

	return subjectID + "\n" + filter, nil
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

	f.Seek(0, 0)

	err = serve(req, w, f, filename)
	if err != nil {
		log.Fatal(err)
		w.Write([]byte(err.Error()))
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
