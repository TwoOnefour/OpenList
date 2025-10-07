package handles

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/fs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/net"
	"github.com/OpenListTeam/OpenList/v4/internal/setting"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"github.com/yuin/goldmark"
	"golang.org/x/net/context"
	"io"
	stdnet "net"
	"net/http"
	stdpath "path"
	"strconv"
	"strings"
)

func Down(c *gin.Context) {
	rawPath := c.Request.Context().Value(conf.PathKey).(string)
	filename := stdpath.Base(rawPath)
	storage, err := fs.GetStorage(rawPath, &fs.GetStoragesArgs{})
	if err != nil {
		common.ErrorPage(c, err, 500)
		return
	}
	if common.ShouldProxy(storage, filename) {
		Proxy(c)
		return
	} else {
		link, _, err := fs.Link(c.Request.Context(), rawPath, model.LinkArgs{
			IP:       c.ClientIP(),
			Header:   c.Request.Header,
			Type:     c.Query("type"),
			Redirect: true,
		})
		if err != nil {
			common.ErrorPage(c, err, 500)
			return
		}
		redirect(c, link)
	}
}

func clientGone(err error) bool {
	if err == nil {
		return false
	}
	// 各平台/库的典型文案（避免遗漏）
	return errors.Is(err, context.Canceled) ||
		errors.Is(err, stdnet.ErrClosed) ||
		strings.Contains(err.Error(), "client disconnected") ||
		strings.Contains(err.Error(), "broken pipe") ||
		strings.Contains(err.Error(), "connection reset by peer") ||
		strings.Contains(err.Error(), "stream closed")
}

func Proxy(c *gin.Context) {
	rawPath := c.Request.Context().Value(conf.PathKey).(string)
	filename := stdpath.Base(rawPath)
	storage, err := fs.GetStorage(rawPath, &fs.GetStoragesArgs{})
	if err != nil {
		common.ErrorPage(c, err, 500)
		return
	}
	if canProxy(storage, filename) {
		if _, ok := c.GetQuery("d"); !ok {
			if url := common.GenerateDownProxyURL(storage.GetStorage(), rawPath); url != "" {
				c.Redirect(302, url)
				return
			}
		}
		link, file, err := fs.Link(c.Request.Context(), rawPath, model.LinkArgs{
			Header: c.Request.Header,
			Type:   c.Query("type"),
		})
		if err != nil {
			common.ErrorPage(c, err, 500)
			return
		}
		proxy(c, link, file, storage.GetStorage().ProxyRange)
	} else {
		common.ErrorPage(c, errors.New("proxy not allowed"), 403)
		return
	}
}

func redirect(c *gin.Context, link *model.Link) {
	defer link.Close()
	var err error
	c.Header("Referrer-Policy", "no-referrer")
	c.Header("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
	if setting.GetBool(conf.ForwardDirectLinkParams) {
		query := c.Request.URL.Query()
		for _, v := range conf.SlicesMap[conf.IgnoreDirectLinkParams] {
			query.Del(v)
		}
		link.URL, err = utils.InjectQuery(link.URL, query)
		if err != nil {
			common.ErrorPage(c, err, 500)
			return
		}
	}
	c.Redirect(302, link.URL)
}

func proxy(c *gin.Context, link *model.Link, file model.Obj, proxyRange bool) {

	var err error
	if link.URL != "" && setting.GetBool(conf.ForwardDirectLinkParams) {
		query := c.Request.URL.Query()
		for _, v := range conf.SlicesMap[conf.IgnoreDirectLinkParams] {
			query.Del(v)
		}
		link.URL, err = utils.InjectQuery(link.URL, query)
		if err != nil {
			common.ErrorPage(c, err, 500)
			link.Close()
			return
		}
	}
	if proxyRange {
		link = common.ProxyRange(c, link, file.GetSize())
	}
	defer func(link *model.Link) {
		err := link.Close()
		if err != nil {
			common.ErrorPage(c, err, http.StatusBadGateway, true)
		}
	}(link)
	Writer := &common.WrittenResponseWriter{ResponseWriter: c.Writer}
	raw, _ := strconv.ParseBool(c.DefaultQuery("raw", "false"))
	if utils.Ext(file.GetName()) == "md" && setting.GetBool(conf.FilterReadMeScripts) && !raw {
		buf := bytes.NewBuffer(make([]byte, 0, file.GetSize()))
		w := &common.InterceptResponseWriter{ResponseWriter: Writer, Writer: buf}
		err = common.Proxy(w, c.Request, link, file)
		if err == nil && buf.Len() > 0 {
			if c.Writer.Status() < 200 || c.Writer.Status() > 300 {
				c.Writer.Write(buf.Bytes())
				return
			}

			var html bytes.Buffer
			if err = goldmark.Convert(buf.Bytes(), &html); err != nil {
				err = fmt.Errorf("markdown conversion failed: %w", err)
			} else {
				buf.Reset()
				err = bluemonday.UGCPolicy().SanitizeReaderToWriter(&html, buf)
				if err == nil {
					Writer.Header().Set("Content-Length", strconv.FormatInt(int64(buf.Len()), 10))
					Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
					_, err = utils.CopyWithBuffer(Writer, buf)
				}
			}
		}
	} else {
		err = common.Proxy(Writer, c.Request, link, file)
	}
	if err == nil {
		return
	}

	// 先把“已开始写”的场景放大兜住
	alreadyStarted := c.Writer.Written() || Writer.IsWritten() ||
		c.Writer.Status() == http.StatusPartialContent ||
		c.Writer.Header().Get("Content-Range") != ""

	// 客户端/前置取消（包含 EOF 类）→ 静默结束
	if alreadyStarted || clientGone(err) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		log.Infof("%s %s client gone/EOF after start: %v", c.Request.Method, c.Request.URL.Path, err)
		return
	}

	// 真正的上游错误才下发错误页；默认 502
	sc := 502
	var hse net.HttpStatusCodeError
	if errors.As(err, &hse) || errors.As(errors.Unwrap(err), &hse) {
		sc = int(hse)
	}
	common.ErrorPage(c, err, sc, true)
}

// TODO need optimize
// when can be proxy?
// 1. text file
// 2. config.MustProxy()
// 3. storage.WebProxy
// 4. proxy_types
// solution: text_file + shouldProxy()
func canProxy(storage driver.Driver, filename string) bool {
	if storage.Config().MustProxy() || storage.GetStorage().WebProxy || storage.GetStorage().WebdavProxyURL() {
		return true
	}
	if utils.SliceContains(conf.SlicesMap[conf.ProxyTypes], utils.Ext(filename)) {
		return true
	}
	if utils.SliceContains(conf.SlicesMap[conf.TextTypes], utils.Ext(filename)) {
		return true
	}
	return false
}
