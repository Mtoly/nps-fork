# 增强功能

## 1. NAT 类型检测

使用 STUN 服务器检测 NAT 类型：
```bash
./npc nat -stun_addr=stun.stunprotocol.org:3478
```
如果 **P2P 双方都是 `Symmetric NAT`** ，则 **无法穿透**，其他 NAT 组合通常可以成功。

📌 **可选参数**

| 参数           | 说明            | 默认值                          |
|--------------|---------------|------------------------------|
| `-stun_addr` | 指定 STUN 服务器地址 | `stun.stunprotocol.org:3478` |

---

## 2. 状态检查

检查 NPC 客户端的运行状态：
```bash
./npc status -config=/path/to/npc.conf
```
📌 **可选参数**

| 参数        | 说明            |
|-----------|---------------|
| `-config` | 指定 NPC 配置文件路径 |

---

## 3. 重载配置文件

重新加载 NPC 客户端配置，而无需重启进程：
```bash
./npc restart -config=/path/to/npc.conf
```
📌 **可选参数**

| 参数        | 说明            |
|-----------|---------------|
| `-config` | 指定 NPC 配置文件路径 |

---

## 4. 通过代理连接 NPS

如果 NPC 运行的机器无法直接访问外网，可以通过 **Socks5 / HTTP 代理** 连接 NPS 服务器。

### **4.1 配置文件方式**
在 `npc.conf` 文件中添加：
```ini
[common]
proxy_url=socks5://111:222@127.0.0.1:8024
```

### **4.2 命令行方式**
```bash
./npc -server=xxx:123 -vkey=xxx -proxy=socks5://111:222@127.0.0.1:8024
```

📌 **支持代理协议**

| 代理类型       | 示例格式                                 |
|------------|--------------------------------------|
| **Socks5** | `socks5://username:password@ip:port` |
| **HTTP**   | `http://username:password@ip:port`   |

---

## 5. 其他命令行参数
📌 **所有参数可与启动命令组合使用** ：

```bash
./npc -server=xxx:123,yyy:456 -vkey=xxx,yyy -type=tls,tcp -log=off -debug=false
```

| 参数                    | 说明                                                     | 默认值       |
|-----------------------|--------------------------------------------------------|-----------|
| `-server`             | 指定 NPS 服务器地址（`ip:port`）                                | 无         |
| `-vkey`               | 客户端认证密钥                                                | 无         |
| `-type`               | 服务器连接方式（`tcp` / `tls` / `kcp` / `quic` / `ws` / `wss`） | `tcp`     |
| `-config`             | 指定配置文件路径                                               | 无         |
| `-proxy`              | 通过代理连接 NPS（支持 Socks5 / HTTP）                           | 无         |
| `-debug`              | 是否启用调试模式                                               | `true`    |
| `-log`                | 日志输出模式（`stdout` / `file` / `both` / `off`）             | `file`    |
| `-log_path`           | NPC 日志路径（为空使用默认路径，`off` 禁用日志）                          | `npc.log` |
| `-log_level`          | 日志级别（trace、debug、info、warn、error、fatal、panic、off）      | `trace`   |
| `-log_compress`       | 是否启用日志压缩                                               | `false`   |
| `-log_max_days`       | 日志最大保留天数（0 关闭）                                         | `7`       |
| `-log_max_files`      | 最大日志文件数（0 关闭）                                          | `10`      |
| `-log_max_size`       | 单个日志文件最大大小（MB）                                         | `5`       |
| `-log_color`          | 控制台输出启用 ANSI 彩色                                        | `true`    |
| `-auto_reconnect`     | 断线后自动重连                                                | `true`    |
| `-disconnect_timeout` | 连接超时秒数                                                 | `60`      |
| `-keepalive`          | 保活（KeepAlive）周期（秒）                                     | 默认        |
| `-pprof`              | 启用 PProf 调试（格式 `ip:port`）                              | 无         |
| `-local_type`         | P2P 目标类型                                               | `p2p`     |
| `-local_port`         | P2P 本地端口                                               | `2000`    |
| `-password`           | P2P 认证密码                                               | 无         |
| `-target`             | P2P 目标（如目标客户端 ID/标识）                                   | 无         |
| `-target_type`        | P2P 目标连接类型（`all` / `tcp` / `udp`）                      | `all`     |
| `-p2p_timeout`        | P2P 超时时间（秒）                                            | `5`       |
| `-p2p_type`           | P2P 连接类型（`quic` / `kcp`）                               | `quic`    |
| `-disable_p2p`        | 禁用 P2P 连接                                              | `false`   |
| `-fallback_secret`    | 当 P2P 不可用时回退到 Secret 直连模式                              | `true`    |
| `-stun_addr`          | STUN 服务器地址                                             | 无         |
| `-dns_server`         | 配置 DNS 服务器                                             | `8.8.8.8` |
| `-ntp_server`         | 配置 NTP 服务器                                             | 无         |
| `-ntp_interval`       | NTP 最小查询间隔（分钟）                                         | `5`       |
| `-timezone`           | 配置时区（Asia/Shanghai）                                    | 无         |
| `-time`               | 客户端注册时间（小时）                                            | `2`       |
| `-gen2fa`             | 生成 TOTP 双因素认证密钥                                        | `false`   |
| `-get2fa`             | 根据提供的密钥输出一次性 TOTP 验证码                                  | 无         |
| `-version`            | 显示当前版本                                                 | 无         |

下面只补充**文档中缺失**的参数说明（基于 `./npc -h`），其余已存在的不重复列出：

---

## 6. 群晖支持

📌 **推荐使用 Docker 部署**
```bash
docker pull duan2001/npc
docker run -d --restart=always --name npc --net=host duan2001/npc -server=xxx:123,yyy:456 -vkey=xxx,yyy -type=tls,tcp -log=off
```
~~曾提供 `.spk` 群晖套件，但已不再维护，建议使用 Docker 方式运行。~~ 
✅[Telegram](https://t.me/npsdev) 内有第三方提供的群晖套件。

---

✅ **如需更多帮助，请查看 [文档](https://github.com/djylb/nps) 或提交 [GitHub Issues](https://github.com/djylb/nps/issues) 反馈问题。**