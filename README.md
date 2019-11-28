## DeepSea Phishing Gear

- v0.8
DeepSea phishing gear aims to help RTOs and pentesters with the delivery of opsec-tight, flexible email phishing campaigns carried out on the outside as well as on the inside of a perimeter.

#### Goals:
- Operate with a minimal footprint deep inside enterprises (Internal phish delivery).
- Seamlessly operate with external and internal mail providers (e.g. O365, Gmail, on-prem mail servers)
- Quickly re-target connectivity parameters.
- Flexibly add headers, targets, attachments
- Correctly format and inline email templates, images and multipart messages.
- Use content templates for personalization
- Account for various secure email communication parameters
- Clearly separate artifacts, mark databases and content delivery for multiple (parallel or sequential) phishing campaigns.

### Usage
Read more: 
- [Here](https://dsnezhkov.github.io/deepsea/) 
- [Here](https://github.com/dsnezhkov/deepsea/blob/master/docs/campaign.md)

### Build

```sh
cd ~/go/src/
export GOPATH=~/go

git clone  https://github.com/dsnezhkov/deepsea
cd deepsea

go get
go build -o deepsea main.go
```



