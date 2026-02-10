package common

import (
	"e-woms/conf"
	"e-woms/controllers/backend"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

// 允许的文件扩展名白名单
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".heic": true, // 苹果手机照片格式
	".heif": true, // 苹果手机照片格式（备用）
	".webp": true, // 现代浏览器常用格式
	".pdf":  true,
	".ppt":  true,
	".pptx": true,
}

// 允许的 MIME 类型白名单
var allowedMimeTypes = map[string]bool{
	"image/jpeg":                    true,
	"image/jpg":                     true,
	"image/png":                     true,
	"image/gif":                     true,
	"image/heic":                    true, // 苹果 HEIC 格式
	"image/heif":                    true, // 苹果 HEIF 格式
	"image/webp":                    true, // WebP 格式
	"application/pdf":               true,
	"application/vnd.ms-powerpoint": true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
}

// UploadController 文件上传控制器
type UploadController struct {
	backend.BaseController
}

// Upload 文件上传
// @Summary 文件上传
// @Title 文件上传
// @Description 上传文件到服务器，支持图片(jpg/jpeg/png/gif/heic/heif/webp)和课件(pdf/ppt/pptx)，大小限制20MB
// @Tags 通用-文件上传
// @Accept json
// @Produce json
// @Param file formData file true "上传的文件"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{"url":"/static/upload/xxx.jpg","filename":"xxx.jpg","size":102400,"ext":".jpg","original":"original.jpg","time":1234567890}}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/common/upload [post]
func (c *UploadController) Upload() {
	logs.Info("[UploadController][Upload] 开始处理文件上传")

	// 1. 获取上传的文件
	file, header, err := c.GetFile("file")
	if err != nil {
		logs.Error("[UploadController][Upload] 获取文件失败: %v", err)
		c.Error(conf.PARAMS_ERROR, "请上传文件")
		return
	}
	defer file.Close()

	// 2. 验证文件大小（20MB = 20 * 1024 * 1024 bytes）
	maxSize := int64(20 * 1024 * 1024)
	if header.Size > maxSize {
		logs.Error("[UploadController][Upload] 文件大小超过限制: %d bytes", header.Size)
		c.Error(conf.PARAMS_ERROR, "文件大小超过限制（最大20MB）")
		return
	}

	// 3. 验证文件扩展名（白名单）
	originalFilename := filepath.Base(header.Filename) // 防止路径遍历攻击
	ext := strings.ToLower(filepath.Ext(originalFilename))
	if !allowedExtensions[ext] {
		logs.Error("[UploadController][Upload] 不支持的文件类型: %s", ext)
		c.Error(conf.PARAMS_ERROR, "不支持的文件类型，仅支持: jpg, jpeg, png, gif, heic, heif, webp, pdf, ppt, pptx")
		return
	}

	// 4. 验证文件 MIME 类型（防止文件伪装）
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		logs.Error("[UploadController][Upload] 读取文件内容失败: %v", err)
		c.Error(conf.SERVER_ERROR, "读取文件失败")
		return
	}
	// 重置文件指针到开头，以便后续保存完整文件
	_, err = file.Seek(0, 0)
	if err != nil {
		logs.Error("[UploadController][Upload] 重置文件指针失败: %v", err)
		c.Error(conf.SERVER_ERROR, "文件处理失败")
		return
	}

	// 检测文件实际 MIME 类型
	contentType := http.DetectContentType(buffer)

	// 特殊处理：pptx/docx/xlsx 等 Office 文件本质上是 ZIP 压缩包
	// Go 的 DetectContentType 无法识别它们的正确 MIME 类型，会返回 application/zip
	// 因此需要根据扩展名进行特殊判断
	if contentType == "application/zip" {
		// 如果检测到是 ZIP，但扩展名是 Office 文件，则允许通过
		if ext == ".pptx" || ext == ".ppt" {
			logs.Info("[UploadController][Upload] 检测到 Office 文件（ZIP 格式）: %s", ext)
			contentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
		} else {
			// 其他情况下的 ZIP 文件不允许
			logs.Error("[UploadController][Upload] 不允许的 ZIP 文件: %s", ext)
			c.Error(conf.PARAMS_ERROR, "不支持 ZIP 压缩文件")
			return
		}
	}

	if !allowedMimeTypes[contentType] {
		logs.Error("[UploadController][Upload] 文件内容类型不匹配: %s", contentType)
		c.Error(conf.PARAMS_ERROR, "文件内容类型不匹配，可能是伪装文件")
		return
	}

	// 5. 生成安全的文件名：时间戳（14位）+ 随机3位数字 + 原扩展名
	// 时间格式：20251223143025（YYYYMMDDHHmmss）
	timestamp := time.Now().Format("20060102150405")
	// 生成随机3位数字：000-999
	randomNum := rand.Intn(1000)
	// 组合文件名（使用验证过的扩展名）
	filename := fmt.Sprintf("%s%03d%s", timestamp, randomNum, ext)

	// 6. 确保保存目录存在
	uploadDir := "static/upload"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logs.Error("[UploadController][Upload] 创建目录失败: %v", err)
		c.Error(conf.SERVER_ERROR, "创建保存目录失败")
		return
	}

	// 7. 保存文件
	savePath := filepath.Join(uploadDir, filename)
	if err := c.SaveToFile("file", savePath); err != nil {
		logs.Error("[UploadController][Upload] 保存文件失败: %v", err)
		c.Error(conf.SERVER_ERROR, "保存文件失败")
		return
	}

	// 8. 返回文件信息
	logs.Info("[UploadController][Upload] 文件上传成功: %s, 原始文件: %s, 大小: %d bytes, MIME: %s",
		filename, originalFilename, header.Size, contentType)
	c.Success(map[string]interface{}{
		"url":      "/" + savePath,    // 返回相对路径：/static/upload/xxx.jpg
		"filename": filename,          // 保存的文件名
		"size":     header.Size,       // 文件大小（字节）
		"ext":      ext,               // 文件扩展名
		"original": originalFilename,  // 原始文件名（已清理）
		"time":     time.Now().Unix(), // 上传时间戳
	})
}
