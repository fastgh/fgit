# fgit
<img src="https://github.com/fastgh/fgit/blob/master/doc/how_fast_annotated.png" alt="image"/>

fgit是一个可以无缝替换git命令行的工具，加速对github.com的git clone/pull/push/update/fetch；fgit实时访问github.com，不是缓存。

## 特点：
  - 目前实测git clone最快时速度超过15MB/S
  - 支持github.com私有库，也支持push
  - 支持包括clone/push/pull/fetch在内的各种git命令，兼容git命令行参数，可以用来无缝替换git命令行
  - 只针对github.com加速，不干扰对非github.com库的使用
  - 两种工作模式（镜像模式和HTTP代理模式）自动切换：执行git命令前自动切换成使用代理或镜像设置，执行完成后自动移除代理和镜像设置
  - 使用go语言开发，不是shell脚本或.bat，跨平台。Windows 10、Linux (Unbuntu)、Mac (x86和M1)都实测通过

## 使用：

   - ```diff
     + 和常规的git命令行几乎相同，支持各种命令行参数，也就是说，平时git命令行怎么用，fgit就怎么用，区别只是git换成了fgit。
     ```

   - 对于公共库，clone/pull/fetch时默认使用镜像模式，基于安全考虑（镜像模式下会镜像服务器能探知到用户名和密码），镜像模式不支持push以及私有库

     镜像模式例如：`fgit clone https://github.com/spring-projects/spring-boot.git --depth=1`

   - 两种情况会判定为需要push或者私有库，此时会使用HTTPs代理模式，代理服务器是HTTPs透明代理，无法探知到用户名和密码，所以代理模式比镜像模式安全

      1. push

      2. URL中包含用户名，那么会被判定为需要push或者私有库。
         对于clone命令，URL是从clone的URL中解析得到，对于其它git命令，则使用`git remote -v`得到

     代理模式例如：

      1. `fgit push origin master`

      2. `fgit clone https://fastgh@github.com/fastgh/fgit.git`

     也可以通过`--use-proxy`选项强制走HTTP代理模式，例如:

      1. `fgit --use-proxy clone https://github.com/fastgh/fgit.git`

     代理服务器的线路成本高，所以大家尽量使用镜像模式，以节省服务器带宽资源

   - 其它功能：

     1. 可以打开调试开关，看一看fgit的工作过程：查询镜像/代理服务器列表 --> 设置镜像或代理 --> 执行git --> 恢复镜像或代理设置

        `fgit --debug clone https://github.com/fastgh/fgit.git`

     2. fgit首次运行时，会在用户主目录下生成一个配置文件.fgit.json，包含服务器地址等信息，必要时可以通过设置这个文件选择接入其它服务方，或指定镜像服务器或代理服务器

## 安装:

  下载页面：[https://github.com/fastgh/fgit/releases](https://github.com/fastgh/fgit/releases)。
  也可以自己编译安装（方法见后面小节）

  - Windows: [https://github.com/fastgh/fgit/releases/download/v1.1.0/fgit.exe](https://github.com/fastgh/fgit/releases/download/v1.1.0/fgit.exe)，下载后把它加入系统路径环境变量

  - Mac和Linux:

    ```shell
       sudo curl -L https://github.com/fastgh/fgit/releases/download/v1.1.0/fgit.$(echo `uname -s` | tr A-Z a-z).$(uname -m) -o /usr/local/bin/fgit
       sudo chmod +x /usr/local/bin/fgit
    ```

## 工作原理：
  git clone为什么慢成每秒几个KB？DNS污染和‘墙’等流传的原因都已经过时了，现在的主要原因是由于我们的访问通常走的是糟糕的国际线路。所以，解决起来其实也很简单直接，fgit使用的镜像服务器和HTTP代理服务器需要接入优化线路。

  知乎上这篇文章做了比较好的解释：[git clone一个github上的仓库，太慢，经常连接失败，但是github官网流畅访问，为什么？](https://www.zhihu.com/question/27159393)

  但是，毕竟那样设置起来比较麻烦，容易出错，而且，如果走的是github.com.cnpmjs.org等公共镜像，只能clone，不能访问私有库，也不能push。fgit把这些包装了起来，是这么做的：
  - 自动判断是使用镜像模式还是使用代理模式，默认情况下，只对git push命令和私有库启用代理模式。
  - 调用git命令前，如果判断是clone/push/pull/fetch，并且是github.com，那么设置镜像或者代理。
  - 调用git命令行（所以需要自己已经安装了git命令行）
  - 调用git命令后，取消镜像或者代理设置

  镜像服务器或者代理服务器的列表是动态更新，我们会不断更新这个列表。
  为防止滥用和被误`墙`，代理服务器使用token来控制连接有效期，默认是1分钟；fgit客户端会自动处理token的获取，对用户是透明的。
  代理服务器同时还做了限制、只允许访问github.com。

## 几种方案的对比：
|               | fgit | 公共镜像(https://github.com.cnpmjs.org等）| 自建代理或VPN | 导入国内GIT(码云等) | 缓存(https://gitclone.com) | 代下载(http://gitd.cc) |
| :------------ | :--- | :-------------------------------------: | :----------: | :---------------: | :----------------------: | :-------------------: |
| 浏览器访问     |      | [x]                                     | [x]          | [x]               |                           |                      |
| git命令行     | [x]  | [x]                                     | [x]          | [x]               | [x]                       |                      |
| 无需手工设置   | [x]  |                                         |              |                   |                           | [x]                  |
| 支持push      | [x]  |                                         | [x]          |                   |                           |                      |
| 支持私有库     | [x]  |                                         | [x]          |                   |                           |                      |
| 免费          | [x]  | [x]                                     |              | [x]               | [x]                       | [x]                  |
| 实时(非缓存)   | [x]  | [x]                                     | [x]          |                   |                           |                      |
| 不容易被误`墙` | [x]  | [x]                                     |              | [x]              |                           |                      |
| 动态更新服务器 | [x]   |                                        |              |                   |                           |                      |


 ### 题外话：

  准备发布1.0版时，发现有一个GITHUB项目做了差不多的事情：[https://github.com/FastGitOrg](https://github.com/FastGitOrg)，看了一下细节，FastGitOrg的缺点是：

  - FastGitOrg只支持镜像模式，通过配置成反向代理的nginx来提供镜像服务器，好像不能push（未实测），而且对于私有库和需要push的场景，反向代理是能解密所有数据的。
    和FastGitOrg相比，fgit使用的HTTP代理服务器只负责转发加密后的HTTPS流量，因为fgit的代理服务器没有github.com的SSL密钥，解密不了，所以更安全。

  - FastGitOrg只支持最基础的clone/pull命令，其它命令不支持。两相比较，fgit因为对全部的可能的git命令行参数/选项都做了解析和识别，所以能支持除了submodule以外的所有其它参数，因此可以透明的替换原git命令行，
    使用体验上和原git没什么区别，兼容其它会调用git命令行的工具场景。

  - FastGitOrg的镜像服务器地址是写死在客户端的，不是动态拉取，需要加镜像服务器时相对不容易，也容易被误`墙`。

  - FastGitOrg客户端不够健壮：它虽然也会在git执行结束后恢复原git设置，但在客户端意外崩溃退出时（譬如常见的Ctrl+C），就会留下镜像设置，就会导致访问
    非github库（譬如公司github库）时出错，需要用户手工修改git配置文件才能恢复。


## 限制：

   1. 不支持SSH协议
  
   2. 不支持sub module

   3. 代理模式下，执行clone时，会临时修改全局的.gitconfig文件（位于用户主目录下）；非clone时，会临时修改当前目录的.git/config文件。一般情况下，fgit会在执行结束后自动恢复原先的设置，即使崩溃退出时也能做到恢复，但不排除意外和有bug的情况，这时需手动修改.gitconfig或.git/config。


## 编译：

  - 安装GO语言开发环境，要求GO版本>=13

  - Linux / Mac环境下，`$ ./build.sh`，编译成功后可执行文件会生成在`build`目录下；Windows环境下类似，参考`build.sh`

## TODO：

  - 给clone命令增加--archive选项，用来下载后立刻打包clone下来的github项目

  - 增加download命令，用来下载github上的release文件

## 开源许可证

   GPL

## 问题反馈和交流：

   因为账号翻墙麻烦，所以几乎不看GITHUB Issues，所以有问题请加微信号：fastgithub

   <img src="https://control.fgit.kailash.cloud:10443/static/wechat_contact.png" alt="image" style="zoom:50%;"/>

