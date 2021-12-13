package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/JhyeonLee/BlockChain/blockchain"
	"github.com/JhyeonLee/BlockChain/utils"
	"github.com/gorilla/mux"
)

var port string

type url string // Stringer()

// type TextMarshaler, MUST correct name and correct signature
func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"` // struct field tag: usually json has lower case, but GO cannot export lower case var. struct filed tag help it to be lower case when it is json.
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"` // omitempty: hide field when field is empty
}

type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add a Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See a Block",
		},
	}

	// rw.Header().Add("Content-Type", "application/json") // make brower read it as json, not text
	// b, err := json.Marshal(data)                        // data interface to b json
	// utils.HandleErr(err)                                 // hadle err(template has helper func but json does not)
	// fmt.Fprintf(rw, "%s", b)
	json.NewEncoder(rw).Encode(data) // same as 3 lines above
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// rw.Header().Add("Content-Type", "application/json") // make brower read it as json, not text
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
	case "POST":
		var addBlockBody addBlockBody                                  // empty variable
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody)) // user request something to variable addBlcokBOdy
		blockchain.Blockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader((http.StatusCreated))
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	// http.HandlerFunc is not fucnction, it is type
	// adapter pattern
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()             // Router of Gorilla MUX
	router.Use(jsonContentTypeMiddleware) // middleware is just a function that is going to be called before the final destination.
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
