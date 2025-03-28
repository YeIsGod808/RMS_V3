package application

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/RMS_V3/config"
	"github.com/RMS_V3/internal/kg/repository"
	"github.com/RMS_V3/log"
	"github.com/RMS_V3/middleware/minioStore"
	"github.com/RMS_V3/middleware/neo4jUtils"
	"github.com/RMS_V3/pkg/response"
	"github.com/RMS_V3/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"go.uber.org/zap"
)

type UploadResourceRequest struct {
	Title              string  `form:"title" binding:"required,min=1,max=100"`
	Description        *string `form:"description" binding:"omitempty,max=65535"`             // 假设text类型最大长度为65535
	Difficulty         *string `form:"difficulty" binding:"omitempty,oneof=easy medium hard"` // 只有exercise需要
	ResourceType       string  `form:"resource_type" binding:"required,oneof=video courseware exercise"`
	KnowledgePointName string  `form:"point_name" binding:"required,min=1,max=100"`
	ResourceLink       *string `form:"resource_link" binding:"omitempty,url"` // 可选，如果是链接则使用此字段
}

func UploadResource(c *gin.Context) {
	var req UploadResourceRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, response.Error(400, "参数验证失败"))
		return
	}

	// 判断是文件上传还是链接上传
	if req.ResourceLink != nil && *req.ResourceLink != "" {
		// 处理链接上传逻辑
		handleLinkUpload(c, req)
		return
	}

	// 处理文件上传逻辑
	handleFileUpload(c, req)
}

// 处理文件上传
func handleFileUpload(c *gin.Context, req UploadResourceRequest) {
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(400, response.Error(400, "无法读取文件"))
		return
	}

	// 查询知识点ID
	knowledgePointId, err := getKnowledgePointIdByName(c, req.KnowledgePointName)
	if err != nil {
		c.JSON(500, response.Error(500, "查询知识点ID失败"))
		return
	}

	// Extract the filename
	filename := filepath.Base(data.Filename)
	// Generate a final filename using a random string and the original filename to prevent overwriting
	finalName := fmt.Sprintf("%s_%s", utils.RandomString(), filename)

	var savePath string
	switch req.ResourceType {
	case "video":
		savePath = config.GetGlobalConfig().MinioConfig.VideoPath
	case "courseware", "exercise":
		savePath = config.GetGlobalConfig().MinioConfig.FilePath
	default:
		c.JSON(400, response.Error(400, "未知的资源类型"))
		return
	}

	// 确保保存目录存在
	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		log.Errorf("Failed to create directory %s: %v\n", savePath, err)
		c.JSON(500, response.Error(500, "创建保存目录失败"))
		return
	}

	saveFile := filepath.Join(savePath, finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}

	if _, err := os.Stat(saveFile); os.IsNotExist(err) {
		log.Errorf("File does not exist at path: %s", saveFile)
		c.JSON(500, response.Error(500, "文件未正确保存"))
		return
	}

	// 根据资源类型进行具体处理
	switch req.ResourceType {
	case "video":
		err := uploadVideo(knowledgePointId, saveFile, req.Title, req.Description)
		if err != nil {
			c.JSON(500, response.Error(500, err.Error()))
			return
		}
	case "courseware":
		err := uploadFile(knowledgePointId, saveFile, req.Title, req.Description, req.ResourceType)
		if err != nil {
			c.JSON(500, response.Error(500, err.Error()))
			return
		}
	case "exercise":
		if req.Difficulty == nil {
			c.JSON(400, response.Error(400, "难度不能为空"))
			return
		}
		err := uploadFile(knowledgePointId, saveFile, req.Title, req.Description, req.ResourceType, *req.Difficulty)
		if err != nil {
			c.JSON(500, response.Error(500, err.Error()))
			return
		}
	}

	c.JSON(200, response.Success(gin.H{
		"message": "资源上传成功",
		"id":      knowledgePointId,
	}))
}

// 处理链接上传
func handleLinkUpload(c *gin.Context, req UploadResourceRequest) {
	// 查询知识点ID
	knowledgePointId, err := getKnowledgePointIdByName(c, req.KnowledgePointName)
	if err != nil {
		c.JSON(500, response.Error(500, "查询知识点ID失败"))
		return
	}

	// 根据资源类型将链接写入数据库
	switch req.ResourceType {
	case "video":
		err := saveLinkToDB(knowledgePointId, *req.ResourceLink, req.Title, "video", req.Description)
		if err != nil {
			c.JSON(500, response.Error(500, err.Error()))
			return
		}
	case "courseware":
		err := saveLinkToDB(knowledgePointId, *req.ResourceLink, req.Title, "courseware", req.Description)
		if err != nil {
			c.JSON(500, response.Error(500, err.Error()))
			return
		}
	case "exercise":
		if req.Difficulty == nil {
			c.JSON(400, response.Error(400, "难度不能为空"))
			return
		}
		err := saveLinkToDB(knowledgePointId, *req.ResourceLink, req.Title, "exercise", req.Description, *req.Difficulty)
		if err != nil {
			c.JSON(500, response.Error(500, err.Error()))
			return
		}
	default:
		c.JSON(400, response.Error(400, "未知的资源类型"))
		return
	}

	c.JSON(200, response.Success(gin.H{
		"message": "链接上传成功",
		"id":      knowledgePointId,
	}))
}

// 将链接保存到数据库
// 将链接保存到数据库
func saveLinkToDB(knowledgePointId int64, url, title, resourceType string, description *string, difficulty ...string) error {
	switch resourceType {
	case "video":
		// 调用 AddVideo 函数保存视频链接
		err := repository.AddVideo(knowledgePointId, title, url, "", description)
		if err != nil {
			log.Errorf("Failed to save video link to DB: %v", err)
			return err
		}
	case "courseware":
		// 调用 AddCourseware 函数保存课件链接
		err := repository.AddCourseware(knowledgePointId, title, url, description)
		if err != nil {
			log.Errorf("Failed to save courseware link to DB: %v", err)
			return err
		}
	case "exercise":
		// 检查难度是否为空
		if len(difficulty) == 0 || difficulty[0] == "" {
			return fmt.Errorf("难度不能为空")
		}
		// 调用 AddExercise 函数保存练习题链接
		err := repository.AddExercise(knowledgePointId, title, url, difficulty[0], description)
		if err != nil {
			log.Errorf("Failed to save exercise link to DB: %v", err)
			return err
		}
	default:
		return fmt.Errorf("未知的资源类型: %s", resourceType)
	}
	return nil
}

// func UploadResource(c *gin.Context) {
// 	var req UploadResourceRequest
// 	if err := c.ShouldBind(&req); err != nil {
// 		c.JSON(400, response.Error(400, "参数验证失败"))
// 		return
// 	}

// 	data, err := c.FormFile("data")
// 	if err != nil {
// 		c.JSON(400, response.Error(400, "无法读取文件"))
// 		return
// 	}

// 	// 查询知识点ID
// 	knowledgePointId, err := getKnowledgePointIdByName(c, req.KnowledgePointName)
// 	if err != nil {
// 		c.JSON(500, response.Error(500, "查询知识点ID失败"))
// 		return
// 	}

// 	// Extract the filename
// 	filename := filepath.Base(data.Filename)
// 	// Generate a final filename using a random string and the original filename to prevent overwriting
// 	finalName := fmt.Sprintf("%s_%s", utils.RandomString(), filename)

// 	var savePath string
// 	switch req.ResourceType {
// 	case "video":
// 		savePath = config.GetGlobalConfig().MinioConfig.VideoPath
// 	case "courseware", "exercise":
// 		savePath = config.GetGlobalConfig().MinioConfig.FilePath
// 	default:
// 		c.JSON(400, response.Error(400, "未知的资源类型"))
// 		return
// 	}

// 	// 确保保存目录存在
// 	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
// 		log.Errorf("Failed to create directory %s: %v\n", savePath, err)
// 		c.JSON(500, response.Error(500, "创建保存目录失败"))
// 		return
// 	}

// 	saveFile := filepath.Join(savePath, finalName)
// 	if err := c.SaveUploadedFile(data, saveFile); err != nil {
// 		c.JSON(500, response.Error(500, err.Error()))
// 		return
// 	}

// 	if _, err := os.Stat(saveFile); os.IsNotExist(err) {
// 		log.Errorf("File does not exist at path: %s", saveFile)
// 		c.JSON(500, response.Error(500, "文件未正确保存"))
// 		return
// 	}

// 	// 根据资源类型进行具体处理
// 	switch req.ResourceType {
// 	case "video":
// 		err := uploadVideo(knowledgePointId, saveFile, req.Title, req.Description)
// 		if err != nil {
// 			c.JSON(500, response.Error(500, err.Error()))
// 			return
// 		}
// 	case "courseware":
// 		err := uploadFile(knowledgePointId, saveFile, req.Title, req.Description, req.ResourceType)
// 		if err != nil {
// 			c.JSON(500, response.Error(500, err.Error()))
// 			return
// 		}
// 	case "exercise":
// 		if req.Difficulty == nil {
// 			c.JSON(400, response.Error(400, "难度不能为空"))
// 			return
// 		}
// 		err := uploadFile(knowledgePointId, saveFile, req.Title, req.Description, req.ResourceType, *req.Difficulty)
// 		if err != nil {
// 			c.JSON(500, response.Error(500, err.Error()))
// 			return
// 		}
// 	}

//		c.JSON(200, response.Success(gin.H{
//			"message": "资源上传成功",
//			"id":      knowledgePointId,
//		}))
//	}
func uploadVideo(pointId int64, videoPath string, title string, description *string) error {
	// 获取minio客户端并且配置清除函数
	client := minioStore.GetMinio()

	// 临时变量用于存储需要删除的文件路径或URL
	var videoToDelete string
	var localVideoFile string
	var picToDelete string
	var localImageFile string

	// 清理函数：删除 MinIO 和本地（如果有的话）上的文件
	cleanup := func() {
		if videoToDelete != "" {
			_ = client.MinioClient.RemoveObject("video", videoToDelete)
			log.Info("MinIO video file deleted:", videoToDelete)
		}
		if localVideoFile != "" {
			_ = os.Remove(localVideoFile) // 删除本地生成的视频文件
			log.Info("Local video file deleted:", localVideoFile)
		}
		if picToDelete != "" {
			_ = client.MinioClient.RemoveObject("pic", picToDelete)
			log.Info("MinIO pic file deleted:", picToDelete)
		}
		if localImageFile != "" {
			_ = os.Remove(localImageFile) // 删除本地生成的封面图片
			log.Info("Local image file deleted:", localImageFile)
		}
	}

	// 使用 defer 确保即使发生 panic 也会执行清理操作
	defer func() {
		if r := recover(); r != nil {
			cleanup()
			panic(r) // 重新抛出panic
		}
	}()

	localVideoFile = videoPath

	// 上传视频
	videoUrl, videoFilename, err := client.UploadFile("video", videoPath, strconv.FormatInt(pointId, 10))
	if err != nil {
		log.Errorf("UploadFile err: %v", zap.Error(err))
		return err
	}
	videoToDelete = videoFilename

	// 生成视频封面（截取第一帧）
	imageFile, err := utils.GetImageFile(videoPath)
	localImageFile = imageFile
	if err != nil {
		log.Errorf("GetImageFile err: %v", zap.Error(err))
		cleanup() // 发生错误时调用清理函数
		return err
	}
	picUrl, picFilename, err := client.UploadFile("image", imageFile, strconv.FormatInt(pointId, 10))
	if err != nil {
		log.Errorf("UploadFile err: %v", zap.Error(err))
		cleanup() // 发生错误时调用清理函数
		picUrl = "https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/7909abe413ec4a1e82032d2beb810157~tplv-k3u1fbpfcp-zoom-in-crop-mark:1304:0:0:0.awebp?"
	} else {
		picToDelete = picFilename
	}

	// 将视频信息保存到数据库
	err = repository.AddVideo(pointId, title, videoUrl, picUrl, description)
	if err != nil {
		log.Errorf("InsertVideo err: %v", zap.Error(err))
		cleanup() // 发生错误时调用清理函数
		return err
	}

	// 如果所有步骤都成功，则清除清理函数中的记录，防止意外删除
	videoToDelete = ""
	picToDelete = ""
	localImageFile = ""

	return nil
}

func uploadFile(pointId int64, filePath string, title string, description *string, resourceType string, difficulty ...string) error {
	// 获取minio客户端并且配置清除函数
	client := minioStore.GetMinio()

	var fileToDelete string
	var localFilePath string

	// 清理函数：删除 MinIO 和本地（如果有的话）上的文件
	cleanup := func() {
		if fileToDelete != "" {
			err := client.MinioClient.RemoveObject("file", fileToDelete)
			if err != nil {
				log.Warnf("Failed to remove object from MinIO: %s", fileToDelete)
			} else {
				log.Infof("MinIO file deleted: %s", fileToDelete)
			}
		}
		if localFilePath != "" {
			err := os.Remove(localFilePath) // 删除本地生成的文件
			if err != nil {
				log.Warnf("Failed to remove local file: %s", localFilePath)
			} else {
				log.Infof("Local file deleted: %s", localFilePath)
			}
		}
	}

	// 使用 defer 确保即使发生 panic 也会执行清理操作
	defer func() {
		if r := recover(); r != nil {
			cleanup()
			panic(r) // 重新抛出panic
		}
	}()

	// 提取文件名
	filename := filepath.Base(filePath)

	// 根据文件扩展名确定filetype
	filetype := getFileTypeByExtension(filepath.Ext(filename))

	// 上传至MinIO
	url, fileName, err := client.UploadFile(filetype, filePath, strconv.FormatInt(pointId, 10))
	if err != nil {
		log.Errorf("UploadFile err: %v", err)
		cleanup() // 发生错误时调用清理函数
		return err
	}

	// 记录如果发生了任何错误，MinIO上要被删除的文件名
	fileToDelete = fileName
	localFilePath = filePath

	// 根据资源类型保存到数据库
	switch resourceType {
	case "courseware":
		err = repository.AddCourseware(pointId, title, url, description)
	case "exercise":
		diffi := ""
		if len(difficulty) > 0 {
			diffi = difficulty[0]
		}
		err = repository.AddExercise(pointId, title, url, diffi, description)
	default:
		err = fmt.Errorf("未知的资源类型: %s", resourceType)
	}

	if err != nil {
		log.Errorf("AddResource err: %v", err)
		cleanup() // 发生错误时调用清理函数
		return err
	}

	// 如果所有步骤都成功，则清除清理函数中的记录，防止意外删除
	fileToDelete = ""
	localFilePath = ""

	return nil
}

// 根据文件扩展名确定filetype
func getFileTypeByExtension(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".pdf":
		return "pdf"
	case ".pptx":
		return "pptx"
	case ".doc":
		return "word"
	case ".docx":
		return "word"
	default:
		return ""
	}
}

// 假设有一个函数可以从Neo4j中获取知识点ID
func getKnowledgePointIdByName(ctx *gin.Context, name string) (int64, error) {
	nodeLabel := "point" // 指定要查询的节点标签
	var ids []int64      // 用于存储所有节点ID
	session := neo4jUtils.GetSession()
	if session == nil {
		ctx.JSON(500, response.Error(500, "无法获取 Neo4j 会话"))
		return 0, fmt.Errorf("无法获取 Neo4j 会话")
	}
	defer session.Close()

	// 查询带有指定名字的知识点ID
	queryNodes := fmt.Sprintf("MATCH (n:%s) WHERE n.name = $name RETURN n", nodeLabel)
	resultNodes, err := session.Run(queryNodes, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return 0, err
	}
	for resultNodes.Next() {
		record := resultNodes.Record()
		nodeBeforeProc, ok := record.Get("n")
		if !ok {
			return 0, fmt.Errorf("类型断言失败: 无法将记录转换为 Node")
		}
		node := nodeBeforeProc.(dbtype.Node)
		id := int(node.Id)
		ids = append(ids, int64(id)) // 收集所有节点ID
	}
	if len(ids) == 0 {
		return 0, fmt.Errorf("未找到任何有效节点")
	}

	// 返回第一个ID或其他逻辑处理
	return ids[0], nil
}

func GetPointVideo(c *gin.Context) {
	pointId, err := strconv.Atoi(c.Query("pointId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的知识点ID"))
		return
	}

	resources, err := repository.GetVideoByPointId(int64(pointId))
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}

	c.JSON(200, response.Success(resources))
}
func GetPointExercise(c *gin.Context) {
	pointId, err := strconv.Atoi(c.Query("pointId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的知识点ID"))
		return
	}

	resources, err := repository.GetFileByPointId(int64(pointId), "exercise")
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	// log.Infof("exercises: ", resources)

	c.JSON(200, response.Success(resources))
}
func GetPointCourseware(c *gin.Context) {
	pointId, err := strconv.Atoi(c.Query("pointId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的知识点ID"))
		return
	}

	resources, err := repository.GetFileByPointId(int64(pointId), "courseware")
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}

	c.JSON(200, response.Success(resources))
}
func DeletePointVideo(c *gin.Context) {
	videoId, err := strconv.Atoi(c.Query("videoId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的视频ID"))
		return
	}
	pointId, err := strconv.Atoi(c.Query("pointId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的知识点ID"))
		return
	}
	err = repository.DeleteVideo(int64(pointId), int64(videoId))
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success("删除成功"))

}
func DeletePointExercise(c *gin.Context) {
	exerciseId, err := strconv.Atoi(c.Query("exerciseId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的练习题ID"))
		return
	}
	pointId, err := strconv.Atoi(c.Query("pointId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的知识点ID"))
		return
	}
	err = repository.DeleteExercise(int64(pointId), int64(exerciseId))
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success("删除成功"))
}
func DeletePointCourseware(c *gin.Context) {
	coursewareId, err := strconv.Atoi(c.Query("coursewareId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的课件ID"))
		return
	}
	pointId, err := strconv.Atoi(c.Query("pointId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的知识点ID"))
		return
	}
	err = repository.DeleteCourseware(int64(pointId), int64(coursewareId))
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success("删除成功"))
}
