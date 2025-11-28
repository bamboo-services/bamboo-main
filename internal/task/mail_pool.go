/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package task

import (
	"crypto/tls"
	"errors"
	"net"
	"net/mail"
	"net/smtp"
	"sync"
	"time"

	"github.com/jordan-wright/email"
)

// TLSMailPool 支持 TLS 直连的邮件连接池
//
// 该连接池支持两种模式：
//   - TLS 直连（端口 465）：使用 tls.Dial 直接建立加密连接
//   - STARTTLS（端口 587）：先建立明文连接，再升级为加密连接
type TLSMailPool struct {
	addr        string           // SMTP 服务器地址（host:port）
	auth        smtp.Auth        // SMTP 认证信息
	tlsConfig   *tls.Config      // TLS 配置
	useTLS      bool             // true=TLS直连(465), false=STARTTLS(587)
	maxConn     int              // 最大连接数
	conns       chan *smtpClient // 连接池通道
	mu          sync.Mutex       // 互斥锁
	activeConns int              // 当前活跃连接数
	closing     chan struct{}    // 关闭信号
	closed      bool             // 是否已关闭
}

// smtpClient 封装 smtp.Client，添加失败计数
type smtpClient struct {
	*smtp.Client
	failCount int // 连续失败次数
}

// NewTLSMailPool 创建支持 TLS 的邮件连接池
//
// 参数说明:
//   - addr: SMTP 服务器地址（格式：host:port）
//   - maxConn: 最大连接数
//   - auth: SMTP 认证信息
//   - tlsConfig: TLS 配置
//   - useTLS: true 表示使用 TLS 直连（465端口），false 表示使用 STARTTLS（587端口）
//
// 返回值:
//   - TLS 邮件连接池实例
func NewTLSMailPool(addr string, maxConn int, auth smtp.Auth, tlsConfig *tls.Config, useTLS bool) *TLSMailPool {
	if maxConn <= 0 {
		maxConn = 4
	}

	return &TLSMailPool{
		addr:      addr,
		auth:      auth,
		tlsConfig: tlsConfig,
		useTLS:    useTLS,
		maxConn:   maxConn,
		conns:     make(chan *smtpClient, maxConn),
		closing:   make(chan struct{}),
	}
}

// buildConnection 创建新的 SMTP 连接
//
// 根据 useTLS 配置选择连接方式：
//   - TLS 直连：使用 tls.Dial 建立加密连接（适用于 465 端口）
//   - STARTTLS：使用 smtp.Dial 建立明文连接后升级（适用于 587 端口）
func (p *TLSMailPool) buildConnection() (*smtpClient, error) {
	var client *smtp.Client
	var err error

	if p.useTLS {
		// TLS 直连模式（465 端口）
		conn, dialErr := tls.Dial("tcp", p.addr, p.tlsConfig)
		if dialErr != nil {
			return nil, dialErr
		}

		// 从地址中提取主机名
		host, _, _ := net.SplitHostPort(p.addr)
		client, err = smtp.NewClient(conn, host)
		if err != nil {
			conn.Close()
			return nil, err
		}
	} else {
		// STARTTLS 模式（587 端口）
		client, err = smtp.Dial(p.addr)
		if err != nil {
			return nil, err
		}

		// 检查是否支持 STARTTLS
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(p.tlsConfig); err != nil {
				client.Close()
				return nil, err
			}
		}
	}

	// 执行 SMTP 认证
	if p.auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(p.auth); err != nil {
				client.Close()
				return nil, err
			}
		}
	}

	return &smtpClient{Client: client, failCount: 0}, nil
}

// getConnection 从连接池获取连接
//
// 获取策略：
//  1. 首先尝试从池中获取现有连接
//  2. 如果池中没有且未达到最大连接数，创建新连接
//  3. 如果已达到最大连接数，等待可用连接或超时
func (p *TLSMailPool) getConnection(timeout time.Duration) (*smtpClient, error) {
	// 检查连接池是否已关闭
	select {
	case <-p.closing:
		return nil, errors.New("连接池已关闭")
	default:
	}

	// 1. 尝试从池中获取现有连接
	select {
	case conn := <-p.conns:
		// 验证连接是否有效
		if err := conn.Noop(); err != nil {
			conn.Close()
			p.mu.Lock()
			p.activeConns--
			p.mu.Unlock()
			// 连接无效，创建新连接
			return p.createNewConnection()
		}
		return conn, nil
	default:
	}

	// 2. 尝试创建新连接
	p.mu.Lock()
	if p.activeConns < p.maxConn {
		p.activeConns++
		p.mu.Unlock()

		conn, err := p.buildConnection()
		if err != nil {
			p.mu.Lock()
			p.activeConns--
			p.mu.Unlock()
			return nil, err
		}
		return conn, nil
	}
	p.mu.Unlock()

	// 3. 等待可用连接
	select {
	case conn := <-p.conns:
		// 验证连接
		if err := conn.Noop(); err != nil {
			conn.Close()
			p.mu.Lock()
			p.activeConns--
			p.mu.Unlock()
			return p.createNewConnection()
		}
		return conn, nil
	case <-time.After(timeout):
		return nil, errors.New("获取连接超时")
	case <-p.closing:
		return nil, errors.New("连接池已关闭")
	}
}

// createNewConnection 创建新连接（内部方法）
func (p *TLSMailPool) createNewConnection() (*smtpClient, error) {
	p.mu.Lock()
	p.activeConns++
	p.mu.Unlock()

	conn, err := p.buildConnection()
	if err != nil {
		p.mu.Lock()
		p.activeConns--
		p.mu.Unlock()
		return nil, err
	}
	return conn, nil
}

// releaseConnection 释放连接回连接池
//
// 如果发送过程中出现错误，会增加失败计数。
// 连续失败超过 3 次的连接会被关闭而不是放回池中。
func (p *TLSMailPool) releaseConnection(conn *smtpClient, sendErr error) {
	if conn == nil {
		return
	}

	// 如果发送失败，增加失败计数
	if sendErr != nil {
		conn.failCount++
		// 连续失败超过 3 次，关闭连接
		if conn.failCount >= 3 {
			conn.Close()
			p.mu.Lock()
			p.activeConns--
			p.mu.Unlock()
			return
		}
	} else {
		conn.failCount = 0
	}

	// 重置连接状态
	if err := conn.Reset(); err != nil {
		conn.Close()
		p.mu.Lock()
		p.activeConns--
		p.mu.Unlock()
		return
	}

	// 尝试放回连接池
	select {
	case p.conns <- conn:
		// 成功放回
	default:
		// 池已满，关闭连接
		conn.Close()
		p.mu.Lock()
		p.activeConns--
		p.mu.Unlock()
	}
}

// Send 使用连接池发送邮件
//
// 参数说明:
//   - e: 邮件对象
//   - timeout: 发送超时时间
//
// 返回值:
//   - 发送错误（如果有）
func (p *TLSMailPool) Send(e *email.Email, timeout time.Duration) error {
	conn, err := p.getConnection(timeout)
	if err != nil {
		return err
	}

	var sendErr error
	defer func() {
		p.releaseConnection(conn, sendErr)
	}()

	// 解析发件人地址
	from, err := mail.ParseAddress(e.From)
	if err != nil {
		sendErr = err
		return err
	}

	// 收集所有收件人
	var recipients []string
	for _, addr := range e.To {
		parsed, err := mail.ParseAddress(addr)
		if err != nil {
			sendErr = err
			return err
		}
		recipients = append(recipients, parsed.Address)
	}
	for _, addr := range e.Cc {
		parsed, err := mail.ParseAddress(addr)
		if err != nil {
			sendErr = err
			return err
		}
		recipients = append(recipients, parsed.Address)
	}
	for _, addr := range e.Bcc {
		parsed, err := mail.ParseAddress(addr)
		if err != nil {
			sendErr = err
			return err
		}
		recipients = append(recipients, parsed.Address)
	}

	// SMTP 发送流程
	if err := conn.Mail(from.Address); err != nil {
		sendErr = err
		return err
	}

	for _, recip := range recipients {
		if err := conn.Rcpt(recip); err != nil {
			sendErr = err
			return err
		}
	}

	w, err := conn.Data()
	if err != nil {
		sendErr = err
		return err
	}

	msg, err := e.Bytes()
	if err != nil {
		sendErr = err
		return err
	}

	if _, err := w.Write(msg); err != nil {
		sendErr = err
		return err
	}

	if err := w.Close(); err != nil {
		sendErr = err
		return err
	}

	return nil
}

// Close 关闭连接池
//
// 关闭所有活跃连接并释放资源
func (p *TLSMailPool) Close() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	close(p.closing)
	p.mu.Unlock()

	// 关闭池中所有连接
	close(p.conns)
	for conn := range p.conns {
		if conn != nil {
			conn.Close()
		}
	}
}
