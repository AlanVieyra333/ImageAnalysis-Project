/**
 * Created by Fenix on 03/10/2016.
set GOPATH=%cd%
go build imageEdit
go install server
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"imageEdit"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

//var templates = template.Must(template.ParseFiles("index.html"))

//var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(upload|json)*(/([a-zA-Z0-9]*))?$")
var validPathImage = regexp.MustCompile("^(https?://)?" + "([a-zA-Z0-9]*.?)+" + "(:[0-9]*)?" + "/uploaded/([a-zA-Z0-9]*.[a-zA-Z0-9]*)" + "$")

/*	---------------------------*	Functions	---------------------------*	*/

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

/*/	---
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}//*/

/*func renderTemplate(w http.ResponseWriter, title string) {
	err := templates.ExecuteTemplate(w, title+".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}//*/

/*	--------------------------->	Handlers	--------------------------->	*/

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		//fmt.Println("" + m[0] + " " + m[1] + " " + m[2]);
		fn(w, r, "index")
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
	/*/fmt.Println(title);
	p, err := loadPage("index")
	if err != nil {
		http.NotFound(w, r)
		return
	}*/
	//renderTemplate(w, title)
	indexTmpl := template.New("index.html").Delims("<<", ">>")
	indexTmpl, _ = indexTmpl.ParseFiles("index.html")
	indexTmpl.Execute(w, nil)
}

/*	------------------------->	Functions operate files	------------------------->	*/

type Status struct {
	Code         int
	Message      string
	FileName     string
	FileNameEdit string
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func uploadHandler(w http.ResponseWriter, req *http.Request, title string) {
	//	fmt.Println("Upl")
	confirm := &Status{
		Code:         1,
		Message:      "Archivo subido exitosamente.",
		FileName:     "",
		FileNameEdit: "",
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		confirm.Code = 2
		confirm.Message = err.Error()
	}
	//	fmt.Println("file: " + confirm.Message);
	confirm.FileName = handler.Filename
	confirm.FileNameEdit = RandStringBytes(10)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		confirm.Code = 2
		confirm.Message = err.Error()
	}
	//	fmt.Println("data: " + confirm.Message);

	/*	--------------------	Create folder and files	--------------------	*/
	err = ioutil.WriteFile("./uploaded/"+confirm.FileName, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		confirm.Code = 2
		confirm.Message = err.Error()
	}

	err = ioutil.WriteFile("./uploaded/"+confirm.FileNameEdit, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		confirm.Code = 2
		confirm.Message = err.Error()
	}

	/*	--------------------	Create folder and files	--------------------	*/
	b, err := json.Marshal(confirm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(confirm)
	w.Write(b)
}

func jsonHandler(w http.ResponseWriter, req *http.Request, title string) {
	//	fmt.Println("json")
	//	fmt.Println(req.Header.Get)
	//	fmt.Println(req.ContentLength)
	//	fmt.Println(req.Method)
	//	fmt.Println(req.Form)
	confirm := &Status{
		Code:         1,
		Message:      "Archivo editado exitosamente.",
		FileName:     "",
		FileNameEdit: "",
	}
	var data imageEdit.InfoJSON

	if req.Body == nil {
		confirm.Code = 2
		confirm.Message = "Please send a request body."
	} else {
		bodyb, err := ioutil.ReadAll(req.Body)
		if err != nil {
			confirm.Code = 2
			confirm.Message = "Please send a request body."
		} else {
			//fmt.Println(string(bodyb))
			err = json.Unmarshal([]byte(string(bodyb)), &data)
			if err != nil {
				confirm.Code = 2
				confirm.Message = err.Error()
			} else {
				//				m := validPathImage.FindStringSubmatch(data.FileName)
				/*	--------------------------	Edit Image	------------------------------	*/
				//				fmt.Println(data.Operation)
				//				fmt.Println(data.FileName)
				//				fmt.Println(data.FileNameEdit)
				//				fmt.Println(data.Args)

				newFile := RandStringBytes(10)
				err = imageEdit.Edit(data, newFile, "./uploaded/")
				//				fmt.Println(newFile)
				if err != nil {
					confirm.Code = 2
					confirm.Message = err.Error()
				} else {
					err = os.Remove("./uploaded/" + data.FileNameEdit)
					if err != nil {
						fmt.Println(err)
						//return
					}

					confirm.FileName = data.FileName
					confirm.FileNameEdit = newFile
				}
			}
		}
	}

	// -----
	b, err := json.Marshal(confirm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	//fmt.Fprintf(w, string(b))
	//fmt.Println("Sending: ............................." + string(b));
}

/*	--------------------------->	Main	--------------------------->	*/

func main() {
	err := os.RemoveAll("./uploaded")
	if err != nil {
		fmt.Println(err)
	} else {
		err = os.Mkdir("./uploaded", 0600)
		if err != nil {
			err = os.Mkdir("./uploaded", 0777)
			if err != nil {
				err = os.Mkdir("./uploaded", 0600)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	http.HandleFunc("/", makeHandler(homeHandler))
	http.HandleFunc("/upload/", makeHandler(uploadHandler))
	http.HandleFunc("/json/", makeHandler(jsonHandler))

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./fonts/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))
	http.Handle("/uploaded/", http.StripPrefix("/uploaded/", http.FileServer(http.Dir("./uploaded/"))))

	port := flag.Int("port", 8080, "port to serve on")
	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	log.Printf("Servidor listo en: %s\n", addr)

	err = http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())

}
