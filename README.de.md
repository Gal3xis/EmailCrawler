# Email Crawler

Der Email Crawler ist ein Tool, das dazu dient, E-Mails von einem E-Mail-Server herunterzuladen und anschließend vom Server zu löschen. Dies hilft dabei, Speicherplatz auf dem Server zu reduzieren, ohne die Notwendigkeit, E-Mails manuell einzeln löschen zu müssen.

## Inhaltsverzeichnis

- [Email Crawler](#email-crawler)
  - [Inhaltsverzeichnis](#inhaltsverzeichnis)
  - [Installation](#installation)
  - [Konfiguration](#konfiguration)
    - [Server-Informationen](#server-informationen)
    - [Anmeldedaten](#anmeldedaten)
    - [Mailbox-Konfiguration](#mailbox-konfiguration)
  - [Verwendung](#verwendung)

## Installation

1. Laden Sie das benötigte Binary für Windows oder Linux herunter. Diese finden Sie unter den [Releases](https://github.com/Gal3xis/EmailCrawler/releases).
2. Neben der Binary wird eine Konfigurationsdatei (`emailCrawlerConfig.conf`) benötigt. Diese muss sich im selben Verzeichnis wie das Binary befinden. Ein Template für diese Datei ist ebenfalls unter den [Releases](https://github.com/Gal3xis/EmailCrawler/releases) zu finden.

## Konfiguration

### Server-Informationen

Die Konfiguration der Verbindung erfolgt im Block `Connection`. Zur Verbindung mit dem Server werden die IMAP-Adresse und der Port benötigt. Die folgende Liste zeigt die Informationen für die wichtigsten Server:

<details>
  <summary>IMAP-Servereinstellungen für verschiedene E-Mail-Dienste (Klicken zum Anzeigen)</summary>

- **AOL Mail**
  - **Serveradresse**: imap.aol.com
  - **Port**: 993
  - **Verschlüsselung**: SSL

- **Gmail**
  - **Serveradresse**: imap.gmail.com
  - **Port**: 993
  - **Verschlüsselung**: SSL

- **GMX Mail**
  - **Serveradresse**: imap.gmx.com
  - **Port**: 993
  - **Verschlüsselung**: SSL

- **iCloud Mail**
  - **Serveradresse**: imap.mail.me.com
  - **Port**: 993
  - **Verschlüsselung**: SSL

- **Mail.com**
  - **Serveradresse**: imap.mail.com
  - **Port**: 993
  - **Verschlüsselung**: SSL

- **Outlook.com / Hotmail**
  - **Serveradresse**: imap-mail.outlook.com
  - **Port**: 993
  - **Verschlüsselung**: SSL

- **Posteo**
  - **Serveradresse**: posteo.de
  - **Port**: 993
  - **Verschlüsselung**: SSL

- **Web.de**
  - **Serveradresse**: imap.web.de
  - **Port**: 993
  - **Verschlüsselung**: SSL

- **Yahoo Mail**
  - **Serveradresse**: imap.mail.yahoo.com
  - **Port**: 993
  - **Verschlüsselung**: SSL

- **Zoho Mail**
  - **Serveradresse**: imap.zoho.com
  - **Port**: 993
  - **Verschlüsselung**: SSL
</details>


### Anmeldedaten

Neben den Server-Informationen müssen auch die Anmeldedaten für den Zugang zum E-Mail-Server hinterlegt werden.

Die Konfiguration der Verbindung sieht wie folgt aus:

```
[Connection]
Url = ...
Port = ...
Username = ...
Password = ...
```
### Mailbox-Konfiguration

Die Konfiguration der Mailboxen erfolgt in individuellen Unterblöcken. Für jede Mailbox, aus der E-Mails heruntergeladen werden sollen, muss ein eigener Block konfiguriert werden.

## Verwendung

Um eine Liste aller verfügbaren Mailboxen anzuzeigen, verwenden Sie den folgenden Befehl. Stellen Sie sicher, dass die Verbindung zuvor erfolgreich konfiguriert wurde.

**Windows:**

```powershell
.\EmailCrawler.exe list
```

**Linux:**
```bash
./emailCrawler list
```

Folgende Konfigurationseinstellungen sind möglich:

- **MailOffset**: Die Anzahl der zuletzt eingegangenen E-Mails, die ignoriert werden sollen. Dies verhindert, dass die neuesten X E-Mails heruntergeladen und gelöscht werden.
- **SaveFolder**: Der absolute Pfad zum Speicherort, an dem alle E-Mails gesichert werden sollen.
- **SavingStructure**: Der relative Speicherpfad für jede einzelne E-Mail. Hierbei können folgende Variablen verwendet werden:
  - `%_FROM_%`: Der Absender der E-Mail.
  - `%_SUBJECT_%`: Der Betreff der E-Mail.
  - `%_DATE_%`: Der Zeitstempel der E-Mail.
- **DeleteMails**: Muss auf `true` gesetzt werden, um das Löschen von E-Mails zu ermöglichen.
- **MinAgeInDaysToDelete**: Das Mindestalter in Tagen, das E-Mails erreicht haben müssen, bevor sie gelöscht werden. Dies soll verhindern, dass kürzlich empfangene E-Mails gelöscht werden.
- **MinEmailsToKeep**: Die minimale Anzahl von E-Mails, die im Posteingang verbleiben sollen, um sicherzustellen, dass die letzten X E-Mails nicht gelöscht werden.

Eine Mailbox-Konfiguration könnte wie folgt aussehen:

```
[INBOX]
MailOffset = 0
SaveFolder = /path/to/SaveFolder
SavingStructure = %_FROM_%/%_SUBJECT_%_%_DATE_%
DeleteMails = true
MinAgeInDaysToDelete = 60
MinEmailsToKeep = 30
```








