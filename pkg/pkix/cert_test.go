package pkix_test

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"testing"

	"github.com/openebl/openebl/pkg/pkix"
)

func TestVerifyWithCustomizedRootCertificates(t *testing.T) {
	rootCert, err := LoadCert("../../credential/root_ca.crt")
	if err != nil {
		t.Fatal(err)
	}

	fingerPrint := sha1.Sum(rootCert.Raw)
	fmt.Println(hex.EncodeToString(fingerPrint[:]))

	cert, err := LoadCert("../../credential/bob_ecc.crt")
	if err != nil {
		t.Fatal(err)
	}
	err = pkix.Verify([]*x509.Certificate{cert}, []*x509.Certificate{rootCert})
	if err != nil {
		t.Fatal(err)
	}
}

func TestVerifyWithIntermediatesCertificates(t *testing.T) {
	rootCert, err := LoadCert("../../credential/root_ca.crt")
	if err != nil {
		t.Fatal(err)
	}
	intermediateCert, err := LoadCert("../../credential/bob_ecc.crt")
	if err != nil {
		t.Fatal(err)
	}
	cert, err := LoadCert("../../credential/bob_ecc2.crt")
	if err != nil {
		t.Fatal(err)
	}
	err = pkix.Verify([]*x509.Certificate{cert, intermediateCert}, []*x509.Certificate{rootCert})
	if err != nil {
		t.Fatal(err)
	}
}

func TestVerifyWithWrongIntermediatesCertificates(t *testing.T) {
	rootCert, err := LoadCert("../../credential/root_ca.crt")
	if err != nil {
		t.Fatal(err)
	}
	intermediateCert, err := LoadCert("../../credential/bob_ecc.crt")
	if err != nil {
		t.Fatal(err)
	}
	cert, err := LoadCert("../../credential/alice_ecc.crt")
	if err != nil {
		t.Fatal(err)
	}
	err = pkix.Verify([]*x509.Certificate{cert, intermediateCert}, []*x509.Certificate{rootCert})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func LoadCert(fileName string) (*x509.Certificate, error) {
	pemFile, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	certBlock, _ := pem.Decode(pemFile)
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func TestParseCertificate(t *testing.T) {
	pemData := `-----BEGIN CERTIFICATE-----
MIIFVDCCBDygAwIBAgIRAMj6vmF8SNMvEAU489YLMGUwDQYJKoZIhvcNAQELBQAw
RjELMAkGA1UEBhMCVVMxIjAgBgNVBAoTGUdvb2dsZSBUcnVzdCBTZXJ2aWNlcyBM
TEMxEzARBgNVBAMTCkdUUyBDQSAxQzMwHhcNMjQwMTA5MDYzMTM5WhcNMjQwNDAy
MDYzMTM4WjAZMRcwFQYDVQQDEw53d3cuZ29vZ2xlLmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBALHfbzZh9gDuq6YU18lZGS1xxvFJ9GWpX+EdqQ2T
iAw6hTS8vFNG/jt76uZhJlRK33derWlvpq+Bbct3pqkYp4kMkFFMURDRvPFrX/3t
Tp2Mv9V9Br1GvB9VXLYFDGpmpPi6LlMDMJMkUOczb4QuDxJ21wdyL62DbVJxGuqv
kAk0cRAPhtMC7ZYGBSqaXOwhHneuzzE5UBlRqODALuUdmBAmbgXd+UxvUsavmqt1
7AYtFiVj8lgEsrXGRFEfGYaaIOXKKzNNQwC4D3B/yEPO1qFVT7ZaIbGWkGf8F0AK
xYClusXyRFIeJNt1atjsfyFwf/gDomwvc+B1LuFJQ1J+ZIECAwEAAaOCAmgwggJk
MA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8E
AjAAMB0GA1UdDgQWBBQRLzybIKbyXeW1ms1CC+hmzv1B2DAfBgNVHSMEGDAWgBSK
dH+vhc3ulc09nNDiRhTzcTUdJzBqBggrBgEFBQcBAQReMFwwJwYIKwYBBQUHMAGG
G2h0dHA6Ly9vY3NwLnBraS5nb29nL2d0czFjMzAxBggrBgEFBQcwAoYlaHR0cDov
L3BraS5nb29nL3JlcG8vY2VydHMvZ3RzMWMzLmRlcjAZBgNVHREEEjAQgg53d3cu
Z29vZ2xlLmNvbTAhBgNVHSAEGjAYMAgGBmeBDAECATAMBgorBgEEAdZ5AgUDMDwG
A1UdHwQ1MDMwMaAvoC2GK2h0dHA6Ly9jcmxzLnBraS5nb29nL2d0czFjMy9RcUZ4
Ymk5TTQ4Yy5jcmwwggEFBgorBgEEAdZ5AgQCBIH2BIHzAPEAdgB2/4g/Crb7lVHC
Ycz1h7o0tKTNuyncaEIKn+ZnTFo6dAAAAYztIlshAAAEAwBHMEUCIQDJzUYIHoy3
xGVtTGPoj5JSC14ZrhUJhEK8PFiDFh7emgIgIbfFX+so1ifIzaaaaDa6u+rzYL/o
qLK6PzVOx0vTrBcAdwBIsONr2qZHNA/lagL6nTDrHFIBy1bdLIHZu7+rOdiEcwAA
AYztIlsWAAAEAwBIMEYCIQCfzSZyKnYqoLc8xw7Djrbezmj9wKeDLiL6UN29hXkn
jQIhAKelO6nIhdU+tGGdY46PXEHpb9REGDlF+mvgb8MZ4waQMA0GCSqGSIb3DQEB
CwUAA4IBAQCpyQ5acH1BrM4jGnXclQKBFh8WfOS/lfzDi4HruRkr2w24qPvbJOy6
3ebm5tKv33cN93GfWFv9Ioy/47O9TyZCzEJYSRH4WDAsj9m9gGiknkvLJfsOaqDo
GX2tmAUu1iUlZv8LDEdfz7lNFDEmWGUM6570bySPX4Ea1w/FOKS0KTNto/OkRkkN
P9Mnj2KGbV6jW3M8TZe5pfGOQk8rllIcnMs72oiKDeaQzPWy2b7Ckil6Ye1ZtKP2
id+tcxoUlhRW/2wTCIVcnvRTBI/gJoVXEg3uKJouUb4NSeA1WKINRy7+1MOi/Jaa
kPTNzO8cc/TgWULoU66wj7DJGI3iIWt/
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFljCCA36gAwIBAgINAgO8U1lrNMcY9QFQZjANBgkqhkiG9w0BAQsFADBHMQsw
CQYDVQQGEwJVUzEiMCAGA1UEChMZR29vZ2xlIFRydXN0IFNlcnZpY2VzIExMQzEU
MBIGA1UEAxMLR1RTIFJvb3QgUjEwHhcNMjAwODEzMDAwMDQyWhcNMjcwOTMwMDAw
MDQyWjBGMQswCQYDVQQGEwJVUzEiMCAGA1UEChMZR29vZ2xlIFRydXN0IFNlcnZp
Y2VzIExMQzETMBEGA1UEAxMKR1RTIENBIDFDMzCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAPWI3+dijB43+DdCkH9sh9D7ZYIl/ejLa6T/belaI+KZ9hzp
kgOZE3wJCor6QtZeViSqejOEH9Hpabu5dOxXTGZok3c3VVP+ORBNtzS7XyV3NzsX
lOo85Z3VvMO0Q+sup0fvsEQRY9i0QYXdQTBIkxu/t/bgRQIh4JZCF8/ZK2VWNAcm
BA2o/X3KLu/qSHw3TT8An4Pf73WELnlXXPxXbhqW//yMmqaZviXZf5YsBvcRKgKA
gOtjGDxQSYflispfGStZloEAoPtR28p3CwvJlk/vcEnHXG0g/Zm0tOLKLnf9LdwL
tmsTDIwZKxeWmLnwi/agJ7u2441Rj72ux5uxiZ0CAwEAAaOCAYAwggF8MA4GA1Ud
DwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0T
AQH/BAgwBgEB/wIBADAdBgNVHQ4EFgQUinR/r4XN7pXNPZzQ4kYU83E1HScwHwYD
VR0jBBgwFoAU5K8rJnEaK0gnhS9SZizv8IkTcT4waAYIKwYBBQUHAQEEXDBaMCYG
CCsGAQUFBzABhhpodHRwOi8vb2NzcC5wa2kuZ29vZy9ndHNyMTAwBggrBgEFBQcw
AoYkaHR0cDovL3BraS5nb29nL3JlcG8vY2VydHMvZ3RzcjEuZGVyMDQGA1UdHwQt
MCswKaAnoCWGI2h0dHA6Ly9jcmwucGtpLmdvb2cvZ3RzcjEvZ3RzcjEuY3JsMFcG
A1UdIARQME4wOAYKKwYBBAHWeQIFAzAqMCgGCCsGAQUFBwIBFhxodHRwczovL3Br
aS5nb29nL3JlcG9zaXRvcnkvMAgGBmeBDAECATAIBgZngQwBAgIwDQYJKoZIhvcN
AQELBQADggIBAIl9rCBcDDy+mqhXlRu0rvqrpXJxtDaV/d9AEQNMwkYUuxQkq/BQ
cSLbrcRuf8/xam/IgxvYzolfh2yHuKkMo5uhYpSTld9brmYZCwKWnvy15xBpPnrL
RklfRuFBsdeYTWU0AIAaP0+fbH9JAIFTQaSSIYKCGvGjRFsqUBITTcFTNvNCCK9U
+o53UxtkOCcXCb1YyRt8OS1b887U7ZfbFAO/CVMkH8IMBHmYJvJh8VNS/UKMG2Yr
PxWhu//2m+OBmgEGcYk1KCTd4b3rGS3hSMs9WYNRtHTGnXzGsYZbr8w0xNPM1IER
lQCh9BIiAfq0g3GvjLeMcySsN1PCAJA/Ef5c7TaUEDu9Ka7ixzpiO2xj2YC/WXGs
Yye5TBeg2vZzFb8q3o/zpWwygTMD0IZRcZk0upONXbVRWPeyk+gB9lm+cZv9TSjO
z23HFtz30dZGm6fKa+l3D/2gthsjgx0QGtkJAITgRNOidSOzNIb2ILCkXhAd4FJG
AJ2xDx8hcFH1mt0G/FX0Kw4zd8NLQsLxdxP8c4CU6x+7Nz/OAipmsHMdMqUybDKw
juDEI/9bfU1lcKwrmz3O2+BtjjKAvpafkmO8l7tdufThcV4q5O8DIrGKZTqPwJNl
1IXNDw9bg1kWRxYtnCQ6yICmJhSFm/Y3m6xv+cXDBlHz4n/FsRC6UfTd
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFYjCCBEqgAwIBAgIQd70NbNs2+RrqIQ/E8FjTDTANBgkqhkiG9w0BAQsFADBX
MQswCQYDVQQGEwJCRTEZMBcGA1UEChMQR2xvYmFsU2lnbiBudi1zYTEQMA4GA1UE
CxMHUm9vdCBDQTEbMBkGA1UEAxMSR2xvYmFsU2lnbiBSb290IENBMB4XDTIwMDYx
OTAwMDA0MloXDTI4MDEyODAwMDA0MlowRzELMAkGA1UEBhMCVVMxIjAgBgNVBAoT
GUdvb2dsZSBUcnVzdCBTZXJ2aWNlcyBMTEMxFDASBgNVBAMTC0dUUyBSb290IFIx
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAthECix7joXebO9y/lD63
ladAPKH9gvl9MgaCcfb2jH/76Nu8ai6Xl6OMS/kr9rH5zoQdsfnFl97vufKj6bwS
iV6nqlKr+CMny6SxnGPb15l+8Ape62im9MZaRw1NEDPjTrETo8gYbEvs/AmQ351k
KSUjB6G00j0uYODP0gmHu81I8E3CwnqIiru6z1kZ1q+PsAewnjHxgsHA3y6mbWwZ
DrXYfiYaRQM9sHmklCitD38m5agI/pboPGiUU+6DOogrFZYJsuB6jC511pzrp1Zk
j5ZPaK49l8KEj8C8QMALXL32h7M1bKwYUH+E4EzNktMg6TO8UpmvMrUpsyUqtEj5
cuHKZPfmghCN6J3Cioj6OGaK/GP5Afl4/Xtcd/p2h/rs37EOeZVXtL0m79YB0esW
CruOC7XFxYpVq9Os6pFLKcwZpDIlTirxZUTQAs6qzkm06p98g7BAe+dDq6dso499
iYH6TKX/1Y7DzkvgtdizjkXPdsDtQCv9Uw+wp9U7DbGKogPeMa3Md+pvez7W35Ei
Eua++tgy/BBjFFFy3l3WFpO9KWgz7zpm7AeKJt8T11dleCfeXkkUAKIAf5qoIbap
sZWwpbkNFhHax2xIPEDgfg1azVY80ZcFuctL7TlLnMQ/0lUTbiSw1nH69MG6zO0b
9f6BQdgAmD06yK56mDcYBZUCAwEAAaOCATgwggE0MA4GA1UdDwEB/wQEAwIBhjAP
BgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTkrysmcRorSCeFL1JmLO/wiRNxPjAf
BgNVHSMEGDAWgBRge2YaRQ2XyolQL30EzTSo//z9SzBgBggrBgEFBQcBAQRUMFIw
JQYIKwYBBQUHMAGGGWh0dHA6Ly9vY3NwLnBraS5nb29nL2dzcjEwKQYIKwYBBQUH
MAKGHWh0dHA6Ly9wa2kuZ29vZy9nc3IxL2dzcjEuY3J0MDIGA1UdHwQrMCkwJ6Al
oCOGIWh0dHA6Ly9jcmwucGtpLmdvb2cvZ3NyMS9nc3IxLmNybDA7BgNVHSAENDAy
MAgGBmeBDAECATAIBgZngQwBAgIwDQYLKwYBBAHWeQIFAwIwDQYLKwYBBAHWeQIF
AwMwDQYJKoZIhvcNAQELBQADggEBADSkHrEoo9C0dhemMXoh6dFSPsjbdBZBiLg9
NR3t5P+T4Vxfq7vqfM/b5A3Ri1fyJm9bvhdGaJQ3b2t6yMAYN/olUazsaL+yyEn9
WprKASOshIArAoyZl+tJaox118fessmXn1hIVw41oeQa1v1vg4Fv74zPl6/AhSrw
9U5pCZEt4Wi4wStz6dTZ/CLANx8LZh1J7QJVj2fhMtfTJr9w4z30Z209fOU0iOMy
+qduBmpvvYuR7hZL6Dupszfnw0Skfths18dG9ZKb59UhvmaSGZRVbNQpsg3BZlvi
d0lIKO2d1xozclOzgjXPYovJJIultzkMu34qQb9Sz/yilrbCgj8=
-----END CERTIFICATE-----
`

	certs, err := pkix.ParseCertificate([]byte(pemData))
	if err != nil {
		t.Fatal(err)
	}
	if len(certs) != 3 {
		t.Fatalf("expected 3 certificates, got %d", len(certs))
	}
}
