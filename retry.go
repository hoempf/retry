package retry

// If you want to use this, please copy it, don't import it. It is too small to
// warrant its own imported package. This is also why I use unexported names.

// Simple retry logic.
type retry struct {
	Max int // Max attempts to try.
	N   int // Attempts we tried until we failed or succeeded.
}

// retryFunc is the func we need to try (usually a closure).
type retryFunc func() (retry bool, err error)

// Do try retryFunc and repeat until retries are exhausted or retry == false.
func (r *retry) Do(f retryFunc) (err error) {
	cont := true
	for r.N < r.Max && cont {
		cont, err = f()
		r.N++
	}
	return
}
