/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-03 20:24:56
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-14 11:30:02
 * @FilePath: /gmvpn/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gmvpn/app"
	"gmvpn/common/config"
)

var _version = "v0.1.0"

func main() {
	config := config.Config{}
	flag.StringVar(&config.Device, "dev", "", "device name")
	flag.StringVar(&config.CIDR, "cidr", "10.8.8.0/24", "tun interface cidr")
	flag.StringVar(&config.CIDRv6, "cidr6", "fced:9999::9999/64", "tun interface ipv6 cidr")
	flag.IntVar(&config.MTU, "mtu", 1500, "tun mtu")
	flag.StringVar(&config.LocalAddr, "local", ":3001", "bind to local address")
	flag.StringVar(&config.RemoteAddr, "remote", "", "remote server address")

	flag.Func("route", "push ipv4 route to client, value example: 192.168.11.0/24,192.168.13.0/24", func(s string) error {
		config.Route = strings.Split(s, ",")
		return nil
	})
	flag.Func("route6", "push ipv6 route to client, value example: 2001:250:4000:2000::1/64", func(s string) error {
		config.Route6 = strings.Split(s, ",")
		return nil
	})
	flag.BoolVar(&config.ServerMode, "s", false, "server mode")
	flag.BoolVar(&config.AutoSnat, "autonat", true, "open iptables auto snat on server mode")
	flag.BoolVar(&config.GlobalMode, "g", false, "client global mode")
	flag.BoolVar(&config.Compress, "compress", false, "enable data compression")
	flag.IntVar(&config.Timeout, "t", 30, "dial timeout in seconds")
	flag.StringVar(&config.CaPath, "ca", "./certs/ca.crt", "gmtls ca file path")
	flag.StringVar(&config.TLSCertificateFilePath, "certificate", "./certs/server.pem", "tls certificate file path")
	flag.StringVar(&config.TLSCertificateKeyFilePath, "privatekey", "./certs/server.key", "tls certificate key file path")
	flag.StringVar(&config.SignCertPath, "signcert", "./certs/signcert.crt", "gmtls sign cert file path")
	flag.StringVar(&config.SignKeyPath, "signkey", "./certs/signkey.key", "gmtls sign key file path")
	flag.StringVar(&config.EncCertPath, "enccert", "./certs/enccert.crt", "gmtls enc cert file path")
	flag.StringVar(&config.EncKeyPath, "enckey", "./certs/enckey.key", "gmtls enc key file path")
	flag.StringVar(&config.TLSSni, "sni", "", "tls handshake sni")
	flag.BoolVar(&config.TLSInsecureSkipVerify, "isv", false, "tls insecure skip verify")
	flag.StringVar(&config.TLSCipher, "cipher", "SM2_WITH_SM4_SM3", "tls cipher suites")
	flag.Parse()

	if config.RemoteAddr == "" {
		config.ServerMode = true
	}

	app := app.NewApp(&config, _version)
	app.InitConfig()
	go app.StartApp()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.StopApp()
}
