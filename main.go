package main

func main()  {
	bc := NewBlockChain("1DERDjJ4eG1ReGfcWdCEh2orNiTKVnm3Pw")
	cli := CLI{bc:bc}
	cli.Run()
}