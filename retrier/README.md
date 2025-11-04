## Retrier

This package has methods to retry blocks of code. It offers options to configure delays and jitter.

```go
err := retrier.Retry(
	func() error {
		// Do something here that could fail but you want to retry.
		// Returning nil means it was successful. Returning an 
		// error will cause this function to retry up to Max Attempts.
	}, 
	retrier.WithMaxAttempts(4),
	retrier.WithDelay(time.Second * 3),
	retrier.WithMaxJitter(200 * time.Millisecond),
)

if err != nil {
	panic(err)
}
```

Default values are:

- Delay: 2 seconds. Each attempt multiples this value
- MaxAttempts: 3
- MaxJitter: 300 milliseconds

