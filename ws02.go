// ws02
//fmt.Println(strings.Join(reg[:],","))
package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

type Userid struct {
	ulogin   string
	upass    string
	Userid   string
	usession string
}

var (
	// компилируем шаблоны, если не удалось, то выходим
	post_template = template.Must(template.ParseFiles(path.Join("public", "index.html")))
)

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title
	if strings.Index(title, ".html") == -1 {
		filename = title + ".txt"
	} //else {
	//filename = title
	//}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func emptHandler(w http.ResponseWriter, r *http.Request, title string) {
	// обработчик запросов
	if err := post_template.ExecuteTemplate(w, "index.html", nil); err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
	/*
		fmt.Print("Salam\n", templates.DefinedTemplates())
		err := templates.ExecuteTemplate(w, "index.html", "yxaxa")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	*/
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {

	if strings.Index(tmpl, ".html") == -1 {
		tmpl = tmpl + ".html"
	}

	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "public/index.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		geUrls := r.URL.Path
		fmt.Print("request: ", r.URL.Path+"\n")

		if geUrls == "/" || geUrls == "" {
			fmt.Print("пустой", "\n")
			//			fn(w, r, "")
			emptHandler(w, r, "")
			return
		} else {
			m := validPath.FindStringSubmatch(geUrls)
			if m == nil {
				fmt.Print("таки пустой", "\n")
				http.NotFound(w, r)
				return
			}
			fn(w, r, m[2])
		}
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/", makeHandler(emptHandler))

	http.ListenAndServe(":7777", nil)
	/*
		p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
		p1.save()
		p2, _ := loadPage("TestPage")
		fmt.Println(string(p2.Body))
	*/
}
