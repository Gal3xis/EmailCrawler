[![README-de](https://img.shields.io/badge/README-Deutsch-red)](README.de.md)


# Email Crawler

The Email Crawler is a tool designed to download emails from an email server and then delete them from the server. This helps to reduce storage space on the server without the need to manually delete emails one by one.

## Table of Contents

- [Email Crawler](#email-crawler)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Configuration](#configuration)
    - [Server Information](#server-information)
    - [Login Credentials](#login-credentials)
    - [Mailbox Configuration](#mailbox-configuration)
  - [Usage](#usage)

## Installation

1. Download the required binary for Windows or Linux. These can be found under the [Releases](https://github.com/Gal3xis/EmailCrawler/releases).
2. In addition to the binary, a configuration file (`emailCrawlerConfig.conf`) is required. This must be located in the same directory as the binary. A template for this file can also be found under the [Releases](https://github.com/Gal3xis/EmailCrawler/releases).

## Configuration

### Server Information

The configuration of the connection is done in the `Connection` block. To connect to the server, the IMAP address and port are needed. The following list shows the information for the most important servers:

<details>
  <summary>IMAP server settings for various email services (Click to show)</summary>

- **AOL Mail**
  - **Server address**: imap.aol.com
  - **Port**: 993
  - **Encryption**: SSL

- **Gmail**
  - **Server address**: imap.gmail.com
  - **Port**: 993
  - **Encryption**: SSL

- **GMX Mail**
  - **Server address**: imap.gmx.com
  - **Port**: 993
  - **Encryption**: SSL

- **iCloud Mail**
  - **Server address**: imap.mail.me.com
  - **Port**: 993
  - **Encryption**: SSL

- **Mail.com**
  - **Server address**: imap.mail.com
  - **Port**: 993
  - **Encryption**: SSL

- **Outlook.com / Hotmail**
  - **Server address**: imap-mail.outlook.com
  - **Port**: 993
  - **Encryption**: SSL

- **Posteo**
  - **Server address**: posteo.de
  - **Port**: 993
  - **Encryption**: SSL

- **Web.de**
  - **Server address**: imap.web.de
  - **Port**: 993
  - **Encryption**: SSL

- **Yahoo Mail**
  - **Server address**: imap.mail.yahoo.com
  - **Port**: 993
  - **Encryption**: SSL

- **Zoho Mail**
  - **Server address**: imap.zoho.com
  - **Port**: 993
  - **Encryption**: SSL
</details>


### Login Credentials

In addition to the server information, the login credentials for access to the email server must also be provided.

The connection configuration looks as follows:

```
[Connection]
Url = ...
Port = ...
Username = ...
Password = ...
```


### Mailbox Configuration

The configuration of the mailboxes is done in individual sub-blocks. A separate block must be configured for each mailbox from which emails are to be downloaded.

## Usage

To display a list of all available mailboxes, use the following command. Make sure the connection has been successfully configured beforehand.

**Windows:**

```powershell
.\EmailCrawler.exe list
```

**Linux:**
```bash
./emailCrawler list
```

The following configuration settings are possible:

- **MailOffset**: The number of the most recently received emails that should be ignored. This prevents the latest X emails from being downloaded and deleted.
- **SaveFolder**: The absolute path to the storage location where all emails should be saved.
- **SavingStructure**: The relative storage path for each individual email. The following variables can be used:
  - `%_FROM_%`: The sender of the email.
  - `%_SUBJECT_%`: The subject of the email.
  - `%_DATE_%`: The timestamp of the email.
- **DeleteMails**: Must be set to `true` to enable the deletion of emails.
- **MinAgeInDaysToDelete**: The minimum age in days that emails must have reached before they can be deleted. This is to prevent recently received emails from being deleted.
- **MinEmailsToKeep**: The minimum number of emails that should remain in the inbox to ensure that the last X emails are not deleted.

A mailbox configuration might look like this:

```makefile
[INBOX]
MailOffset = 0
SaveFolder = /path/to/saveFolder
SavingStructure = %_FROM_%/%_SUBJECT_%_%_DATE_%
DeleteMails = true
MinAgeInDaysToDelete = 60
MinEmailsToKeep = 30
