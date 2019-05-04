package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	flag "github.com/spf13/pflag"
)

var limit int

func init() {
	flag.IntVarP(&limit, "limit", "l", 10, "Enter the limit")
}

func strToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func getNumbers(question string) (int, int, error) {
	questionVars := strings.Split(question, "+")
	n1, err := strToInt(questionVars[0])
	if err != nil {
		log.Fatalf("Can't Convert Number %s", questionVars[0])
		return 0, 0, err
	}
	n2, err := strToInt(questionVars[1])
	if err != nil {
		log.Fatalf("Can't Convert Number %s", questionVars[1])
		return 0, 0, err
	}
	return n1, n2, err
}

func userInput(reader *csv.Reader, totalQuestion *int, score *int, quit chan bool) {

	line, err := reader.Read()
	if err == io.EOF {
		quit <- true
		return
	} else if err != nil {
		log.Fatal(err)
	}
	a, b, err := getNumbers(line[0])
	if err != nil {
		log.Fatalf("Invalid question %s", line[0])
	}
	ans, err := strToInt(line[1])
	if err != nil {
		log.Fatalf("Invalid answer %s", line[1])
	}
	fmt.Printf("%d + %d = ", a, b)
	ioReader := bufio.NewReader(os.Stdin)
	text, err := ioReader.ReadString('\n')
	if err != nil {
		log.Fatalf("Unable to receive input %s", err)
	}
	text = strings.TrimSpace(text)
	userAns, err := strToInt(text)
	if err != nil {
		log.Fatal("Invalid answer format")
	}
	(*totalQuestion)++
	correctAns := (userAns == ans)
	if correctAns {
		(*score)++
	}
	quit <- false
}

func timer(seconds int, timeout chan bool) {
	time.Sleep(time.Duration(seconds) * time.Second)
	timeout <- true
}

func dumpy(line chan bool) {
	time.Sleep(1 * time.Second)
	line <- true
}

func startQuiz() {
	csvFile, _ := os.Open("./problems.csv")
	score, totalQuestion := 0, 0
	reader := csv.NewReader(bufio.NewReader(csvFile))
	defer csvFile.Close()
	timeout := make(chan bool, 1)
	cancelChan := make(chan bool, 1)
	go timer(limit, timeout)
	go userInput(reader, &totalQuestion, &score, cancelChan)
Loop:
	for {
		select {
		case quit := <-cancelChan:
			if quit {
				break Loop
			} 
			userInput(reader, &totalQuestion, &score, cancelChan)
		case <-timeout:
			break Loop
		}
	}

	fmt.Printf("Your score is %d/%d \n", score, totalQuestion)

}

var rootCmd = &cobra.Command{
	Use:   "Quiz",
	Short: "Starts the quiz",
	Long:  `Answer the questions as fast as you can`,
	Run: func(cmd *cobra.Command, args []string) {
		startQuiz()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
