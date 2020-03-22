package main

func main()  {
	bc := NewBlockChain("1LExBMxTqpjyMv358NKbYVPhdVC3EhDH4r")
	cli := CLI{bc:bc}
	cli.Run()
}