# UF2C v0.0.1

_A tool made by lazy people, for lazy people._

This Go package provides a command-line tool to assist with Git workflows by leveraging the OpenAI API. It can generate commit messages and progress updates based on the changes in your Git working directory.

## Features

-   **Generate Commit Messages:** Automatically create concise and relevant commit messages based on the `git diff`.
-   **Generate Progress Updates:** Produce clear and simple summaries of your daily work, including what you accomplished and any challenges faced. The output is formatted for easy understanding.
-   **Clipboard Integration (macOS):** Optionally copies the generated output to your clipboard for quick use.

## Prerequisites

-   **Go Installation:** Ensure you have Go installed on your system. You can download it from [https://go.dev/dl/](https://go.dev/dl/).
-   **Git:** Git must be initialized in your project directory.
-   **.env File:** You need to create a `.env` file in the root of your project with your OpenAI API key:

    ```
    OPENAI_API_KEY=your_openai_api_key_here
    ```

    Replace `your_openai_api_key_here` with your actual OpenAI API key. You can obtain one from [https://platform.openai.com/api-keys](https://platform.openai.com/api-keys).

## Installation

1.  **Clone the repository (if applicable):**

    ```bash
    git clone <repository_url>
    cd <repository_directory>
    ```

2.  **Build the Go package:**

    ```bash
    go build -o uf2c main.go
    ```

    This will create an executable file named `uf2c` in your current directory.

## Usage

Navigate to your Git repository in the terminal and run the `uf2c` executable with the desired flags:

### Show help
```bash
./uf2c -h
```

### Generate Commit Message

```bash
./uf2c -cm
```

### Generate Progress Update

```bash
./uf2c -pu
```

## Hey! Consider supporting me in my FOSS journey. I plan to keep on making free tools like this, and well, I need the $$ to buy me some good ol' coffee~ Thanksss

If you like my work, consider supporting me on [Buy Me A Coffee](https://buymeacoffee.com/fossoctopus). Your support helps me continue developing open-source projects and bringing more value to the community! - Foss

[![Buy Me A Coffee](https://img.shields.io/badge/Support_Me-Buy_Me_A_Coffee-yellow?logo=buymeacoffee)](https://buymeacoffee.com/fossoctopus)
