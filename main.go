package main

import (
	_ "embed"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, nil)
}

//go:embed icon.png
var icon []byte

func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle("Waving Leaf")
	systray.SetTooltip("一叶飞来细浪生")

	sftpItem := systray.AddMenuItemCheckbox("SFTP", "基于SSH的安全文件传输协议", false)
	ftpItem := systray.AddMenuItemCheckbox("FTP", "传统的文件传输协议", false)
	webdavItem := systray.AddMenuItemCheckbox("WebDav", "较新兴的基于http的文件传输协议", false)
	systray.AddSeparator()
	configItem := systray.AddMenuItem("配置", "配置监听端口与用户信息")
	aboutItem := systray.AddMenuItem("关于", "关于本应用")
	exitItem := systray.AddMenuItem("退出", "退出程序")

	go func() {
		for {
			select {
			case <-sftpItem.ClickedCh:
				print("sftp 启动")
			case <-ftpItem.ClickedCh:
				print("ftp 启动")
			case <-webdavItem.ClickedCh:
				print("webdav 启动")
			case <-configItem.ClickedCh:
				systray.Quit()
			case <-aboutItem.ClickedCh:
				systray.Quit()
			case <-exitItem.ClickedCh:
				systray.Quit()
			}
		}
	}()
}
