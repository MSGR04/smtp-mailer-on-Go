package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MailRequest struct {
	From    string `form:"from" binding:"required,email"`
	To      string `form:"to"`
	Subject string `form:"subject" binding:"required"`
	Body    string `form:"body" binding:"required"`
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "form.html", nil)
	})

	router.POST("/send", func(c *gin.Context) {
		var req MailRequest
		if err := c.ShouldBind(&req); err != nil {
			c.HTML(http.StatusBadRequest, "form.html", gin.H{"error": err.Error()})
			return
		}

		var recipients []string

		file, err := c.FormFile("recipientsFile")
		if err == nil {
			f, err := file.Open()
			if err != nil {
				c.HTML(http.StatusBadRequest, "form.html", gin.H{"error": "Ошибка открытия файла: " + err.Error()})
				return
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				email := strings.TrimSpace(scanner.Text())
				if email != "" {
					recipients = append(recipients, email)
				}
			}
			if err := scanner.Err(); err != nil {
				c.HTML(http.StatusBadRequest, "form.html", gin.H{"error": "Ошибка чтения файла: " + err.Error()})
				return
			}
		} else {
			if req.To == "" {
				c.HTML(http.StatusBadRequest, "form.html", gin.H{"error": "Укажите хотя бы один email получателя или загрузите файл."})
				return
			}
			addresses := strings.Split(req.To, ",")
			for _, addr := range addresses {
				email := strings.TrimSpace(addr)
				if email != "" {
					recipients = append(recipients, email)
				}
			}
		}

		if len(recipients) == 0 {
			c.HTML(http.StatusBadRequest, "form.html", gin.H{"error": "Список получателей пуст."})
			return
		}

		var lastErr error
		for _, recipient := range recipients {
			msg := buildMessage(req, recipient)
			err := sendMail("localhost:25", req.From, recipient, msg)
			if err != nil {
				log.Printf("Ошибка отправки письма для %s: %v", recipient, err)
				lastErr = err
			} else {
				log.Printf("Письмо успешно отправлено для %s", recipient)
			}
			time.Sleep(5 * time.Second)
		}

		if lastErr != nil {
			c.HTML(http.StatusInternalServerError, "form.html", gin.H{"error": fmt.Sprintf("Ошибка отправки письма: %v", lastErr)})
		} else {
			c.HTML(http.StatusOK, "form.html", gin.H{"message": "Письма успешно отправлены!"})
		}
	})

	router.Run(":8080")
}

func buildMessage(req MailRequest, recipient string) string {
	date := time.Now().Format(time.RFC1123Z)
	message := fmt.Sprintf(
		"Date: %s\r\n"+
			"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n"+
			"%s",
		date, req.From, recipient, req.Subject, req.Body,
	)
	return message
}

func sendMail(addr, from, to, body string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	readResponse := func() {
		line, err := r.ReadString('\n')
		if err != nil {
			log.Printf("Ошибка при чтении ответа: %v", err)
		} else {
			log.Printf("<< %s", strings.TrimSpace(line))
		}
	}

	writeCommand := func(cmd string) {
		fmt.Fprintf(w, "%s\r\n", cmd)
		w.Flush()
		log.Printf(">> %s", cmd)
		readResponse()
	}

	readResponse()
	writeCommand("HELO localhost")
	writeCommand("MAIL FROM:<" + from + ">")
	writeCommand("RCPT TO:<" + to + ">")
	writeCommand("DATA")

	fmt.Fprintf(w, "%s\r\n.\r\n", body)
	w.Flush()
	readResponse()

	writeCommand("QUIT")
	return nil
}
