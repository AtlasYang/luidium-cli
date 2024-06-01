package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

var certify = &cobra.Command{
	Use:   "certify",
	Short: "Authenticate the CLI with CLI token",
	Long:  `Authenticate the CLI with CLI token.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a token")
			return
		}
		ValidateToken(args[0])
	},
}

func ValidateToken(token string) {
	url := "https://luidium-main-api.lighterlinks.io/cli/validate/" + token

	client := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	if err := client.Do(req, resp); err != nil {
		fmt.Println("Error validating token")
		return
	}

	var body SingleStringResponse

	err := json.Unmarshal(resp.Body(), &body)
	if err != nil {
		fmt.Println("Error parsing response")
		return
	}

	fmt.Println(body.Content)

	if resp.StatusCode() != 200 {
		return
	}

	// generate token file in /etc/luidium/cli_token.toml
	configName := "cli_token.toml"
	configContent := fmt.Sprintf(`cli_token = "%s"`, token)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}
	configDir := filepath.Join(homeDir, ".luidium")
	configPath := filepath.Join(configDir, configName)

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Println("Error creating config directory:", err)
		return
	}

	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error creating config file:", err)
		return
	}
	defer file.Close()

	file.WriteString(configContent)
}

func init() {
	rootCmd.AddCommand(certify)
}
