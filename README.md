# smtp2api
*smtp2api* is an SMTP server written in Golang that receives email data and sends it through a *BREVO* email server API.

## What is it for
*smtp2api* is useful when an application installed on your server cannot access the SMTP port to send email. In this case, *smtp2api* works as a proxy, receiving requests via the SMTP protocol and forwarding them via **HTTP** or **HTTPS API** to some email service that accepts this type of request.

For example, if you have an application installed on your server that uses the SMTP protocol to send email, but the firewall block any SMTP port. So, you could use *smtp2api* to forward SMTP requests to an email service **API** that has this service.


## Installation
To install smtp2api, you can use the following command:

```
go get github.com/dmslabsbr/smtp2api
```

You need to set BREVO_APIKEY as enviromment variable.

```
export BREVO_APIKEY=<your_brevo_api_key>
```
### Docker

 Build the image
```
docker build -t smtp2api .
```

 Run the image

 

## Documentation
The documentation for smtp2api is available at the following link:

https://github.com/bravocode/smtp2api/blob/master/README.md


## Contributions
Contributions are welcome. To contribute to *smtp2api*, you can do the following:

Create a fork of the repository.
Make your changes to the fork.
Submit a pull request.
License
smtp2api is licensed under the MIT license.