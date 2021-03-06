# k8s-extended-apiserver-DIY




## Kube-APIserver

- [x] [Kube Api Server](https://www.youtube.com/watch?v=EJGwWP_qFVw)
- `kc get pods -n kube-system` then can see the pods of kube-system where kube-apiserver, codedns, etcd, controller, scheduler etc belongs. Can see those pods by describing them.
- `kc describe pod kube-apiserver-kind-control-plane -n kube-system`


## TLS & Others

- [x] [How does HTTPS work? What's a CA? What's a self-signed Certificate?](https://www.youtube.com/watch?v=T4Df5_cojAs)
- Prerequisites
    - You need to trust that public key cryptography & signature works
    - Any message encrypted with Bob's public key can only be decrypted with Bob's private key
    - Anyone with access to Alice's public key can verify that a message (signature) could only have been created by someone with access to Alice's private key

- CA (Certificate Authority) && How is a Certificate Signed
    - without CA signed certificate it is basically using HTTP, to use HTTPS need certificate from CA
    - There are a list of CA who is considered as a trusted CA, they can give the certificate, for example: Google CA
    - A CA and a web server etc have a pair of key (everyone can create this type of key pair)
    - When a server who was using HTTP and now want to use HTTPS, it need to generate it's key-pair (public-private key) and give a Certificate Signing Request (created by it's key pair) to a trusted CA for signing it
    - After getting the Certificate Signing Request the CA sign it with it's private key
    - Now, anyone who has the respective CA's public key can verify that it was actually signed by that CA (which is trusted one by that client, ex: my browser)
    - Most Browsers by default have a list of Certificates(CA's certificates) which are issued by a trusted CA, in those certificates it get the public key of that CA
    - It's a good way to prevent "A man in the middle attack"
    - After completing the infos & verifications then the client and server shared a secret key, until that they used asymmetric key encryption(used two key, public-private key pair) but after that when they start passing data by encrypted/decrypted with the same secret key(which they got from each other) they basically start using symmetric key encryption(use only one key)

- Self-Signed Certificate
    - You can create your own CA (create a private-public key pair) and do the same process like previous section said
    - now your different apps can get interact with your another app (with HTTPS) which signed it's certificate by your new CA 
    - it's limited only in your zone/environments

- [x] [SSL, TLS, HTTP, HTTPS](https://www.youtube.com/watch?v=hExRDVZHhig)
    - In HTTP(HyperText Transfer Protocol) the data is clear text, no encryption
    - In HTTPS(Secure HTTP) the data is encrypted
    - SSL(Secure Sockets Layer), uses public key encryption to secure data
        - An SSL certificate is used to authenticate the identity of a website (web server give to to client)
        - Browser make sure it trust the certificate, then the ssl seesion can proceed and encrypted data can be passed
    - TLS (Transport Layer Security)
        - The latest industry standard cryptographic protocol
        - The successor to SSL
        - Authenticates the server, client and encrypts the data

## Self Generating CA certificate Infrastructure

- [x] [How TLS and self-signed certificates work](https://www.youtube.com/watch?v=gH5X7hLAWeU)
- Parts of CA certificate
    - Each CA certificate has three parts
    - private part + public part + crt(public key + claim = CN(common name), SANS(subject alternative names), O(organization) + signature(claim, signed by private key))
    - has extra two part
        - isCa: bool (true | false)
        - Usage:
            - have many things in usage, but 4 things are mainly useful
            - digital signature
            - key encipherment
            - server auth
            - client auth
            - cert sign (not important for us)
- To make new CA | To issue new CA | To make new CA infrastructure | All are self signed not by known CA
    - Step1: Generate CA cert pair(public and private key) (use `openssl` for this and see in certmanager site for the command)
    - After Step1 will generate two files
        - ca.key (base64 encoded private-key)
        - ca.crt (all the parts that told in Parts of CA certificate will be here. base64 encoded)  [it's provided to every browser/pc by default]
            - isCA needs to be `true` of the cert [isCA: true]
            - Issuer will be same as it's CN(Common Name)   [Issuer: ${CN}]
            - SANS generally `empty`  [SANS: empty]
            - Usage:
                - digital signature
                - key encipherment
                - cert sign
    - Now, with this can generate a lot of server certificate or client certificate
    - For generating server certificate
        - again we need to generate a private key using openssl
            - server.key
        - now during making server.crt we will give the CN and others  things as usual but we need to `sign` the server certificate with our generated CA's private key
            - server.crt er sign part er private key ta hobe `ca.key`(ca's private key)
        - Example of server certificate
            - (pub + claim = pg.default.svc, SANS=[domain,ip] + sign(claim, ca.key))
    - Now server will provide it's certificate to client when through curl/browser client opens the HTTPS connection 
    - Now client will valided that whether the signing part(of server.crt) is actually signed by trusted CA or not by using known CA's public key(which generally belongs to client/browswer/crul etc)
    - After sign validation and CN, O, SANS, expire date etc then client make sure yes it is the actual site that client was looking for
    - Note that, claim part is not encrypted, claim just need to be valid but the sign part is encrypted so sign need to be decrypted by CA's public key and need to be match.
    - Also note that, the `ca.key` need to be private, cause otherwise anyone can issue certificate with this by the main CA's name            
    - `dig +short <site>` : to get the ip of the site
    - Note that for `server.crt`
        - as it is not a CA so it's `isCa: false`
        - it's Usage:
            - digital signature
            - key encipherment
            - server auth (it's need to be on, it's particularly important, if server auth is not on then you cannot use it as a server certificate)
                - when write a server in go and give listen by tls config then if this part is not given or on then that will fail (TLS protocol will fail)
    - Client certificate also need to validated(if it is on) by server, in this case client_ca.crt also need to give to server, it can be same or different than server ca.crt(in kubernetes we give client ca)
    - Client certificate also generate like the server certificate
        - client.key (private key)
        - client.crt / ca.key
        - it's `isCa: false`
        - Usage:
            - digital signature
            - key encipherment
            - client auth
        - (pub + claim = CN:pg-admin(means it will be user name) + sign(claim, client_ca.key(with this the signature will be signed)))  
    - For kubernetes webhooks we need to make a ca, that ca is basically a client ca, for that we need to give ca certificate bundle when we make apiservice object, by this we tell k8s server when we will connect with you, you verify my client certificates with this public key.
    - when we maintained both server and client then we can use self-signed certificate

- mTLS / mutual TLS
    - client and server both will identified each other, this is basically mutual TLS
    - during making client object in go, need to provide client.key and client.crt pair in tls config    
    - then we connect in server side, in client need server ca.crt and it's certs pair(of client.key and client.crt)
        - ```go
          http.Client{
            RootCA: ca.crt
            Client certs: [client.key, client.crt]
          }
          ``` 
- CRL (Certificate Revokation List)
    - CA authority has a list, where the id number of revoked certificate
    - when a client connect then server will check whether the certificate is valid or not, then it will check the revoke list, if it is there then server will reject that client connection
    - Standard TLS client library of Go does not support this CRL method
    
- `certlogik.com/decoder`: For decoding the key, note that never give your private key here
- `maxmind.com` for seeing the IP location & Info
- `traceroute google.com` to see the route from your ip to google.com's ip


## Webhooks in K8s

There are three webhooks in kubernetes.  
   - Mutation Webhook  
   - Validation Webhook  
   - Conversion Webhook  


## Encryption

### Symmetric key Encryption


### Asymmetric key encryption





# Resources

- [x] [Basic concepts of web applications, how they work and the HTTP protocol](https://www.youtube.com/watch?v=RsQ1tFLwldY)
