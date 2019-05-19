# A simple retry logic for Go
This is a small experiment of mine to manage function retries in Go. I use it in
a few places by copying this. Don't import this as a library. It is very small
and probably not updated that much ;)

# Use

See example in tests:

```go
func Example_retry_Do() {
	max := 10
	i := 0

	try := &retry{Max: max}
	err := try.Do(func() (bool, error) {
		if i++; i > max {
			return false, nil // Finished.
		}
		fmt.Println(i)
		return true, nil // Repeat.
	})

	// You can capture the error if any and act on it. The closure you pass
	// should decide if we should retry, err here is just the last error after
	// we reached max retries.
	fmt.Println(err)
	// And also the number of executions.
	fmt.Println(try.N)

	//Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// 10
	// <nil>
	// 10
}
```
# Author
Mathias Seiler, mathias.seiler@mironet.ch