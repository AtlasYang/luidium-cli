package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var upload = &cobra.Command{
	Use:   "upload",
	Short: "Upload files to the storage server",
	Long:  `Upload files to the storage server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a directory")
			return
		}
		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			fmt.Println("Directory does not exist")
			return
		}
		if viper.GetString("cli_token") == "NOT_AUTHENTICATED" {
			fmt.Println("Luidium CLI is not authenticated. Please run `luidium certify <YOUR_CLI_TOKEN>` to authenticate")
			return
		}
		if viper.GetString("config") == "NOT_CONFIGURED" {
			fmt.Println("Luidium CLI is not configured. Please run `luidium init <APPLICATION/VERSION/BLOCK_NAME>` to configure")
			return
		}
		start_time := time.Now()
		SendClearFolderRequest()
		ignore_list := ReadIgnoreFile()
		ProcessDir(args[0], ignore_list)
		elapsed := time.Since(start_time)
		fmt.Println("Time taken: ", elapsed)
		fmt.Println("Upload complete")
	},
}

func ReadIgnoreFile() []string {
	gitignore := []string{}
	use_gitignore := viper.GetString("use_gitignore")
	if use_gitignore == "false" {
		return gitignore
	}

	file, err := os.Open(".gitignore")
	if err != nil {
		return gitignore
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content := scanner.Text()
		if content == "" {
			continue
		}
		if content[0] == '#' {
			continue
		}
		if content[0] == '/' {
			gitignore = append(gitignore, content[1:])
		} else {
			gitignore = append(gitignore, "**/"+content)
		}
	}

	if err := scanner.Err(); err != nil {
		return gitignore
	}

	return gitignore
}

func ProcessDir(dir string, exclude_patterns []string) {
	appName := viper.GetString("application")
	version := viper.GetString("version")
	name := viper.GetString("name")

	fmt.Println("Uploading to ", appName+"/"+version+"/"+name)
	prefix := appName + "/" + version + "/" + name + "/"

	wg := &sync.WaitGroup{}
	file_list := []string{}

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(dir, path)
		for _, pattern := range exclude_patterns {
			matched, _ := filepath.Match(pattern, relPath)
			if matched {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		if d.IsDir() {
			return nil
		}

		file_list = append(file_list, path)
		return nil
	},
	)

	fmt.Println("Total files: ", len(file_list))

	for _, file := range file_list {
		wg.Add(1)
		go func(file string) {
			UploadFile(prefix, file)
			wg.Done()
		}(file)
	}

	wg.Wait()
}

func UploadFile(prefix string, path string) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return
	}

	url := "https://luidium-storage-api.lighterlinks.io/storage/upload/" + prefix + path

	// multipart form data: key: file, value: file data, key: size, value: file size
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", path)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	part.Write(fileData)
	writer.WriteField("size", fmt.Sprintf("%d", len(fileData)))
	err = writer.Close()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	client := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPut)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBody(body.Bytes())
	resp := fasthttp.AcquireResponse()
	err = client.Do(req, resp)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func init() {
	rootCmd.AddCommand(upload)
}
