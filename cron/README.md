# Cron

The **cron** package is a thin wrapper around the `github.com/robfig/cron/v3` library for running functions on a cron schedule. 

```go
cron.Add("*/1 * * * *", func() {
    fmt.Printf("tick\n")
})

cron.Start()

<-waiter.Wait()
cron.Stop()
```
