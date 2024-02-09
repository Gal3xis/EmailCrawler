package main

import (
	"EmailCrawler/src/conf"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var Config = conf.Config{}

func main() {
	// Config einlesen
	err := Config.Read()
	if err != nil {
		return
	}

	// Verbindung zum IMAP-Server herstellen
	c, err := client.DialTLS(Config.Url+":"+strconv.Itoa(Config.Port), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Connected")
	defer c.Terminate()

	err = c.Login(Config.Username, Config.Password)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Logged in")
	defer c.Logout()

	for _, mailboxConfig := range Config.Mailboxes {
		err := Work(c, mailboxConfig)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
	}
}

func Work(client *client.Client, mailboxConfig conf.MailboxConfig) error {
	_, err := client.Select(mailboxConfig.Mailbox, false)
	if err != nil {
		return err // Verwenden Sie return statt log.Fatal
	}

	seqset := new(imap.SeqSet)
	seqset.Add("1:*")
	if err != nil {
		return err
	}
	items := []imap.FetchItem{imap.FetchItem("BODY[]"), imap.FetchEnvelope, imap.FetchFlags, imap.FetchInternalDate, imap.FetchRFC822Size, imap.FetchBody, imap.FetchBodyStructure}

	messages := make(chan *imap.Message)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(seqset, items, messages)
	}()

	for msg := range messages {
		section := &imap.BodySectionName{} // Keine spezifischen Teile angegeben, holt den kompletten Inhalt
		body := msg.GetBody(section)
		if body == nil {
			log.Println("Konnte den Nachrichtenkörper nicht abrufen")
			continue
		}
		content, err := io.ReadAll(body)
		if err != nil {
			log.Fatal("Fehler beim Lesen der Nachricht:", err)
		}

		filename := sanitizeSubject(msg.Envelope.Subject) + ".eml"
		filepath := filepath.Join(Config.SaveFolder, filename)
		file, err := os.Create(filepath)
		if err != nil {
			return err
		}
		file.WriteString(string(content))
		file.Close()
	}
	return nil
}

func sanitizeSubject(subject string) string {
	// Entfernen Sie alle unerwünschten Zeichen aus dem Betreff
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	safeSubject := reg.ReplaceAllString(subject, "_")

	// Begrenzen Sie die Länge des Dateinamens
	if len(safeSubject) > 50 {
		safeSubject = safeSubject[:50]
	}
	return safeSubject
}
