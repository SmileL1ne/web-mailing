package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"code.sajari.com/docconv"
	"github.com/SmileL1ne/web-mailing/model"
)

func ReadForm(r *http.Request, sub model.Subscriber) (model.Subscriber, error) {
	// err := json.NewDecoder(r.Body).Decode(&subs)
	// if err != nil {
	// 	log.Println(err)
	// 	return model.Subscriber{}, err
	// }
	// return subs, nil
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return model.Subscriber{}, err
	}
	sub = model.Subscriber{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Interest:  r.Form.Get("interest"),
	}
	return sub, nil
}

func JSONWriter(w http.ResponseWriter, msg string, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func ReadMultiForm(w http.ResponseWriter, r *http.Request, mail model.MailUpload) (model.MailUpload, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Fatalln(err)
	}

	form := r.MultipartForm
	mail.DocxName = form.Value["docx_name"][0]
	mail.Date = time.Now()

	file, ok := form.File["docx"]
	if !ok {
		return model.MailUpload{}, fmt.Errorf("unable to find uploaded document")
	}

	if file[0].Filename != "" {
		fileExtension := filepath.Ext(file[0].Filename)

		f, err := file[0].Open()
		if err != nil {
			return model.MailUpload{}, fmt.Errorf("unable to open uploaded document")
		}
		defer f.Close()

		switch fileExtension {
		case ".txt":
			scanner := bufio.NewScanner(f)

			for scanner.Scan() {
				line := fmt.Sprintf("%s<br>", scanner.Text())
				mail.DocxContent += line
			}
			if scanner.Err() != nil {
				log.Fatalln(err)
			}
		case ".doc", ".docx":
			res, _, err := docconv.ConvertDocx(f)
			if err != nil {
				log.Fatalln(err)
			}

			lines := strings.Split(res, "\n")
			var content string
			for _, line := range lines {
				content += line + "<br>"
			}
			mail.DocxContent = content
		default:
			return model.MailUpload{}, fmt.Errorf("file format not allowed; try using .doc, .docx or .txt")
		}
	}

	return mail, nil
}

func HTMLRender(w http.ResponseWriter, r *http.Request, dt any) error {
	filePath := "./templates/index.html"
	tmp, err := template.ParseFiles(filePath)
	if err != nil {
		return fmt.Errorf("HTMLReader Error: failed to parse file: %v\n", err)
	}

	err = tmp.Execute(w, dt)
	if err != nil {
		return fmt.Errorf("HTMLReader Error: failed to execute template: %v\n", err)
	}

	return nil
}
