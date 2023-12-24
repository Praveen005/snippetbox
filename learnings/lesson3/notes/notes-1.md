## Fprintf:

Purpose: 
- Formats a string using placeholders (format specifiers) and writes the formatted output to a specified writer.

- Offers precise control over how data is presented in text-based output.

Syntax:
```
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
```
Parameters:
- ```w```: An ```io.Writer``` interface that represents the destination where the formatted output will be written. It can be a file, a network socket, a console, or any other object that implements the ```io.Writer``` interface.

    we did,
    ```
    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
    ```
    means in place of `io.writer` and we passed it our `http.ResponseWriter` object instead — and it worked fine.This happened beacuse
    `http.ResponseWriter` object satisfies the interface as it has a `w.Write()` method.


- ```format```: A string containing text and format specifiers (like ```%d```, ```%s```, ```%f```, etc.) that define how the values in ```a``` will be formatted.

- ```a ...interface{}```: A variadic argument list containing the values to be formatted and inserted into the format string.

Return Values:
- ```n```: The number of bytes written to the writer.

- ```err```: An error value, if any, encountered during the writing process.


## strconv.Atoi()

- Converts a string representation of an integer in base 10 (decimal) to an actual integer value.

Key Points:
- Base 10 Only: It specifically handles base 10 numbers. For other bases, use ```strconv.ParseInt()```.

- Error Handling: Always check the returned error to ensure successful conversion since invalid input can cause a panic.

- Whitespace: It ignores leading and trailing whitespace in the string.


## Key differences between ```fmt``` and ```log``` packages in Golang:

Purpose:

- log: Primary focus on logging events and messages for debugging, monitoring, and analysis.

- fmt: General-purpose formatting and printing of text output for user interaction and data presentation.

Common Use Cases:

- log:
    - Tracking application events and errors

    - Debugging issues

    - Monitoring application health

    - Generating logs for analysis

- fmt:
    - Printing user-facing messages

    - Formatting data for output

    - Creating structured text reports

    - Generating custom output formats

Key Points:

- Logging Levels: `log` offers different logging levels (e.g., Info, Debug, Error) for filtering and prioritizing messages.

- Timestamps: `log` automatically adds timestamps to log entries for chronological tracking.

- Concurrent Safety: `log` is thread-safe, ensuring proper output in multi-threaded environments.

- Custom Output: `log` allows customization of output format and destination.

- Destination: `fmt` primarily writes to standard output (stdout), while log defaults to standard error (stderr).


