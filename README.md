# Repo Copy CLI 📦

Repo Copy CLI is a powerful command-line tool designed to streamline the process of extracting code from Git repositories for developers working with large language models (LLMs) like ChatGPT. This tool ensures that all relevant code is collected into a single file, adhering to `.gitignore` settings, and provides extensive statistics to aid in data understanding and model training.

## Key Features 🌟

- **Efficient Code Extraction 🚀**: Extracts code while ignoring files specified in `.gitignore`.
- **Single File Output 📄**: Outputs all extracted code into a `codebase.txt`, facilitating easier context management for LLMs.
- **Codebase Statistics 📊**: Provides detailed insights, including:
  - Total number of files 🗂️
  - Total number of lines 📃
  - Total number of words 📝
  - Total number of characters 🔠
  - Distribution of programming languages by file extensions 💻
  - Total number of tokens 🔤
  - Estimated number of requests needed to process the codebase with LLMs (assuming 4,096 tokens per request) 🤖
- **Clipboard Integration 📋**: Automatically copies the `codebase.txt` contents to the clipboard.

## Installation 🛠️

### macOS

1. Download the `repo-copy` binary for macOS from the [GitHub Releases](https://github.com/MalteBoehm/repo-copy/releases) page.
2. Open a terminal and navigate to the directory where you downloaded the binary.
3. Move the binary to a directory in your system's `PATH`, for example, `/usr/local/bin`:
   sudo mv repo-copy /usr/local/bin/
4. Make the binary executable:
   sudo chmod +x /usr/local/bin/repo-copy
5. You can now run the CLI from any directory by typing `repo-copy` in the terminal.

### Windows

1. Download the `repo-copy.exe` binary for Windows from the [GitHub Releases](https://github.com/MalteBoehme/repo-copy/releases) page.
2. Move the `repo-copy.exe` file to a directory in your system's `PATH`, e.g., `C:\Windows\System32`.
3. Open a command prompt and navigate to the directory where you moved the binary.
4. You can now run the CLI from any directory by typing `repo-copy` in the command prompt.

### Linux

1. Download the `repo-copy` binary for Linux from the [GitHub Releases](https://github.com/MalteBoehme/repo-copy/releases) page.
2. Open a terminal and navigate to the directory where you downloaded the binary.
3. Move the binary to a directory in your system's `PATH`, for example, `/usr/local/bin`:
   sudo mv repo-copy /usr/local/bin/
4. Make the binary executable:
   sudo chmod +x /usr/local/bin/repo-copy
5. You can now run the CLI from any directory by typing `repo-copy` in the terminal.

## Usage 📘

1. Open a terminal (or command prompt on Windows) and navigate to the root directory of your Git repository.
2. Run the following command:
   repo-copy
3. The CLI will process the files, ignoring those specified in the `.gitignore`. It then generates a `codebase.txt` file in the current directory and copies its contents to the clipboard.
4. After processing, the CLI displays the detailed statistics about the codebase.

## License 📝

This project is licensed under the [MIT License](LICENSE). See the `LICENSE` file for more details.

