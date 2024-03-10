package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	MemSize = 10 * 1024 * 1024 // 10 MB
	Value   = 0b00001111       // Ones and zeros to detect both possible changes
)

var (
	chatID string
	token  string
)

func init() {
	flag.StringVar(&chatID, "chat_id", "", "Telegram Chat ID")
	flag.StringVar(&token, "token", "", "Telegram Bot Token")
	flag.Parse()
}

type TelegramWriter struct {
	alternate io.Writer
}

func (t *TelegramWriter) Write(p []byte) (n int, err error) {
	if err := SendMessage(string(p)); err != nil {
		return t.alternate.Write(p)
	}
	return len(p), nil
}

func SendMessage(text string) error {
	msg := map[string]string{
		"chat_id": chatID,
		"text":    text,
	}
	// io.Reader from encoded to json msg
	data, _ := json.Marshal(msg)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("failed to send message")
	}
	return nil
}

func FileOutput() io.WriteCloser {
	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func Run() {
	log.Println("Starting...")

	// Create a empty bytes slice.
	data := make([]byte, MemSize)

	// Fill slice
	for i := range data {
		data[i] = Value
	}

	for i, a := range data {
		// If value by some reason (cosmos particular) changed, log it.
		if (a ^ Value) != 0 {
			log.Printf("CHANGE DETECTED: i=%d value=%08b expected=%08b", i, a, Value)
		}
		// Sleep for a minute.
		time.Sleep(time.Minute)
	}
}

func main() {
	// Output all logs into w
	w := FileOutput()
	defer w.Close()
	output := &TelegramWriter{
		alternate: w,
	}
	log.SetOutput(output)
	go Run()
	// fetch interrupt to garante correct file closing.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
}
