package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type DataItem struct {
	ID    int
	Name  string
	Value float64
}

type DataPage struct {
	Items       []DataItem
	PageNum     int
	TotalPages  int
	PrevPageNum int
	NextPageNum int
}

var data = func() []DataItem {
	items := make([]DataItem, 100) // Simulating 100 data items
	for i := 0; i < 100; i++ {
		items[i] = DataItem{ID: i + 1, Name: "Item " + strconv.Itoa(i+1), Value: float64(i+1) * 10}
	}
	return items
}()

func renderPage(w http.ResponseWriter, r *http.Request) {
	pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	itemsPerPage := 10
	start := (pageNum - 1) * itemsPerPage
	end := start + itemsPerPage
	if end > len(data) {
		end = len(data)
	}

	totalPages := (len(data) + itemsPerPage - 1) / itemsPerPage
	prevPage := pageNum - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := pageNum + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	page := DataPage{
		Items:       data[start:end],
		PageNum:     pageNum,
		TotalPages:  totalPages,
		PrevPageNum: prevPage,
		NextPageNum: nextPage,
	}

	tmpl, err := template.New("page.html").ParseFiles("page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/page", renderPage)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
