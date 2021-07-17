package confs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

// func CertInit() {
// 	_, keyErr := os.Stat(`./key.pem`)
// 	_, certErr := os.Stat(`./cert.pem`)
// 	if (keyErr == nil || os.IsExist(keyErr)) && (certErr == nil || os.IsExist(certErr)) {
// 		return
// 	} else {
// 		fmt.Println(`[警告] 密钥、证书不存在，正在创建...`)

// 		cmd := exec.Command(`openssl`,
// 			`genrsa`, `-out`, `key.pem`, `2048`)
// 		if out, err := cmd.CombinedOutput(); err != nil {
// 			logs.ErrorPanic(err, `密钥创建错误 -> `+string(out))
// 		}

// 		cmd = exec.Command(`openssl`,
// 			`req`, `-new`, `-x509`, `-key`, `key.pem`,
// 			`-out`, `cert.pem`, `-days`, `3650`,
// 			`-subj`, `/C=CN/ST=BJ/L=BJ/O=None/OU=None/CN=None/emailAddress=None`)
// 		if out, err := cmd.CombinedOutput(); err != nil {
// 			logs.ErrorPanic(err, `证书创建错误 -> `+string(out))
// 		}

// 		fmt.Println(`[警告] 密钥、证书创建成功`)
// 		return
// 	}
// }

func CertInit() {
	_, keyErr := os.Stat(`./key.pem`)
	_, certErr := os.Stat(`./cert.pem`)
	if (keyErr == nil || os.IsExist(keyErr)) && (certErr == nil || os.IsExist(certErr)) {
		return
	} else {
		max := new(big.Int).Lsh(big.NewInt(1), 128)
		serialNumber, _ := rand.Int(rand.Reader, max)

		template := x509.Certificate{
			SerialNumber: serialNumber,
			NotBefore:    time.Now(),
			NotAfter:     time.Now().Add(365 * 24 * time.Hour),
			KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			IPAddresses:  []net.IP{net.ParseIP("39.107.92.179")},
		}

		pk, _ := rsa.GenerateKey(rand.Reader, 2048)

		derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
		certOut, _ := os.Create("cert.pem")
		pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
		certOut.Close()

		keyOut, _ := os.Create("key.pem")
		pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
		keyOut.Close()
	}
}
