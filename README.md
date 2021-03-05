# k8s-extended-apiserver-DIY




## Kube-APIserver

- [x] [Kube Api Server](https://www.youtube.com/watch?v=EJGwWP_qFVw)
- `kc get pods -n kube-system` then can see the pods of kube-system where kube-apiserver, codedns, etcd, controller, scheduler etc belongs. Can see those pods by describing them.
- `kc describe pod kube-apiserver-kind-control-plane -n kube-system`


## HTTPS, CA, Self-Signed Certificate

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
        - Tu successor to SSL
        - Authenticates the server, client and encrypts the data
        