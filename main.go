package main

import (
	"archive/zip"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"github.com/rs/cors"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

const maxFileNameLength = 255                              // 文件名的最大长度
var maxStorageSize int64                                   // 新增变量来存储最大存储量
var validFileName = regexp.MustCompile(`^[^<>:"/\\|?*]*$`) // 修改正则表达式，排除Windows中不允许的字符
var rootPath string                                        // 服务所在的文件目录

type File struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	IsDir   bool      `json:"isDir"`
}

type LoggingHandler struct {
	Handler http.Handler
	Logger  *log.Logger
}

func (h *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	h.Handler.ServeHTTP(w, r)
	h.Logger.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
}

func main() {
	rootPath = path.Join(".", "data")
	dir := flag.String("d", rootPath, "启动服务所在的文件目录（默认是当前目录）")
	port := flag.String("p", "8000", "监听的端口（默认是8000）")
	maxStorage := flag.String("s", "1G", "文件夹的最大存储量（默认是1G，如果设置为0则表示不限制）")

	// 解析命令行参数
	flag.Parse()

	// 解析最大存储量
	var err error
	maxStorageSize, err = parseStorageSize(*maxStorage)
	if err != nil {
		log.Fatalf("无法解析最大存储量：%v", err)
	}

	// 创建文件夹
	rootPath = *dir
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		os.MkdirAll(rootPath, os.ModePerm)
	}

	// 创建CORS中间件
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // 允许所有来源，你可以替换为具体的域名
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
	})

	// 使用CORS中间件包装HTTP服务器
	handler := corsConfig.Handler(http.DefaultServeMux)

	// 处理静态文件的请求
	fs := http.FS(staticFiles)
	http.Handle("/", http.FileServer(fs))

	// 将API路由更改为以api/开头
	http.HandleFunc("/api/addFolder", addFolderHandler)
	http.HandleFunc("/api/delete", deleteHandler)
	http.HandleFunc("/api/rename", renameHandler)
	http.HandleFunc("/api/download", downloadHandler)
	http.HandleFunc("/api/upload", uploadHandler)
	http.HandleFunc("/api/files", filesHandler)

	ip := getLocalIP()
	log.Printf("启动HTTP服务: http://%s:%s/static 服务所在的文件目录是 %s \n", ip, *port, *dir)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+*port, handler))
}

func addFolderHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	path := r.URL.Query().Get("path") // 获取新的 path 参数
	if err := validateFileName(name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 在指定的路径下创建新的文件夹
	os.MkdirAll(filepath.Join(rootPath, path, name), os.ModePerm)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	path := filepath.Join(rootPath, name)
	err := os.Remove(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renameHandler(w http.ResponseWriter, r *http.Request) {
	oldName := r.URL.Query().Get("oldName")
	newName := r.URL.Query().Get("newName")
	newPath := filepath.Join(filepath.Dir(oldName), newName)
	err := os.Rename(filepath.Join(rootPath, oldName), filepath.Join(rootPath, newPath))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	filePath := filepath.Join(rootPath, name)
	fileName := filepath.Base(filePath)
	info, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		// 创建一个zip文件
		w.Header().Add("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName)+".zip")
		w.Header().Add("Content-Type", "application/zip")
		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()

		// 将文件夹中的所有文件添加到zip文件中
		filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(filePath, path)
			if err != nil {
				return err
			}

			zipFile, err := zipWriter.Create(relPath)
			if err != nil {
				return err
			}

			fsFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fsFile.Close()

			_, err = io.Copy(zipFile, fsFile)
			return err
		})
	} else {
		encodedFileName := url.QueryEscape(fileName)
		w.Header().Add("Content-Disposition", "attachment; filename="+encodedFileName)
		w.Header().Add("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, filePath)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(10 << 20) // limit your max input length!
		file, handler, err := r.FormFile("file")
		path := r.FormValue("path") // 获取新的 path 参数
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// 检查文件大小
		if maxStorageSize > 0 {
			fileSize := handler.Size
			dirSize, _ := getDirSize(rootPath)
			if fileSize+dirSize > maxStorageSize {
				http.Error(w, "上传的文件会超过文件夹的最大存储量", http.StatusBadRequest)
				return
			}
		}

		fileName := handler.Filename
		if err := validateFileName(fileName); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filePath := filepath.Join(rootPath, path, fileName) // Append the rootPath and path parameter to the file name
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		io.Copy(f, file)
	}
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" || path == "/" {
		path = "."
	}
	files, err := os.ReadDir(filepath.Join(rootPath, path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileInfos := make([]File, len(files))
	for i, file := range files {
		info, _ := file.Info()
		fileInfos[i] = File{
			Name:    file.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   file.IsDir(),
		}
	}

	json.NewEncoder(w).Encode(fileInfos)
}

func getLocalIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 {
			continue // interface down, or loopback interface
		}
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
		}
	}
	return ""
}

// 获取文件夹的大小
func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// 解析存储量的参数
func parseStorageSize(sizeStr string) (int64, error) {
	if sizeStr == "0" || sizeStr == "0K" || sizeStr == "0M" || sizeStr == "0G" || sizeStr == "0T" {
		return 0, nil // 0表示无限制
	}

	var size float64
	var err error
	var multiplier int64

	if strings.HasSuffix(sizeStr, "T") {
		size, err = strconv.ParseFloat(sizeStr[:len(sizeStr)-1], 64)
		multiplier = 1 << 40 // 转换为字节
	} else if strings.HasSuffix(sizeStr, "G") {
		size, err = strconv.ParseFloat(sizeStr[:len(sizeStr)-1], 64)
		multiplier = 1 << 30 // 转换为字节
	} else if strings.HasSuffix(sizeStr, "M") {
		size, err = strconv.ParseFloat(sizeStr[:len(sizeStr)-1], 64)
		multiplier = 1 << 20 // 转换为字节
	} else if strings.HasSuffix(sizeStr, "K") {
		size, err = strconv.ParseFloat(sizeStr[:len(sizeStr)-1], 64)
		multiplier = 1 << 10 // 转换为字节
	} else {
		size, err = strconv.ParseFloat(sizeStr, 64)
	}

	if err != nil {
		return 0, err
	}

	return int64(size * float64(multiplier)), nil
}

func validateFileName(fileName string) error {
	if len(fileName) > maxFileNameLength {
		return errors.New("文件名过长")
	}
	if !validFileName.MatchString(fileName) {
		return errors.New("文件名包含非法字符")
	}
	if strings.Contains(fileName, "..") {
		return errors.New("文件名包含非法路径")
	}
	return nil
}
