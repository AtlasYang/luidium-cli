package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var initialization = &cobra.Command{
	Use:   "init",
	Short: "Initialize the luidium cli",
	Long:  `Initialize the luidium cli.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide proper block url `application/version/blockName`")
			return
		}
		if len(strings.Split(args[0], "/")) != 3 {
			fmt.Println("Please provide proper block url `application/version/blockName`")
			return
		}
		if viper.GetString("cli_token") == "NOT_AUTHENTICATED" {
			fmt.Println("Luidium CLI is not authenticated. Please run `luidium certify <YOUR_CLI_TOKEN>` to authenticate")
			return
		}
		if viper.GetString("config") == "CONFIGURED" {
			fmt.Println("Luidium CLI is already configured")
			return
		}
		blockConfig := strings.Split(args[0], "/")
		application := blockConfig[0]
		version := blockConfig[1]
		blockName := blockConfig[2]
		CheckConfig(application, version, blockName)
	},
}

func GenerateConfig(application string, version string, blockName string) {
	config_name := "./" + blockName + "/luidium.toml"
	config_content := fmt.Sprintf(`# This file is used to configure the Luidium CLI
# You can download template configuration from Application Dashboard

# Block configuration (Do not modify)
application = "%s"
version = "%s"
name = "%s"

# If this option is disabled, luidium cli will upload all files in the project directory
# Use this option to specify which files should be uploaded
use_gitignore = true

# If this option is enabled, it will load the .env file from given path (default is project root)
# This will *override* the environment variables you set in Application Dashboard
use_dotenv = false`, application, version, blockName)

	configPath := filepath.Join(".", blockName)
	err := os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to create directories:", err)
		return
	}

	file, err := os.Create(config_name)
	if err != nil {
		fmt.Println("Error creating config file")
		return
	}
	defer file.Close()

	file.WriteString(config_content)
}

func GetObjectKeys(application string, version string, blockName string) []string {
	payload := StorageRequest{
		Action:    StorageActionListObject,
		Bucket:    application,
		ObjectKey: fmt.Sprintf("%s/%s", version, blockName),
	}

	url := "https://luidium-storage-api.lighterlinks.io/storage/single"

	client := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Content-Type", "application/json")
	req.SetBody([]byte(fmt.Sprintf(`{"action":"%s","bucket":"%s","object_key":"%s"}`, payload.Action, payload.Bucket, payload.ObjectKey)))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	if err := client.Do(req, resp); err != nil {
		fmt.Println("Error getting object keys")
		return nil
	}

	var body map[string]interface{}
	err := json.Unmarshal(resp.Body(), &body)
	if err != nil {
		fmt.Println("Error parsing response")
		return nil
	}

	var objectKeys []string
	for _, object := range body["object_keys"].([]interface{}) {
		objectKeys = append(objectKeys, object.(string))
	}

	return objectKeys
}

func ReadLuidiumConfig(application string, version string, blockName string) LuidiumConfig {
	url := "https://luidium-storage-api.lighterlinks.io/storage/config/" + application + "/" + version + "/" + blockName

	client := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	if err := client.Do(req, resp); err != nil {
		fmt.Println("Error reading config")
		return LuidiumConfig{}
	}

	var config LuidiumConfig
	err := json.Unmarshal(resp.Body(), &config)
	if err != nil {
		fmt.Println("Error parsing response")
		return LuidiumConfig{}
	}

	return config
}

func DownloadObjectToPath(application string, objectKey string, folderPath string) {
	url := "https://luidium-storage-api.lighterlinks.io/storage/download/" + application + "/" + objectKey

	client := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	if err := client.Do(req, resp); err != nil {
		fmt.Println("Error downloading object")
		return
	}

	// to get file name, exclude first 2 parts of the object key (version and block name)
	fileName := strings.Join(strings.Split(objectKey, "/")[2:], "/")
	filePath := filepath.Join(folderPath, fileName)

	dirPath := filepath.Dir(filePath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to create directories:", err)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(resp.Body())
	if err != nil {
		fmt.Println("Failed to write to file:", err)
		return
	}
}

func CheckConfig(application string, version string, blockName string) {
	config_name := "luidium.toml"
	_, err := os.ReadFile(config_name)
	if err != nil {
		GenerateConfig(application, version, blockName)
		fmt.Println("Config file generated")

		objectKeys := GetObjectKeys(application, version, blockName)
		config := ReadLuidiumConfig(application, version, blockName)

		var newObjectKeys = []string{}
		for _, objectKey := range objectKeys {
			ignore := false
			for _, ignoreFile := range config.IgnoreFiles {
				if strings.Contains(objectKey, ignoreFile) || strings.Contains(objectKey, "luidium-config.json") {
					ignore = true
					break
				}
			}
			if !ignore {
				newObjectKeys = append(newObjectKeys, objectKey)
			}
		}

		fmt.Println("Pulling files from storage...")

		for _, objectKey := range newObjectKeys {
			DownloadObjectToPath(application, objectKey, "./"+blockName)
		}
	} else {
		fmt.Println("Config file exists")
	}
}

func init() {
	rootCmd.AddCommand(initialization)
}
