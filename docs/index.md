
_DeepSea_ phishing gear aims to help RTOs and pentesters with the delivery of opsec-tight, flexible email phishing campaigns carried out on the inside of a perimeter. 

## Goals:

* Operate with a minimal footprint deep inside enterprises (Internal phish delivery).
* Seamlessly operate with external and internal mail providers (e.g. O365, Gmail, on-prem mail servers)
* Quickly re-target connectivity parameters.
* Flexibly add headers, targets, attachments
* Correctly format and inline email templates, images and multipart messages.
* Use content templates for personalization
* Account for various secure email communication parameters
* Clearly separate artifacts, mark databases and content delivery for multiple (parallel or sequential) phishing campaigns.


## Operating Instructions

DeepSea relies on directives specified in a `YAML` configuration file for it's runtime. Most of the directives are grouped by logical steps of executing mail delivery, loading marks, transforming messages, creating embedded resource URLs, and querying databases. Most of the directives can also be overriden on the command line if the operator finds this method more acceptable in the context.

### Some Quick Workflow Example:
To illustrate the workflow, here is an example of a campaign, follwed by a smaple configuration to satisfy the outlined scenario.

> We are mailing a phish via our Outlook Live account, connecting to it via a standard TLS SMTP port, with authentication (SMTP password is provided interactively). The message has been written in HTML, CSS inlined  and trasformed to contain a TXT counterpart. The contents of the email have been personalized from variables introduced by loading a list of marks (targeted users) from an external CSV file. For example, the First name and the Last names of a person. During preprocessing of the message DeepSea has generated a unique identifier for the user, and it was included in the content link. Here we also have chosen to attach a PDF document to the email, embed two logo files in the message, and add a set of SMTP headers to the mail envelope.

### Operational Workflow Example:

1. To begin working on the campaign we load marks from a CSV into a campaign database, taking directives from the campaign configuration file.

```sh
$ ./deepsea  storage load \
--config ./campaigns/campaign1/campaign.yml   
-d ./campaigns/campaign1/campaign.db   
-s ./campaigns/campaign1/campaign.csv 
```
2. Mail the campaign with all the parameters  specified in the configuration file which is depicted below.

```sh
$ ./deepsea  mailclient  --config ./conf/.deepsea.yaml
```

Here is the example of a mark records in the CSV file for DB import for reference. 

```sh
$cat campaigns/campaign_irc/campaign.csv

ident,email,firstname,lastname
<dynamic>,xxxx@xxxxxx.com,FName,LName
```

Here is the example of the configuration file for our scenario.

```yaml
mailclient:

  connection:

    #SMTPUser: "info"
    SMTPUser: "XXX@outlook.com"

    #SMTPServer: "smtp.gmail.com"
    SMTPServer: "smtp.office365.com"

    SMTPPort: 587
    #SMTPPort: 465
    TLS: "yes"

  message:

    Subject: "Here you go."
    ## Some providers, namely MSFT does not like to relay arbitrary emails.
    ## Make sure the "From" is your@outlook.com
    ## Or you get: `554 5.2.0 STOREDRV.Submission.Exception:SendAsDeniedException.MapiExceptionSendAsDenied`

    ## Google/Gmail is still ok with the following:
    #From: "Joe B <info@xxx.org>"

    ## Gather Marks from CSV import
    ## This is the most common use of the marks:
    To: "campaign1/db/deepsea.db"

    ## However, you could send a direct one-off email:
    #To: "xxxx@gmail.com"

    #mark:
    #  FirstName: "First"
    #  LastName: "Last"
    #  Identifier: "345345sdfsdf"

    ## We would like the email messages to have access to additional metadata
    template-data:
      # This directive is used to construct `URLCustom` property exposed in the templates
      # In previous example http://evil.com/Identifier/file is {{URLTop}}/{{IdentifierRegex Result}}
      URLTop: "https://xxx.com"

    headers:
      Return-Receipt-To: "info@xxxx.com"
      Disposition-Notification-To: "info@xxxx.com"
      List-Unsubscribe: "<https://www.xxxx.com/unsubscribe?u=876>, <mailto:info@xxxx.com?subject=unsubscribe>"
      List-Unsubscribe-Post: "List-Unsubscribe=One-Click"

    body:
      # Templated HTML / TEXT multipart delivery
      # Templates can substitute dynamic vartiables (See. Template Section for details).
      html: "campaign1/message.htpl"
      text: "campaign1/message.ttpl"

    attach:
      - "/tmp/evil_report.pdf"
    embed:
      - "campaign1/artifacts/logo_header.png"
      - "campaign1/artifacts/logo_footer.png"

##
## Storage module
##
storage:
  DBFile: "campaign1/campaign.db"
  load:
    SourceFile: "campaign1/marks.csv"
    ## Identifiers are used to track marks across the campaign. Identifiers can be of any format, as long as you can carry them in URLs as resources. For example, the URL: http://evil.com/Identifier/file can track access to a hosted payload. You could generate your own Identifiers in the CSV file. Or, you can talk DeepSea to generate unique Identifiers for each mark based on a Regex pattern. You would then need to place "<dynamic>" in place of identifier field in CSV, and use `IdentifierRegex` directive to notify the program the format of the Identifier you want to generate. If for some reason custom generation of Identifiers fail, you will get a 8 Int rand string

    IdentifierRegex: "^[a-z0-9]{8}$"
  query:
    DBTask: "showmarks"
```

Here is the example of a _templated_ email message (`HTPL`). The text multipart (`TTPL`) is similar in content, but different in the format (without HTML markup). 

_Note_: You can generate `TTPL` from `HTPL` by using `dsh2t` tool provided in the distribution. 

```html
<html>
<head></head>
<body>


<div> <img src="cid:{{ index .EmbedImage 0 }}" alt="Header"/> </div>

<h3> Greetings from Frequent Flyers, `{(printf "%s %s" .Mark.Firstname .Mark.Lastname) }` ! </h3>
<p>A new correspondence is waiting for you in the portal 

<p>Kindly review the information so we can assist you in reserving your travel.</p>

<p>
Please visit {(printf "%s/%s" .URLTop .Mark.Identifier) }
</p>

<table width="100%" border="0" cellspacing="0" cellpadding="0">
  <tr>
    <td>
      <div>
        <!--[if mso]>
          <v:roundrect xmlns:v="urn:schemas-microsoft-com:vml" xmlns:w="urn:schemas-microsoft-com:office:word" href="http://litmus.com" style="height:36px;v-text-anchor:middle;width:150px;" arcsize="5%" strokecolor="#FEC426" fillcolor="#FEC426">
            <w:anchorlock/>
            <center style="color:#000000;font-family:Helvetica, Arial,sans-serif;font-size:16px;">Account #E5589-A344</center>
          </v:roundrect>
        <![endif]-->
        <a href="{{.URLTop}}" style="background-color:#FEC426;border:1px solid #FEC426;border-radius:3px;color:#000000;display:inline-block;font-family:sans-serif;font-size:16px;line-height:44px;text-align:center;text-decoration:none;width:200px;-webkit-text-size-adjust:none;mso-hide:all;">Account #E5589-A344</a>
      </div>
    </td>
  </tr>
</table>


<div> <img src="cid:{{ index .EmbedImage 1 }}" alt="Footer"/> </div>

</body>
</html>

```

## Email Message Construction
_This section will most likely be changed as the email mod tools will be rolling into DeepSeea binary_

### Inlining CSS
_Problem:_ There are many editors and mail clients that can be used to create the initial HTML content for the email message. Conversely, you may want to clone content of a web page and modify it. Many times, when straing HTML is used, it is not suitable for proper rendering by some email clients as the CSS is not properly inlined.

_Solution:_ You can use supplied `ds2inline` tool to convert a _normal_ HTML to inlined CSS HTML.

```sh
./dsh2inline  ./css.in ./css.out
```

### Generating Multipart TXT From HTML

_Problem:_ Most email clients support and prefer HTML, some do not. Email gateways check for both HTML and TXT versions of the email document to be present in the delivered envelope. Security stack also takes advantage of that fact when checking for a phishing email. 

_Solution:_ DeepSea automatically constructs multipart email. Expectation is that RTOs provide it with  both version. However, it may be tedious to _textify_ an already existing HTML document by hand, and RTOs may not have access to tools to do so when operating deep in the network. 

DeepSea provides `dsh2t` tool to convert an HTML to a TEXT structured document, to be included into the envelope.

```sh
./dsh2t samples/email.html samples/email.txt
```

## Usage 

```sh
./deepsea 

Usage:
  DeepSea [flags]
  DeepSea [command]

Available Commands:
  config      dump config file
  help        Help about any command
  mailclient  Email a phish
  storage     Manage persistent record storage

Flags:
      --config string   config file (no default)
  -h, --help            help for DeepSea

Use "DeepSea [command] --help" for more information about a command.
```

```sh
$ ./deepsea mailclient --help
Email a phish with features

Usage:
  DeepSea mailclient [flags]

Flags:
  -F, --From string           Message From: header
  -H, --HTMLTemplate string   HTML Template file (.htpl)
  -p, --SMTPPort int          SMTP server port (default 25)
  -s, --SMTPServer string     SMTP server (default "127.0.0.1")
  -U, --SMTPUser string       SMTP user (default "testuser")
  -S, --Subject string        Message Subject: header
  -t, --TLS string            Use TLS handshake (STARTTLS) (default "yes")
  -P, --TextTemplate string   Text Template file (.ttpl)
  -T, --To string             Message To: header
  -h, --help                  help for mailclient

Global Flags:
      --config string   config file (no default)


$ ./deepsea storage --help
STORAGE: TODO

Usage:
  DeepSea storage [flags]
  DeepSea storage [command]

Available Commands:
  load        Load Marks from a file
  query       Query storage

Flags:
  -d, --DBFile string   Path to QL DB file
  -h, --help            help for storage

Global Flags:
      --config string   config file (no default)

Use "DeepSea storage [command] --help" for more information about a command.

```

## Build

Normally, just:

```
go build 
```
However, a static x-platform  build may be desired:

```
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w -extldflags "-static"' .
```

## Email relay support:
Any SMTP compliant server: 
  - Ports 25, 587,486 (SSL), IMAP/S

## Links
[DeepSea](http://github.com/dsnezhkov/deepsea)


## TODO:
- Testing
- Implement direct sendmail (printers): 
  - https://gist.github.com/alok87/56aaecb6c2e102bcf625
  - https://github.com/ansonl/middiff/blob/934ee8e998b27133d9f84d72b3d119f86c671c07/mail.go
- Test against lighter weight C runtime [how to build against musl](https://dominik.honnef.co/posts/2015/06/statically_compiled_go_programs__always__even_with_cgo__using_musl/)