package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const OPERATIONS = 100

func SendCookie(w http.ResponseWriter, req *http.Request, newFile string) {
	/*	--------------------------->	Cookie	--------------------------->	*/
	cookiePrin, err := req.Cookie("imageAnalysis")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}
	vals := strings.Split(cookiePrin.Value, ",")
	ID := vals[0]

	var name, value string
	opCount, err := strconv.ParseInt(vals[1], 10, 32)
	if err == nil {
		switch opCount {
		case 0:
			opCount++
			name = string("op1")
			value = newFile
			cookie := http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: true}
			http.SetCookie(w, &cookie)

			name = cookiePrin.Name
			value = ID + "," + strconv.Itoa(int(opCount))
			cookie = http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: false}
			http.SetCookie(w, &cookie)
		default:
			cookieOp, err := req.Cookie("op1")
			if err != nil {
				http.Redirect(w, req, "/", http.StatusFound)
				return
			}
			if opCount < OPERATIONS {
				name = cookieOp.Name
				vals = strings.Split(cookieOp.Value, ",")
				for i := 0; i < int(opCount); i++ {
					value += vals[i] + ","
				}
				value += newFile
				cookie := http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: true}
				http.SetCookie(w, &cookie)

				opCount++
				name = cookiePrin.Name
				value = ID + "," + strconv.Itoa(int(opCount))
				cookie = http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: false}
				http.SetCookie(w, &cookie)
			} else {
				name = cookieOp.Name
				vals = strings.Split(cookieOp.Value, ",")

				// Delete file
				err = os.Remove("./uploaded/" + ID + "/" + vals[0])
				if err != nil {
					fmt.Println(err)
				} //*/

				for i := 1; i < int(opCount); i++ {
					value += vals[i] + ","
				}
				value += newFile
				cookie := http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: true}
				http.SetCookie(w, &cookie)
			}
		}
	}
	/*	<---------------------------	Cookie	<---------------------------	*/
}

func Undo(w http.ResponseWriter, req *http.Request) bool {
	var name, value string
	var cookie http.Cookie

	cookiePrin, err := req.Cookie("imageAnalysis")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusFound)
		return false
	}
	vals := strings.Split(cookiePrin.Value, ",")
	ID := vals[0]
	//fmt.Printf("vals:%s", vals)
	opCount, _ := strconv.ParseInt(vals[1], 10, 32)

	//	Undo
	if opCount > 1 {
		opCount--
		// Cookie
		name = cookiePrin.Name
		value = ID + "," + strconv.Itoa(int(opCount))
		cookie = http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: false}
		http.SetCookie(w, &cookie)
		return true
	}
	return false
}

func Redo(w http.ResponseWriter, req *http.Request) bool {
	var name, value string
	var cookie http.Cookie
	var paths int

	cookiePrin, err := req.Cookie("imageAnalysis")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusFound)
		return false
	}
	vals := strings.Split(cookiePrin.Value, ",")
	ID := vals[0]
	//fmt.Printf("vals:%s", vals)
	opCount, _ := strconv.ParseInt(vals[1], 10, 32)

	if opCount > 0 {
		cookieOp, _ := req.Cookie("op1")
		vals2 := strings.Split(cookieOp.Value, ",")
		paths = len(vals2)
		//fmt.Printf("%d:%s", paths, vals2)
	}

	//	Redo
	if int(opCount) < paths {
		opCount++
		// Cookie
		name = cookiePrin.Name
		value = ID + "," + strconv.Itoa(int(opCount))
		cookie = http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: false}
		http.SetCookie(w, &cookie)
		return true
	}
	return false
}
