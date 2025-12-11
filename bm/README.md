# bm - CLI Bookmark Manager

`bm` is a command-line interface (CLI) tool for managing bookmarks. It allows you to add, find, and read bookmarks from the command line.

## Installation

```sh
go install github.com/bart-jaskulski/toolbox/bm@latest
```

## Usage

### Adding a Bookmark

To add a new bookmark, use the `add` command:

```
bm add <url>
```

Alternatively, you can provide the URL through standard input:

```
echo "https://example.com" | bm add
```

### Finding Bookmarks

To find bookmarks based on a regular expression pattern, use the `find` command:

```
bm find <pattern>
```

If no pattern is provided, all bookmarks will be listed.

### Reading a Bookmark

To open a bookmark in your default browser, use the `read` command:

```
bm read <url|pattern>
```

If you provide a pattern instead of a URL, the tool will search for matching bookmarks and let you choose one to open.

## Examples

1. Add a new bookmark:
   ```
   bm add https://www.example.com
   ```

2. Find bookmarks containing the word "google":
   ```
   bm find google
   ```

3. Open a bookmark by URL:
   ```
   bm read https://www.google.com
   ```

4. Open a bookmark by pattern:
   ```
   bm read google
   ```

## Contributing

If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request on the [GitHub repository](https://github.com/bart-jaskulski/toolbox/bm).

## License

This project is licensed under the [MIT License](LICENSE).
