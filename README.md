# Worker Pool

## install
go get github.com/whyming/workerpool

## Demo

### Basic
```go
func main() {
	p := workerpool.NewPool(4, 10)
	go request1(p)
	go request2(p)
	time.Sleep(time.Second)
}

func request1(p *workerpool.Pool) {
	for i := 0; i < 5; i++ {
		i := i
		p.AddJob(context.Background(), func() error {
			fmt.Println("request-1", i)
			return nil
		}, nil)
	}
	fmt.Println("request-1 DONE!")
}

func request2(p *workerpool.Pool) {
	for i := 5; i < 10; i++ {
		i := i
		p.AddJob(context.Background(), func() error {
			fmt.Println("request-2", i)
			return nil
		}, nil)
	}
	fmt.Println("request-2 DONE!")
}
```

### use group and wait
```go
var p = workerpool.NewPool(4, 10)

func main() {
	go request1(p)
	go request2(p)
	time.Sleep(time.Second)
}

func request1(p *workerpool.Pool) {
	g := p.NewGroup()
	for i := 0; i < 5; i++ {
		i := i
		g.AddJob(context.Background(), func() error {
			fmt.Println("request-1", i)
			return nil
		})
	}
	g.Done()
	fmt.Println("request-1 DONE!")
}

func request2(p *workerpool.Pool) {
	g := p.NewGroup()
	for i := 5; i < 10; i++ {
		i := i
		g.AddJob(context.Background(), func() error {
			fmt.Println("request-2", i)
			return nil
		})
	}
	g.Done()
	fmt.Println("request-2 DONE!")
}
```
