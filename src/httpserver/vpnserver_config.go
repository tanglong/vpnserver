package main

import ()

type xmlLog struct {
	Config string `xml:"config"`
}

type xmlRedis struct {
	Address string `xml:"address"`
	Index   int    `xml:"index"`
}

type xmlConfig struct {
	Log   xmlLog   `xml:"log"`
	Http  string   `xml:"http"`
	Redis xmlRedis `xml:"redis"`
	OutIp string   `xml:"outip"`
}
