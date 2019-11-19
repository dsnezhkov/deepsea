### Build tool

cd ~/go/src/
export GOPATH=~/go

git clone  https://github.com/dsnezhkov/deepsea
cd deepsea

go get 
go build -o deepsea main.go

### Setup campaigns workspace
mkdir -p campaigns/campaign1
cp conf/.deepsea.yaml  campaigns/campaign/campaign1.yaml
cd campaigns/campaign1

### Workspace tasks
- edit campaign.yaml

mailclient:

  connection:

    SMTPUser: "user@outlook.com"
    SMTPServer: "smtp.office365.com"
    SMTPPort: 587
    TLS: "yes"

  message:

    Subject: "Subject."
    ## Some providers, namely MSFT does not like to relay arbitrary emails.
    ## Make sure the "From" is your@outlook.com
    ## Or you get: `554 5.2.0 STOREDRV.Submission.Exception:SendAsDeniedException.MapiExceptionSendAsDenied`

    From: "user@outlook.com"
    # File means we take marks from the database
    To: "campaign.db"

    ## We would like the email messages to have access to additional metadata
    template-data:
      # This directive is used to construct `URLCustom` property exposed in the templates
      # In previous example http://evil.com/Identifier/file is {{URLTop}}/{{IdentifierRegex Result}}
      URLTop: "https://example.com"

    headers:
      Return-Receipt-To: "user@outlook.com"
      Disposition-Notification-To: "user@outlook.com"
      List-Unsubscribe: "<https://www.microsoft.com/unsubscribe?u=876>, <mailto:user@outlook.com?subject=unsubscribe>"
      List-Unsubscribe-Post: "List-Unsubscribe=One-Click"

    body:
      # Templated HTML / TEXT multipart delivery
      # Templates can substitute dynamic vartiables (See. Template Section for details).
      html: "message.htpl"
      text: "message.ttpl"

    # No attachemnts
    attach:
    # No img embeds
    embed:

##
## Storage module
##
storage:
  # location of the database of imported marks
  DBFile: "campaign.db"
  load:
    # locations of the CSV file of marks to import into database ^^
    SourceFile: "marks.csv"

    # Generate identified based on pattern
    IdentifierRegex: "^[a-z0-9]{8}$"
  query:
    DBTask: "showmarks"
```

- edit marks.csv

```
ident,email,firstname,lastname
<dynamic>,dsnezhkov@gmail.com,,
```

#### Load Marks
- create database 
```
../../deepsea  --config ./campaign.yaml  storage query -d ./campaign.db -t createtable
Using config file: ./campaign1.yaml
2019/11/18 13:16:16 Task: createtable
2019/11/18 13:16:16 Creating Marks table
```
- load marks from CSV 
```
../../deepsea  --config ./campaign.yaml  storage load -d ./campaign.db -s ./marks.csv
Using config file: ./campaign.yaml
2019/11/18 13:21:11 Dropping table Mark if exists
2019/11/18 13:21:11 Creating Marks table
2019/11/18 13:21:11 Pointing to mark table
2019/11/18 13:21:11 Removing existing rows if any
2019/11/18 13:21:11 Inserting a row
2019/11/18 13:21:11 Querying for result : find()
2019/11/18 13:21:11 Getting all results
2019/11/18 13:21:11 Printing Marks

vobi97v7, dsnezhkov@gmail.com, , .
```
- you can verify the marks are loaded

```
../../deepsea  --config ./campaign.yaml  storage query -d ./campaign.db -t showmarks
Using config file: ./campaign.yaml
2019/11/18 13:22:17 Task: showmarks
2019/11/18 13:22:17 Querying for result : find()

-= Table: Marks =-
vobi97v7, dsnezhkov@gmail.com, , .
```


#### Create Content

- Get a decent HTML template
    Ex: ` wget https://raw.githubusercontent.com/leemunroe/responsive-html-email-template/master/email.html`
- write content
 
- 1. Inline CSS (if needed)
```
../../tools/dsh2inline message.html message.htpl
```
- 2. Create a TXT verson from the HTML version
```
../../tools/dsh2t message.htpl message.ttpl
```

#### Mail Campaign
Note: We ask for interactive password on the email provider account for now.
```
../../deepsea mailclient --config ./campaign.yaml 

Using config file: ./campaign.yaml
SMTP Server : smtp.office365.com
SMTP Port   : 587
SMTP User : user@outlook.com
SMTP TLS : yes
From: user@outlook.com
To: campaign.db
Subject: Subject.
Text Template: message.ttpl
HTML Template: message.htpl

-= SMTP Authentication Credentials for smtp.office365.com =-
Enter Password: 

2019/11/18 18:14:18 Pointing to mark table
2019/11/18 18:14:18 Querying for result : find()
2019/11/18 18:14:18 Getting all results
2019/11/18 18:14:18 -= Marks =-
Emailing: dsnezhkov@gmail.com [id:vobi97v7] 
```

### Testing
If you need to run campaign to a test emails, you can reload test marks.
For that, just recycle the data in the marks table like so:

``` 
../../deepsea  --config ./campaign.yaml storage query -t recycletable
Using config file: ./campaign.yaml
2019/11/18 18:39:17 Task: recycletable
2019/11/18 18:39:17 Dropping table Mark if exists
2019/11/18 18:39:17 Creating Marks table
```

```
../../deepsea  --config ./campaign.yaml storage query showmarks
Using config file: ./campaign.yaml
2019/11/18 18:39:24 Task: showmarks
2019/11/18 18:39:24 Querying for result : find()
-= Table: Marks =-
``` 
