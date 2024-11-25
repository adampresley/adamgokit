# Waiter

Provides a function to wait on a channel for an interrupt signal. Here is a sample usage:

```go
<-waiter.Wait()
```

That code will wait for an interrupt signal, such as `CTRL+C`.
