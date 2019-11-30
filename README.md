
## DeepSea Phishing Gear

<img src="https://github.com/dsnezhkov/deepsea/blob/master/docs/images/logo.png" width="180" height="150">

DeepSea phishing gear aims to help RTOs and pentesters with the delivery of opsec-tight, 
flexible email phishing campaigns carried out on the outside as well as on the inside of a perimeter.

> 45 config lines is all you need to consistently send a decent phish ... 

Here's to that

:-------------------------:
<img src="https://github.com/dsnezhkov/deepsea/blob/master/docs/images/config.png" width="350" height="500">


Current Release: *v0.9* 

*Goals*
- Operate with a minimal footprint deep inside enterprises (Internal phish delivery).
- Seamlessly operate with external and internal mail providers (e.g. O365, Gmail, on-prem mail servers)
- Quickly re-target connectivity parameters.
- Flexibly add headers, targets, attachments
- Correctly format and inline email templates, images and multipart messages.
- Use content templates for personalization
- Account for various secure email communication parameters
- Clearly separate artifacts, mark databases and content delivery for multiple (parallel or sequential) phishing campaigns.
- Help create content with minimal dependencies. Embedded tools to support Markdown->HTML->TXT workflow.  
---

### Usage
Read more [here](https://dsnezhkov.github.io/deepsea/) 

### Build

```sh
cd ~/go/src/
export GOPATH=~/go

git clone  https://github.com/dsnezhkov/deepsea
cd deepsea

go get
go build -o deepsea main.go
```
## Operations

### Setup campaigns workspace

```sh
mkdir -p campaigns/campaign1
cp conf/template.yaml  campaigns/campaign/campaign1.yaml
cd campaigns/campaign1
```

### Set Workspace tasks
- edit `campaign.yaml` 

See descriptions of directives in [template](https://github.com/dsnezhkov/deepsea/blob/master/conf/template.yaml)

- edit marks.csv

```csv
ident,email,firstname,lastname
<dynamic>,user@gmail.com,,
```

#### Load Marks
- load marks from CSV defined in the `yml` (creates db. schema automatically)

```
../../deepsea  --config ./campaign.yaml  storage load 
```

Alternatively, split db management tasks:

- create DB
```sh
../../deepsea  --config ./campaign.yaml  storage query -d ./campaign.db -t createtable
Using config file: ./campaign1.yaml
2019/11/18 13:16:16 Task: createtable
2019/11/18 13:16:16 Creating Marks table
```
- load marks from CSV 

```sh
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

vobi97v7, user@gmail.com, , .
```
- you can verify the marks are loaded

```sh
../../deepsea  --config ./campaign.yaml  storage manager -D ./campaign.db -T showmarks
Using config file: ./campaign.yaml
2019/11/18 13:22:17 Task: showmarks
2019/11/18 13:22:17 Querying for result : find()

-= Table: Marks =-
vobi97v7, user@gmail.com, , .
```


#### Create Content

Tow methods: templated and hand-rolled 
##### Templated 
1. Get a decent HTML template
    Ex: ` wget https://raw.githubusercontent.com/leemunroe/responsive-html-email-template/master/email.html`
2. write content
   introduce key/value pairs from `yml`'s `template-data`/`dictonary` and interpolate in the template


3. Inline CSS (if needed) when done with the template (.htpl)

```sh
../../deepsea mailclient --config ./campaign.yaml  content inline

```

4. Create a TXT verson from the HTML version (.ttpl)

```sh
../../deepsea mailclient --config ./campaign.yaml  content multipart
```

##### Hand rolled. Tools
DeepSea provides tools to help roll yourt own html. Most likely you might want to:
- Cretate HTML snippets from Markdown for fast prototyping
- HTML to TEXT for seeing how HTML structure looks in terminal and multipart testing
- Inline CSS Styling for older clients
- Multipart messages

Example (MD2HTML):

```sh
../../deepsea mailclient --config ./campaign.yaml  content md2html  -M ./campaigns/campaign1.md -H ./campaigns/campaign1.html

#STDOUT
../../deepsea mailclient --config ./campaign.yaml  content md2html  -M ./campaigns/campaign1.md 
```

```sh
../../deepsea mailclient --config ./campaign.yaml  content html2text  -K ./campaigns/campaign1.html -L ./campaigns/campaign1.txt
```

#### Mail Campaign

```sh
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
Emailing: user@gmail.com [id:vobi97v7] 
```
Note: We ask for password on the email provider account interactively for now.

### Testing
If you need to run campaign to a test emails, you can reload test marks.
For that, just recycle the data in the marks table like so:

```sh
../../deepsea  --config ./campaign.yaml storage manager -T recycletable
Using config file: ./campaign.yaml
2019/11/18 18:39:17 Task: recycletable
2019/11/18 18:39:17 Dropping table Mark if exists
2019/11/18 18:39:17 Creating Marks table
```

- edit `marks.csv`
- load test marks
```
../../deepsea  --config ./campaign.yaml storage load
```



