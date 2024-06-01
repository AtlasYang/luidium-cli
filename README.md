# Luidium CLI

Luidium CLI is a command-line interface tool that allows you to interact with [Luidium Cloud Service](https://app.luidium.com). It provides a convenient way to authenticate, initialize projects, upload files, and deploy blocks.

## Installation

To install Luidium CLI, follow these steps:

1. Download the appropriate version of Luidium CLI for your operating system from the [releases page](https://github.com/AtlasYang/luidium-cli/releases).

2. Extract the downloaded archive.

3. Open a terminal and navigate to the extracted directory.

4. Run the setup script:

   - For Linux and macOS: `./setup.sh`
   - For Windows: Double-click on `setup.bat`

5. The Luidium CLI will be installed and ready to use.

## User Guide

### Authentication

Before using Luidium CLI, you need to authenticate with your Luidium Cloud Service account. Follow these steps:

1. Log in to your Luidium Cloud Service account.

2. Go to Settings > API.

3. Claim and Copy your CLI token.

4. Open a terminal and run the following command:
   ```
   luidium certify <CLI_TOKEN>
   ```
   Replace `<CLI_TOKEN>` with the token you copied from the Luidium Cloud Service.

### Initializing a Project

To initialize a new project with Luidium CLI, use the `luidium init` command:

```
luidium init <applicationname/version/blockname>
```

Replace `<applicationname/version/blockname>` with the appropriate values for your project. This command will pull the necessary code from the server and create a basic configuration for your project.

### Uploading Files

To upload files to your Luidium project, use the `luidium upload` command:

```
luidium upload <path>
```

Replace `<path>` with the path to the directory containing the files you want to upload. If you want to upload all files in the current directory, use `.` as the path.

### Deploying a Block

To deploy a block to your Luidium project, use the `luidium deploy` command:

```
luidium deploy <path>
```

Replace `<path>` with the path to the directory containing the files you want to deploy. This command will upload the files and add the block to the deployment queue. The block will be deployed automatically.

## Configuration File

Luidium CLI uses a configuration file named `luidium.toml` to store project-specific settings. You can download a template configuration file from the Application Dashboard.

Here's an explanation of the options in the configuration file:

- `application`, `version`, `name`: These fields are used to configure the block and should not be modified.

- `use_gitignore`: If this option is set to `true`, Luidium CLI will only upload files specified in the `.gitignore` file. If set to `false`, all files in the project directory will be uploaded.

- `use_dotenv`: If this option is set to `true`, Luidium CLI will load environment variables from a `.env` file located in the specified path (default is the project root). These variables will override the environment variables set in the Application Dashboard.

Make sure to place the `luidium.toml` file in the root directory of your project.

## Support

If you encounter any issues or have questions regarding Luidium CLI, please contact me at atlas@lighterlinks.io.
