package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/GhostNet-Dev/glambda/gtoken"
	"github.com/GhostNet-Dev/glambda/lexer"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lexer := lexer.NewLexer(line)

		for tok := lexer.NextTokenMake(); tok.Type != gtoken.EOF; tok = lexer.NextTokenMake() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
