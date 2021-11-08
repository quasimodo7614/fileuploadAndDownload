package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const BaseUploadPath = "/var/file"

var indexHtml = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">

    <title>Upload Files</title>

</head>

<body>
<form method="post" enctype="multipart/form-data">
    <input type="file" name="file" multiple>
    <input type="submit" value="Upload File" name="submit">
</form>

<script type="text/javascript">
    var puburi = "%s"

    async function makeList() {
        const url = 'http://' + puburi + ':9000/files';
        const downLoadurl = 'http://' + puburi + ':9000/download?filename=';
        const response = await fetch(url);
        var data = await response.json();
        console.info(data)
        // Establish the array which acts as a data source for the list
        let listData = data,
            // Make a container element for the list
            listContainer = document.createElement('div'),
            // Make the list
            listElement = document.createElement('ul'),
            // Set up a loop that goes through the items in listItems one at a time
            numberOfListItems = listData.length,
            listItem,
            i;


        // Add it to the page
        document.getElementsByTagName('body')[0].appendChild(listContainer);
        listContainer.appendChild(listElement);

        for (i = 0; i < numberOfListItems; ++i) {
            // create an item for each one
            listItem = document.createElement('li');

            // Add the item text
            listItem.innerHTML = "<a href='" + downLoadurl + listData[i] + "'>" + listData[i] + "</a>";
            // Add listItem to the listElement
            listElement.appendChild(listItem);
        }
    }

    // Usage
    makeList();
</script>

<script>
    var puburi = "%s"
    const url = 'http://'+puburi+':9000/upload';
    const form = document.querySelector('form');

    form.addEventListener('submit', e => {
        e.preventDefault();

        const files = document.querySelector('[type=file]').files;
        const formData = new FormData();

        for (let i = 0; i < files.length; i++) {
            let file = files[i];

            formData.append('file', file);
        }

        fetch(url, {
            method: 'POST',
            body: formData
        }).then(response => {
            return response.text();
        }).then(data => {
            alert(data+" 上传成功，请刷新页面");
        });
    });
</script>
</body>

</html>
`

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/download", handleDownload)
	http.HandleFunc("/files", handleFiles)

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("Server run fail:", err)
	}
}

func handleIndex(res http.ResponseWriter, req *http.Request) {

	host := "localhost"
	if s := os.Getenv("PUB_HOST"); s != "" {
		host = s
	}
	fmt.Fprint(res, fmt.Sprintf(indexHtml, host, host))
}
func handleUpload(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//文件上传只允许POST方法
	if request.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method not allowed"))
		return
	}

	//从表单中读取文件
	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		log.Println(err)
		_, _ = io.WriteString(w, "Read file error")
		return
	}
	//defer 结束时关闭文件
	defer file.Close()
	log.Println("filename: " + fileHeader.Filename)

	//创建文件
	newFile, err := os.Create(BaseUploadPath + "/" + fileHeader.Filename)
	if err != nil {
		_, _ = io.WriteString(w, "Create file error")
		return
	}
	//defer 结束时关闭文件
	defer newFile.Close()

	//将文件写到本地
	_, err = io.Copy(newFile, file)
	if err != nil {
		_, _ = io.WriteString(w, "Write file error")
		return
	}
	_, _ = io.WriteString(w, "Upload success")
}

func handleDownload(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//文件上传只允许GET方法
	if request.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method not allowed"))
		return
	}
	//文件名
	filename := request.FormValue("filename")
	if filename == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
	log.Println("filename: " + filename)
	//打开文件
	file, err := os.Open(BaseUploadPath + "/" + filename)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
	//结束后关闭文件
	defer file.Close()

	//设置响应的header头
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment; filename=\""+filename+"\"")
	//将文件写至responseBody
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
}

func handleFiles(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//文件上传只允许GET方法
	if request.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method not allowed"))
		return
	}

	files, err := ioutil.ReadDir(BaseUploadPath)
	if err != nil {
		log.Println(err)
		return
	}

	fs := []string{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fs = append(fs, file.Name())
	}
	b, err := json.Marshal(fs)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(b)

}

