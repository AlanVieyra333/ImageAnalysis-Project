/**
 * Created by Fenix on 03/10/2016.
set GOPATH=%cd%
go build editImage
go install server
*/

package main

import (
	"editImage"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"server"
	"strconv"
	"strings"
	"tools"
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

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		//fmt.Println("" + m[0] + " " + m[1] + " " + m[2]);
		fn(w, r)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func newID() string {
	id := tools.GetRandInteger(10)
	//fmt.Printf("id:%s\n", id)
	for exist, _ := exists("./uploaded/" + id + "/"); exist; exist, _ = exists("./uploaded/" + id) {
		id = tools.GetRandInteger(10)
	}
	os.Mkdir("./uploaded/"+id+"/", 0777)
	return id
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	/*/fmt.Println(title);
	p, err := loadPage("index")
	if err != nil {
		http.NotFound(w, r)
		return
	}*/
	//renderTemplate(w, title)

	/*	----	Cookie	---	*/
	if len(r.Cookies()) == 0 { //	Create new cookie
		// Value = [IDClient],[currentOperation]
		cookie := http.Cookie{Name: "imageAnalysis", Value: newID() + ",0", Path: "/", HttpOnly: false}
		http.SetCookie(w, &cookie)
	}

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
	Data         editImage.DataOutJSON
}

func uploadHandler(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("Upl")
	//fmt.Println(req.Header)
	cookiePrin, err := req.Cookie("imageAnalysis")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}
	vals := strings.Split(cookiePrin.Value, ",")
	ID := vals[0]

	confirm := &Status{
		Code:         1,
		Message:      "Archivo subido exitosamente.",
		FileName:     "",
		FileNameEdit: "",
	}
	isNewImage := req.FormValue("type")
	file, handler, err := req.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		confirm.Code = 2
		confirm.Message = err.Error()
	}
	//	fmt.Println("file: " + confirm.Message);
	confirm.FileName = handler.Filename
	confirm.FileNameEdit = tools.GetRandString(10)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		confirm.Code = 2
		confirm.Message = err.Error()
	}
	//	fmt.Println("data: " + confirm.Message);

	//	New Image.
	if isNewImage == "0" {
		tools.RemovePath("./uploaded/" + ID)
		/*	----	Cookie	---	*/
		name := cookiePrin.Name
		value := ID + ",1"
		cookie := http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: false}
		http.SetCookie(w, &cookie)

		name = string("op1")
		value = confirm.FileNameEdit
		cookie = http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: true}
		http.SetCookie(w, &cookie)
	}
	/*	--------------------	Create folder and files	--------------------	*/
	err = ioutil.WriteFile("./uploaded/"+ID+"/"+confirm.FileNameEdit, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		confirm.Code = 2
		confirm.Message = err.Error()
	}
	/*	--------------------	Create folder and files	--------------------	*/

	/*	----	JSON	---	*/
	b, err := json.Marshal(confirm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func jsonHandler(w http.ResponseWriter, req *http.Request) {
	//fmt.Println(req.Header)
	cookiePrin, err := req.Cookie("imageAnalysis")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}
	vals := strings.Split(cookiePrin.Value, ",")
	ID := vals[0]

	confirm := &Status{
		Code:         1,
		Message:      "Archivo editado exitosamente.",
		FileName:     "",
		FileNameEdit: "",
	}
	var data editImage.InfoJSON

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
			} else if data.Operation == -1 { // Cookie
				args := strings.Split(data.Args, ";")
				aux, err := strconv.ParseInt(args[0], 10, 64)

				if err != nil {
					confirm.Code = 2
					confirm.Message = err.Error()
				} else {
					confirm.FileName = data.FileName
					confirm.FileNameEdit = data.FileNameEdit
					var nothing editImage.DataOutJSON
					confirm.Data = nothing

					cookieOp, err := req.Cookie("op1")
					var paths []string
					var op int64
					if err != nil {
						confirm.Code = 2
						confirm.Message = "Invalid cookie."
					} else {
						paths = strings.Split(cookieOp.Value, ",")
						op, _ = strconv.ParseInt(vals[1], 10, 32)
					}

					switch aux {
					case 1: //	Undo.
						if server.Undo(w, req) {
							op--
							confirm.FileNameEdit = paths[int(op-1)]
						}
					case 2: //	Redo.
						if server.Redo(w, req) {
							op++
							confirm.FileNameEdit = paths[int(op-1)]
						}
					default:
						confirm.Code = 2
						confirm.Message = "Invalid option."
					}
				}
			} else {
				//				m := validPathImage.FindStringSubmatch(data.FileName)
				/*	--------------------------	Edit Image	------------------------------	*/

				newFile := tools.GetRandString(10)
				dataJSON, err := editImage.Edit(data, newFile, "./uploaded/"+ID+"/")
				//				fmt.Println(newFile)
				if err != nil {
					confirm.Code = 2
					confirm.Message = err.Error()
				} else {
					confirm.FileName = data.FileName
					confirm.FileNameEdit = newFile
					confirm.Data = dataJSON

					server.SendCookie(w, req, newFile)
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
}

/*	--------------------------->	Main	--------------------------->	*/
func main() {
	tools.RemovePath("./uploaded")

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

	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())

}
