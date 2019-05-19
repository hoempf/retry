package retry

import (
	"fmt"
	"testing"
)

type executor struct {
	N     int // Reporting number of executions.
	Max   int // Retry this many times before reporting error.
	Retry int // Retry this many times before reporting success.
}

// Do does simulate execution of "something" Max times and reports executions to
// N so we can test if the retry struct is doing its job.
func (e *executor) Do() (retry bool, err error) {
	if e.N+1 > e.Max {
		return false, fmt.Errorf("max retries reached") // Failure.
	}
	if e.N+1 > e.Retry {
		return false, nil // Success.
	}
	e.N++
	return true, nil
}

// Return a new executor.
func newExecutor(max, retry int) *executor {
	e := &executor{
		Max:   max,
		Retry: retry,
	}
	return e
}

// Test a few retry cases.
func Test_retry_Do(t *testing.T) {
	type args struct {
		f retryFunc
	}
	tests := []struct {
		name    string
		r       *retry
		args    args
		wantErr bool
	}{
		{
			name: "Test 0",
			r:    &retry{Max: 0},
			args: args{
				f: newExecutor(0, 2).Do,
			},
			wantErr: false,
		},
		{
			name: "Test 1 which errors",
			r:    &retry{Max: 1},
			args: args{
				f: newExecutor(0, 2).Do,
			},
			wantErr: true,
		},
		{
			name: "Test 3 retries max, success after 2",
			r:    &retry{Max: 3},
			args: args{
				f: newExecutor(3, 2).Do,
			},
			wantErr: false,
		},
		{
			name: "Test 3 retries max, error after 2",
			r:    &retry{Max: 3},
			args: args{
				f: newExecutor(2, 3).Do,
			},
			wantErr: true,
		},
		{
			name: "Test 1000000 retries max, error after 100000",
			r:    &retry{Max: 1000000},
			args: args{
				f: newExecutor(100000, 100001).Do,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Do(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("retry.Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
