# go-generator-filesystem

## Overview

The go-generator-filesystem project is a tool designed to facilitate the conversion of text and binary files into byte arrays in the Go programming language (Golang). This enables users to reference these files directly within their Go code without needing to read them from the file system.

## Features

File Conversion: Convert both text and binary files into byte arrays.
Go Integration: Seamless integration with Go code, allowing for easy referencing of files.
Efficient Access: Avoids the need to read files from the file system at runtime, enhancing performance in certain scenarios.
Simplified Deployment: Eliminates the necessity of deploying separate files alongside the application executable, as the file contents are embedded within the code.

## How It Works

The go-generator-filesystem tool works by taking input files and generating Go code containing byte array representations of these files. This generated code can then be included in Go projects, allowing direct access to the file contents without requiring file I/O operations.

## Usage

* Install the go-generator-filesystem tool.
* Specify the files to be converted.
* Run the tool to generate Go code with byte array representations of the input files.
* Include the generated code in your Go project.
* Access the file contents directly within your Go code using the generated byte arrays.

## Install

``` bash
go install github.com/SebastianPozoga/go-generator-filesystem
```

## Example

Suppose you have a text file named "example.txt" containing some text data. Using go-generator-filesystem, you can convert this file into a byte array representation and include it in your Go project. Then, you can access the contents of "example.txt" within your Go code without needing to read it from the file system at runtime.

``` go
package main

import (
    "fmt"
    staticdata "github.com/username/yourproject/staticdata"
)

func main() {
    // Accessing the contents of the example.txt file
    fmt.Println(string(staticdata.ExampleTxt.Data))
}
```

``` bash
go-generator-filesystem --from=example/from --to=example/to --cache=example/cache  --package.prefix=github.com/user/repo/fs/public/
```

Find result in [example directory](https://github.com/SebastianPozoga/go-generator-filesystem/blob/main/example)


## Contribution

Contributions to the go-generator-filesystem project are welcome. Feel free to submit bug reports, feature requests, or pull requests through GitHub.

## License

This project is licensed under the MIT License. Feel free to use, modify, and distribute it according to the terms of the license.

## Author

This project is maintained by Sebastian Pozoga. For any inquiries or assistance, please contact via github.