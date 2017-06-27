//https://github.com/golang-samples/http/blob/master/fileupload/main.go

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"encoding/json"
	"github.com/satori/go.uuid"
	"strings"
	"path/filepath"
)

// 100 MB
const MAX_MEMORY = 100 * 1024 * 1024
const extension string = "heic"

type HeifConfig struct {
	General General `json:"general"`
	Content []Content `json:"content"`
}

type General struct {
	Output Output `json:"output"`
	Brands Brands `json:"brands"`
}

type Output struct {
	FilePath string `json:"file_path"`
}

type Brands struct {
	Major string `json:"major"`
	Other []string `json:"other"`
}

type Content struct {
	Master Master `json:"master"`
	Thumbs []Thumbs `json:"thumbs"`
}

type Master struct {
	FilePath string `json:"file_path"`
	HdlrType string `json:"hdlr_type"`
	CodeType string `json:"code_type"`
	EncpType string `json:"encp_type"`
}

type Thumbs struct {
	FilePath string `json:"file_path"`
	HdlrType string `json:"hdlr_type"`
	CodeType string `json:"code_type"`
	EncpType string `json:"encp_type"`
	SyncRate int	`json:"sync_rate"`
}

func upload(w http.ResponseWriter, r *http.Request) {
	guid := uuid.NewV4()
	var UPLOAD_FOLDER string = "files/" + guid.String()
	var ffmpeg_cmd = "/usr/local/bin/ffmpeg"
	var heif_cmd = "/usr/local/bin/writerapp"
	var args = []string{}
	var main_file = UPLOAD_FOLDER + "/bitstream.265"
	var thumb_file = UPLOAD_FOLDER + "/bitstream.thumb.265"
	var configfile = UPLOAD_FOLDER + "/config.json"

	if _, err := os.Stat(UPLOAD_FOLDER); os.IsNotExist(err) {
		if create_folder_err := os.Mkdir(UPLOAD_FOLDER, 0755); create_folder_err != nil {
			log.Println("Unable to create upload folder '" + UPLOAD_FOLDER + "'.")
			return
		} else {
			log.Println("Upload folder '" + UPLOAD_FOLDER + "' created.")
		}
	}

	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
	}

	for key, value := range r.MultipartForm.Value {
		fmt.Fprintf(w, "%s:%s", key, value)
		log.Printf("%s:%s", key, value)
	}

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			filename := fileHeader.Filename
			file, _ := fileHeader.Open()
			path := fmt.Sprintf(UPLOAD_FOLDER + "/%s", filename)
			buf, _ := ioutil.ReadAll(file)
			ioutil.WriteFile(path, buf, 0644)
			new_file_name := UPLOAD_FOLDER + "/" + strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + extension
			log.Println("File '" + path + "' saved.")
			log.Println("New filename: " + new_file_name)


			// Create json
			// https://stackoverflow.com/a/24356483/319826
			heif_config := &HeifConfig {
				General: General {
					Output: Output {
						FilePath: new_file_name,
					},
					Brands: Brands {
						Major: "mif1",
						Other: []string{"mif1", "heic", "hevc"},
					},
				},
				Content: []Content {
					{Master: Master {
						FilePath: main_file,
						HdlrType: "pict",
						CodeType: "hvc1",
						EncpType: "meta",
					},
						Thumbs: []Thumbs {
							{
								FilePath: thumb_file,
								HdlrType: "pict",
								CodeType: "hvc1",
								EncpType: "meta",
								SyncRate: 1,
							},
						},
					},
				},
			}
			// End create json
			b, err := json.Marshal(heif_config)
			if err != nil {
				log.Println(err)
			} else {
				fileErr := ioutil.WriteFile(configfile, b, 0644)
				if fileErr != nil {
					log.Println(fileErr)
					return
				}
			}

			pwd, err := exec.Command("pwd").Output()
			if err != nil {
				log.Println(pwd)
			}

			args = []string{"-y", "-i", path, "-crf", "12", "-preset", "slower", "-pix_fmt", "yuv420p", "-f", "hevc", main_file}
			if err := exec.Command(ffmpeg_cmd, args...).Run(); err != nil {
				log.Println("Unable to create main file!")
				return
			}

			args = []string{"-y", "-i", path, "-vf", "scale=320:240", "-crf", "28", "-preset", "slower", "-pix_fmt", "yuv420p", "-f", "hevc", thumb_file}
			if err := exec.Command(ffmpeg_cmd, args...).Run(); err != nil {
				log.Println("Unable to create thumbnail!")
				return
			}

			//build_heif_file = "/usr/local/bin/writerapp " + UPLOAD_FOLDER + "/" + configfile
			args = []string{configfile}
			if err = exec.Command(heif_cmd, args...).Run(); err != nil {
				log.Println("Unable to encode to HEIF-format!")
				log.Println(err)
				return
			}

			log.Println("File '" + path + "' converted to HEIF-format!")
		}
	}
}

func main() {
	http.HandleFunc("/upload", upload)
	http.Handle("/", http.FileServer(http.Dir("static")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
