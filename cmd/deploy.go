package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var deploy = &cobra.Command{
	Use:   "deploy",
	Short: "Upload files to the storage server and deploy the service",
	Long:  `Upload files to the storage server. Automatically deploy the service.`,
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
		ignore_list := ReadIgnoreFile()
		SendClearFolderRequest()
		ProcessDir(".", ignore_list)
		SendDeployRequest()
		elapsed := time.Since(start_time)
		fmt.Println("Time taken: ", elapsed)
		fmt.Println("Upload completed and service deployed.")
	},
}

func SendClearFolderRequest() {
	url := "https://luidium-storage-api.lighterlinks.io/storage/single"

	payload := StorageRequest{
		Action:    StorageActionDeleteFolderExceptConfig,
		Bucket:    viper.GetString("application"),
		ObjectKey: fmt.Sprintf("%s/%s", viper.GetString("version"), viper.GetString("name")),
	}

	client := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Content-Type", "application/json")
	req.SetBody([]byte(fmt.Sprintf(`{"action":"%s","bucket":"%s","object_key":"%s"}`, payload.Action, payload.Bucket, payload.ObjectKey)))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	if err := client.Do(req, resp); err != nil {
		fmt.Println("Error sending clear folder request")
		return
	}
}

func SendDeployRequest() {
	url := "https://luidium-main-api.lighterlinks.io/block/run_by_cli"

	payload := RunByCliRequest{
		Token:     viper.GetString("cli_token"),
		AppName:   viper.GetString("application"),
		Version:   viper.GetString("version"),
		BlockName: viper.GetString("name"),
	}

	client := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Content-Type", "application/json")
	req.SetBody([]byte(fmt.Sprintf(`{"token":"%s","app_name":"%s","version":"%s","block_name":"%s"}`, payload.Token, payload.AppName, payload.Version, payload.BlockName)))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	if err := client.Do(req, resp); err != nil {
		fmt.Println("Error sending deploy request")
		return
	}
}

func init() {
	rootCmd.AddCommand(deploy)
}
