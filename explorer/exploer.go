package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/JhyeonLee/BlockChain/blockchain"
)

const (
	port        string = ":4000"
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

// rw : data to user
// r : request
func home(rw http.ResponseWriter, r *http.Request) {
	// template from html, not from text
	// must handle error, no exception in GO
	/*tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Fatal(err)
	}*/
	// template.MUST : deal with error, same action as above lines
	//tmpl := template.Must(template.ParseFiles("templates/home.gohtml"))
	//tmpl := template.Must(template.ParseFiles("templates/pages/home.gohtml"))

	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()                   // ParseForm populates r.Form and r.PostForm
		data := r.Form.Get("blockData") // Form : url.Values // Values map[string][]string // Must be same name "blockData" as name at <input>, add.gohtml
		blockchain.GetBlockchain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect) // when user add data, redirect to home(/) // redirect: refresh the site
	}
}

func Start() {
	// Go, Cannot use **/*.filetype
	// template.Must : helper function that deals with error // GO, Must deal with error and template.Must is doing it
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))     // Standard Library Package template, loading all gohtml templates using a pattern(ParseGlob)
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml")) // Variable templates, loading all gohtml templates using a pattern(ParseGlob)

	http.HandleFunc("/", home)   // register function home for the given patern "/"
	http.HandleFunc("/add", add) // register function add for the given patern "/add"

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
