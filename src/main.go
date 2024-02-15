package main

import (
	"EmailCrawler/src/conf"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var Config = conf.Config{}

func main() {
	// Config einlesen
	err := Config.Read()
	if err != nil {
		log.Println(err.Error())
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
	err = seqset.Add("1:*")
	if err != nil {
		return err
	}
	items := []imap.FetchItem{imap.FetchItem("BODY[]"), imap.FetchEnvelope, imap.FetchFlags, imap.FetchInternalDate, imap.FetchRFC822Size, imap.FetchBody, imap.FetchBodyStructure, imap.FetchUid}

	channel := make(chan *imap.Message)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(seqset, items, channel)
	}()
	var messages []*imap.Message
	for msg := range channel {
		messages = append(messages, msg)
	}
	var toDelete []uint32
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Envelope.Date.After(messages[j].Envelope.Date)
	})

	if mailboxConfig.MailOffset > len(messages) {
		return nil
	}
	messages = messages[mailboxConfig.MailOffset-1:]

	for i, msg := range messages {
		section := &imap.BodySectionName{} // Keine spezifischen Teile angegeben, holt den kompletten Inhalt
		body := msg.GetBody(section)
		if body == nil {
			log.Println("Konnte den Nachrichtenkörper nicht abrufen")
			continue
		}
		content, err := io.ReadAll(body)
		if err != nil {
			log.Println("Fehler beim Lesen der Nachricht:", err)
		}
		filename := getRelSavePath(mailboxConfig, msg)
		filePath := filepath.Join(mailboxConfig.SaveFolder, filename) + ".eml"
		if !fileExists(filePath) {
			dirPath := filepath.Dir(filePath)
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				log.Print(err.Error())
				continue
			}
			file, err := os.Create(filePath)
			if err != nil {
				file.Close()
				log.Print(err.Error())
				continue
			}
			file.WriteString(string(content))
			//setCreationTime(filePath, msg.InternalDate)
			file.Close()
		}

		ageInDays := int(time.Since(msg.InternalDate).Hours() / (float64(24)))
		if i < mailboxConfig.MinEmailsToKeep || ageInDays < mailboxConfig.MinAgeInDays {
			continue
		}
		toDelete = append(toDelete, msg.Uid)
	}
	if mailboxConfig.DeleteMails {
		DeleteMails(toDelete, client, mailboxConfig)
	}
	return nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil { // Kein Fehler, Datei existiert
		return true
	}
	if os.IsNotExist(err) { // Spezifischer Fehler für nicht existierende Datei
		return false
	}
	return false // Andere Fehler könnten auftreten, z.B. Berechtigungsprobleme
}

func sanitizePathSegment(segment string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, segment)
}

func getRelSavePath(mailboxConfig conf.MailboxConfig, msg *imap.Message) string {
	path := mailboxConfig.SavingStructure
	path = strings.ReplaceAll(path, "%_FROM_%", sanitizePathSegment(msg.Envelope.From[0].Address()))
	path = strings.ReplaceAll(path, "%_SUBJECT_%", sanitizePathSegment(getSubject(msg)))
	path = strings.ReplaceAll(path, "%_DATE_%", sanitizePathSegment(msg.InternalDate.Format("02-01-2006T15-04-05")))
	return path
}

func getSubject(msg *imap.Message) string {
	if msg.Envelope.Subject == "" {
		return "No Subject"
	}
	return msg.Envelope.Subject
}

func DeleteMails(uids []uint32, client *client.Client, mailboxConfig conf.MailboxConfig) {
	if uids == nil {
		return
	}
	_, err := client.Select(mailboxConfig.Mailbox, false)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}
	if err := client.UidStore(seqSet, item, flags, nil); err != nil {
		log.Printf("Fehler beim Markieren der E-Mails als gelöscht: %v", err)
		return
	}

	if err := client.Expunge(nil); err != nil {
		log.Printf("Fehler beim Ausführen von Expunge: %v", err)
		return
	}

	log.Println("E-Mails erfolgreich gelöscht")
}
