package main

func main()  {
	bc := NewBlockChain("kk")
	cli := CLI{bc:bc}
	cli.Run()
}