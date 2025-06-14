<div align="center">
  <img width="100px" alt="logo" src="https://raw.githubusercontent.com/OpenListTeam/Logo/main/OpenList.svg"/></a>
  <p><em>🗂一个支持多存储的文件列表程序，使用 Gin 和 Solidjs，基于 alist 项目 fork 开发</em></p>
<div>
  <a href="https://goreportcard.com/report/github.com/OpenListTeam/OpenList/v3">
    <img src="https://goreportcard.com/badge/github.com/OpenListTeam/OpenList/v3" alt="latest version" />
  </a>
  <a href="https://github.com/OpenListTeam/OpenList/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/OpenListTeam/OpenList" alt="License" />
  </a>
  <a href="https://github.com/OpenListTeam/OpenList/actions?query=workflow%3ABuild">
    <img src="https://img.shields.io/github/actions/workflow/status/OpenListTeam/OpenList/build.yml?branch=main" alt="Build status" />
  </a>
  <a href="https://github.com/OpenListTeam/OpenList/releases">
    <img src="https://img.shields.io/github/release/OpenListTeam/OpenList" alt="latest version" />
  </a>
</div>
<div>
  <a href="https://github.com/OpenListTeam/OpenList/discussions">
    <img src="https://img.shields.io/github/discussions/OpenListTeam/OpenList?color=%23ED8936" alt="discussions" />
  </a>
  <a href="https://github.com/OpenListTeam/OpenList/releases">
    <img src="https://img.shields.io/github/downloads/OpenListTeam/OpenList/total?color=%239F7AEA&logo=github" alt="Downloads" />
  </a>
</div>
</div>

---

> [!IMPORTANT]
>
> **这是FORK版本！！！**
>
> 我们诚挚地感谢原作者 [xhofe/alist](https://github.com/xhofe/alist) 为原项目做出的重大贡献。
>
> 本 Fork 尚未稳定, 大量外部链接指向的内容尚未得到审计, 存在投毒风险, 含相关文档网站等. 切勿盲目信任!
>
> 当前仅移除了被确认超出原作者控制范围的外部链接（如`alistgo.com`），项目具体迁移进度可在[OpenList 迁移工作总结](https://github.com/OpenListTeam/OpenList/issues/6)中查看

[English](./README.md) | 中文 | [日本語](./README_ja.md) | [Contributing](./CONTRIBUTING.md) | [CODE_OF_CONDUCT](./CODE_OF_CONDUCT.md)

## 功能

- [x] 多种存储
    - [x] 本地存储
    - [x] [阿里云盘](https://www.alipan.com/)
    - [x] OneDrive / Sharepoint（[国际版](https://www.office.com/), [世纪互联](https://portal.partner.microsoftonline.cn),de,us）
    - [x] [天翼云盘](https://cloud.189.cn) (个人云, 家庭云)
    - [x] [GoogleDrive](https://drive.google.com/)
    - [x] [123云盘](https://www.123pan.com/)
    - [x] FTP / SFTP
    - [x] [PikPak](https://www.mypikpak.com/)
    - [x] [S3](https://aws.amazon.com/cn/s3/)
    - [x] [Seafile](https://seafile.com/)
    - [x] [又拍云对象存储](https://www.upyun.com/products/file-storage)
    - [x] WebDav(支持无API的OneDrive/SharePoint)
    - [x] Teambition（[中国](https://www.teambition.com/ )，[国际](https://us.teambition.com/ )）
    - [x] [分秒帧](https://www.mediatrack.cn/)
    - [x] [和彩云](https://yun.139.com/) (个人云, 家庭云，共享群组)
    - [x] [Yandex.Disk](https://disk.yandex.com/)
    - [x] [百度网盘](http://pan.baidu.com/)
    - [x] [UC网盘](https://drive.uc.cn)
    - [x] [夸克网盘](https://pan.quark.cn)
    - [x] [迅雷网盘](https://pan.xunlei.com)
    - [x] [蓝奏云](https://www.lanzou.com/)
    - [x] [蓝奏云优享版](https://www.ilanzou.com/)
    - [x] [阿里云盘分享](https://www.alipan.com/)
    - [x] [谷歌相册](https://photos.google.com/)
    - [x] [Mega.nz](https://mega.nz)
    - [x] [一刻相册](https://photo.baidu.com/)
    - [x] SMB
    - [x] [115](https://115.com/)
    - [X] Cloudreve
    - [x] [Dropbox](https://www.dropbox.com/)
    - [x] [飞机盘](https://www.feijipan.com/)
    - [x] [多吉云](https://www.dogecloud.com/product/oss)
- [x] 部署方便，开箱即用
- [x] 文件预览（PDF、markdown、代码、纯文本……）
- [x] 画廊模式下的图像预览
- [x] 视频和音频预览，支持歌词和字幕
- [x] Office 文档预览（docx、pptx、xlsx、...）
- [x] `README.md` 预览渲染
- [x] 文件永久链接复制和直接文件下载
- [x] 黑暗模式
- [x] 国际化
- [x] 受保护的路由（密码保护和身份验证）
- [x] WebDav (详细文档待补充)
- [ ] Docker 部署（重建中）
- [x] Cloudflare workers 中转
- [x] 文件/文件夹打包下载
- [x] 网页上传(可以允许访客上传)，删除，新建文件夹，重命名，移动，复制
- [x] 离线下载
- [x] 跨存储复制文件
- [x] 单线程下载/串流的多线程下载加速

## 文档

<https://openlistteam.github.io/docs/> (临时)

## Demo

N/A

## 讨论

一般问题请到[讨论论坛](https://github.com/OpenListTeam/OpenList/discussions) ，**issue仅针对错误报告和功能请求。**

## 贡献者

感谢这些开源作者们：

[![Contributors](https://contrib.rocks/image?repo=OpenListTeam/OpenList)](https://github.com/OpenListTeam/OpenList/graphs/contributors)

## 许可

`OpenList` 是在 AGPL-3.0 许可下许可的开源软件。

## 免责声明
- 本程序为免费开源项目，旨在分享网盘文件，方便下载以及学习golang，使用时请遵守相关法律法规，请勿滥用；
- 本程序通过调用官方sdk/接口实现，无破坏官方接口行为；
- 本程序仅做302重定向/流量转发，不拦截、存储、篡改任何用户数据；
- 在使用本程序之前，你应了解并承担相应的风险，包括但不限于账号被ban，下载限速等，与本程序无关；
- 如有侵权，请联系[OpenListTeam](https://github.com/OpenListTeam)，团队会及时处理。

---

> [@GitHub](https://github.com/OpenListTeam) · [@Telegram 交流群]( https://t.me/OpenListTeam)
