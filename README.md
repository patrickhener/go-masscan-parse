# go-masscan-parse

Easy tool to parse masscan xml results.

## Installation

```bash
go install github.com/patrickhener/go-masscan-parse@latest
```

## Usage

```bash
go-masscan-parse <inputfile.xml>
```

## Output

This tool will write a file called `output.txt`.

It will contain the results in this format:

```
IP:PORT
```