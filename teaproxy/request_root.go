package teaproxy

import (
	"fmt"
	"github.com/dchest/siphash"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 调用本地静态资源
func (this *Request) callRoot(writer *ResponseWriter) error {
	if len(this.uri) == 0 {
		this.notFoundError(writer)
		return nil
	}

	if !filepath.IsAbs(this.root) {
		this.root = Tea.Root + Tea.DS + this.root
	}

	requestPath := this.uri
	uri, err := url.ParseRequestURI(this.uri)
	query := ""
	if err == nil {
		requestPath = uri.Path
		query = uri.RawQuery
	}

	// 去掉其中的奇怪的路径
	requestPath = strings.Replace(requestPath, "..\\", "", -1)

	if requestPath == "/" {
		// 根目录
		indexFile := this.findIndexFile(this.root)
		if len(indexFile) > 0 {
			this.uri = requestPath + indexFile
			if len(query) > 0 {
				this.uri += "?" + query
			}
			err := this.configure(this.server, 0)
			if err != nil {
				logs.Error(err)
				this.addError(err)
				this.serverError(writer)
				return nil
			}
			return this.call(writer)
		} else {
			this.notFoundError(writer)
			return nil
		}
	}
	filename := strings.Replace(requestPath, "/", Tea.DS, -1)
	filePath := ""
	if filename[0:1] == Tea.DS {
		filePath = this.root + filename
	} else {
		filePath = this.root + Tea.DS + filename
	}

	this.filePath = filePath

	stat, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			this.notFoundError(writer)
			return nil
		} else {
			this.serverError(writer)
			logs.Error(err)
			this.addError(err)
			return nil
		}
	}
	if stat.IsDir() {
		indexFile := this.findIndexFile(filePath)
		if len(indexFile) > 0 {
			this.uri = requestPath + indexFile
			if len(query) > 0 {
				this.uri += "?" + query
			}
			err := this.configure(this.server, 0)
			if err != nil {
				logs.Error(err)
				this.serverError(writer)
				this.addError(err)
				return nil
			}
			return this.call(writer)
		} else {
			this.notFoundError(writer)
			return nil
		}
	}

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// 响应header
	respHeader := writer.Header()

	// mime type
	if !hasIgnoreHeaders || !ignoreHeaders.Has("CONTENT-TYPE") {
		ext := filepath.Ext(requestPath)
		if len(ext) > 0 {
			mimeType := mime.TypeByExtension(ext)
			if len(mimeType) > 0 {
				if _, found := textMimeMap[mimeType]; found {
					if len(this.charset) > 0 {
						// 去掉里面的charset设置
						index := strings.Index(mimeType, "charset=")
						if index > 0 {
							respHeader.Set("Content-Type", mimeType[:index+len("charset=")]+this.charset)
						} else {
							respHeader.Set("Content-Type", mimeType+"; charset="+this.charset)
						}
					} else {
						respHeader.Set("Content-Type", mimeType)
					}
				} else {
					respHeader.Set("Content-Type", mimeType)
				}
			}
		}
	}

	// length
	respHeader.Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

	// 支持 Last-Modified
	modifiedTime := stat.ModTime().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	if len(respHeader.Get("Last-Modified")) == 0 {
		respHeader.Set("Last-Modified", modifiedTime)
	}

	// 支持 ETag
	eTag := "\"et" + fmt.Sprintf("%0x", siphash.Hash(0, 0, []byte(filename+strconv.FormatInt(stat.ModTime().UnixNano(), 10)+strconv.FormatInt(stat.Size(), 10)))) + "\""
	if len(respHeader.Get("ETag")) == 0 {
		respHeader.Set("ETag", eTag)
	}

	// proxy callback
	if this.responseCallback != nil {
		this.responseCallback(writer)
	}

	// 支持 If-None-Match
	if this.requestHeader("If-None-Match") == eTag {
		// 自定义Header
		this.WriteResponseHeaders(writer, http.StatusNotModified)

		writer.WriteHeader(http.StatusNotModified)

		return nil
	}

	// 支持 If-Modified-Since
	if this.requestHeader("If-Modified-Since") == modifiedTime {
		// 自定义Header
		this.WriteResponseHeaders(writer, http.StatusNotModified)

		writer.WriteHeader(http.StatusNotModified)

		return nil
	}

	// 自定义Header
	this.WriteResponseHeaders(writer, http.StatusOK)

	var contentReader io.Reader = nil
	isOpen := false
	if this.server.CacheStatic {
		reader, shouldClose, err := ShareStaticDelivery.Read(filePath, stat)
		if err != nil {
			this.serverError(writer)
			logs.Error(err)
			this.addError(err)
			return nil
		}
		contentReader = reader
		if shouldClose {
			defer contentReader.(*os.File).Close()
		}
	} else {
		reader, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
		if err != nil {
			this.serverError(writer)
			logs.Error(err)
			this.addError(err)
			return nil
		}
		contentReader = reader
		isOpen = true
	}

	writer.Prepare(stat.Size())
	buf := make([]byte, 1024) // TODO buffer size应该可以设置，或者根据stat.Size()动态调整
	_, err = io.CopyBuffer(writer, contentReader, buf)

	// 不使用defer，以便于加快速度
	if isOpen {
		contentReader.(*os.File).Close()
	}

	if err != nil {
		if this.debug {
			logs.Error(err)
		}
		return nil
	}

	return nil
}
