package main

import (
	_ "embed"
	// "leaf/wavingleaf/ftp"
	"leaf/wavingleaf/sftp"

	"github.com/getlantern/systray"
)

func main() {
	// systray.Run(onReady, onExit)
	// Init()
	// ftp.FTPTest()
	sftp.TestSFTP()
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
			// 怎么会有这么傻逼的check uncheck逻辑？
			case <-sftpItem.ClickedCh:
				if sftpItem.Checked() {
					sftpItem.Uncheck()
				} else {
					sftpItem.Check()
				}
				print(sftpItem.Checked())
			case <-ftpItem.ClickedCh:
				if ftpItem.Checked() {
					ftpItem.Uncheck()
				} else {
					ftpItem.Check()
				}
				print(ftpItem.Checked())
			case <-webdavItem.ClickedCh:
				if webdavItem.Checked() {
					webdavItem.Uncheck()
				} else {
					webdavItem.Check()
				}
				print(webdavItem.Checked())
			case <-configItem.ClickedCh:
				systray.Quit()
			case <-aboutItem.ClickedCh:
				Info("测试")
			case <-exitItem.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	CloseDB()

}
