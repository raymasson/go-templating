package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type friend struct {
	Fname string
}

type user struct {
	FirstName string
	LastName  string
	Age       int
	Emails    []string
	Friends   []friend
}

func main() {
	// Simple template
	simpleTemplate()

	// Nested fields template
	nestedFieldsTemplate()

	// Conditions template
	conditionsTemplate()

	// Variables template
	variablesTemplate()

	// Templates validation
	templatesValidation()

	// Sub-Templates
	subTemplates()

	// Web page
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", index).Methods("GET")
	fmt.Println("Web page accessible at localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

// EmailDealWith ...
func EmailDealWith(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	// replace "at" to @ symbol
	return strings.Replace(s, "at", "@", 1)
}

func index(w http.ResponseWriter, r *http.Request) {
	t := template.New("welcome template") // Create a template.
	t, _ = t.ParseFiles("welcome.html")   // Parse template file.
	u := user{
		FirstName: "Ray",
		LastName:  "MASSON",
		Age:       31,
	}
	t.ExecuteTemplate(w, "welcome.html", u) // merge.
}

func simpleTemplate() {
	fmt.Println("Simple template:")
	t := template.New("FirstName example")
	t, _ = t.Parse("hello {{.FirstName}}!")
	u := user{FirstName: "Ray"}
	t.Execute(os.Stdout, u)
	fmt.Println()
	fmt.Println()
}

func nestedFieldsTemplate() {
	fmt.Println("Nested fields template:")
	f1 := friend{Fname: "Laurent"}
	f2 := friend{Fname: "Julien"}
	t := template.New("friends example")
	t = t.Funcs(template.FuncMap{"emailDeal": EmailDealWith})
	t, _ = t.Parse(`hello {{.FirstName}}!
            {{range .Emails}}
                an email {{.|emailDeal}}
            {{end}}
            {{with .Friends}}
            {{range .}}
                my friend name is {{.Fname}}
            {{end}}
            {{end}}
            `)
	u := user{FirstName: "Ray",
		Emails:  []string{"rmatbeego.me", "raymassatgmail.com"},
		Friends: []friend{f1, f2}}
	t.Execute(os.Stdout, u)
	fmt.Println()
	fmt.Println()
}

func conditionsTemplate() {
	fmt.Println("Conditions template:")
	tEmpty := template.New("template tEmpty")
	tEmpty = template.Must(tEmpty.Parse("Empty pipeline if demo: {{if ``}} will not be outputted. {{end}}\n"))
	tEmpty.Execute(os.Stdout, nil)

	tWithValue := template.New("template tWithValue")
	tWithValue = template.Must(tWithValue.Parse("Not empty pipeline if demo: {{if `anything`}} will be outputted. {{end}}\n"))
	tWithValue.Execute(os.Stdout, nil)

	tIfElse := template.New("template tIfElse")
	tIfElse = template.Must(tIfElse.Parse("if-else demo: {{if `anything`}} if part {{else}} else part.{{end}}\n"))
	tIfElse.Execute(os.Stdout, nil)

	fmt.Println()
	fmt.Println()
}

func variablesTemplate() {
	fmt.Println("Variables template:")
	t := template.New("Variables example")
	t = template.Must(t.Parse(`{{with $x := "output" | printf "%q"}}{{$x}}{{end}}`))
	t.Execute(os.Stdout, nil)

	t2 := template.New("Variables example 2")
	t2 = template.Must(t2.Parse(`{{with $x := "output"}}{{printf "%q" $x}}{{end}}`))
	t2.Execute(os.Stdout, nil)

	t3 := template.New("Variables example 3")
	t3 = template.Must(t3.Parse(`{{with $x := "output"}}{{$x | printf "%q"}}{{end}}`))
	t3.Execute(os.Stdout, nil)
	fmt.Println()
	fmt.Println()
}

func templatesValidation() {
	fmt.Println("Templates validation:")
	tOk := template.New("first")
	template.Must(tOk.Parse(" some static text /* and a comment */"))
	fmt.Println("The first one parsed OK.")

	template.Must(template.New("second").Parse("some static text {{ .Name }}"))
	fmt.Println("The second one parsed OK.")

	/*fmt.Println("The next one ought to fail.")
	tErr := template.New("check parse error with Must")
	template.Must(tErr.Parse(" some static text {{ .Name }"))*/
	fmt.Println()
	fmt.Println()
}

func subTemplates() {
	var allFiles []string
	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".tmpl") {
			allFiles = append(allFiles, "./templates/"+filename)
		}
	}

	templates, err := template.ParseFiles(allFiles...)

	s1 := templates.Lookup("header.tmpl")
	s1.ExecuteTemplate(os.Stdout, "header", nil)
	fmt.Println()
	s2 := templates.Lookup("content.tmpl")
	s2.ExecuteTemplate(os.Stdout, "content", nil)
	fmt.Println()
	s3 := templates.Lookup("footer.tmpl")
	s3.ExecuteTemplate(os.Stdout, "footer", nil)
	fmt.Println()
	s3.Execute(os.Stdout, nil)
}
