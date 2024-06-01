package cmd

type LuidiumConfig struct {
	BlockName            string   `json:"block_name"`
	Framework            string   `json:"framework"`
	PortBinding          string   `json:"port_binding"`
	VolumeBinding        string   `json:"volume_binding"`
	EnvironmentVariables []string `json:"environment_variables"`
	IgnoreFiles          []string `json:"ignore_files"`
}

type RunByCliRequest struct {
	Token     string `json:"token"`
	AppName   string `json:"app_name"`
	Version   string `json:"version"`
	BlockName string `json:"block_name"`
}

type StorageRequest struct {
	Action    string `json:"action"`
	Bucket    string `json:"bucket"`
	ObjectKey string `json:"object_key"`
}

type SingleStringResponse struct {
	Content string `json:"content"`
}

const (
	StorageActionListObject               = "list-object"
	StorageActionDeleteFolderExceptConfig = "delete-folder-except-config"
)
