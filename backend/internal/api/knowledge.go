package api

import (
	"strconv"
	"strings"

	"github.com/RMS_V3/internal/model"
	"github.com/RMS_V3/pkg/response"
	"github.com/gin-gonic/gin"
)

// 模拟数据
var mockData = struct {
	Chapters  []model.Node
	Sections  map[int][]model.Node
	Points    map[int][]model.Node
	Videos    []model.Video
	Exercises []model.Exercise
}{
	Chapters: []model.Node{
		{ID: 1, Name: "数据库基础知识", Type: "chapter"},
		{ID: 2, Name: "关系代数", Type: "chapter"},
		{ID: 3, Name: "SQL语言", Type: "chapter"},
		{ID: 4, Name: "数据库设计", Type: "chapter"},
		{ID: 5, Name: "完整性与安全性", Type: "chapter"},
		{ID: 6, Name: "事务与并发处理", Type: "chapter"},
		{ID: 7, Name: "故障与恢复", Type: "chapter"},
	},
	Sections: map[int][]model.Node{
		1: {
			{ID: 101, Name: "数据库系统概述", Type: "section"},
			{ID: 102, Name: "数据模型", Type: "section"},
			{ID: 103, Name: "数据库系统结构", Type: "section"},
		},
		2: {
			{ID: 201, Name: "关系代数基本运算", Type: "section"},
			{ID: 202, Name: "扩展关系代数", Type: "section"},
			{ID: 203, Name: "关系演算", Type: "section"},
		},
		3: {
			{ID: 301, Name: "SQL基础查询", Type: "section"},
			{ID: 302, Name: "SQL高级查询", Type: "section"},
			{ID: 303, Name: "SQL数据定义", Type: "section"},
			{ID: 304, Name: "SQL数据操作", Type: "section"},
		},
		4: {
			{ID: 401, Name: "ER模型", Type: "section"},
			{ID: 402, Name: "规范化理论", Type: "section"},
			{ID: 403, Name: "数据库设计步骤", Type: "section"},
		},
		5: {
			{ID: 501, Name: "实体完整性", Type: "section"},
			{ID: 502, Name: "参照完整性", Type: "section"},
			{ID: 503, Name: "用户定义的完整性", Type: "section"},
			{ID: 504, Name: "数据库安全性", Type: "section"},
		},
		6: {
			{ID: 601, Name: "事务的基本概念", Type: "section"},
			{ID: 602, Name: "并发控制机制", Type: "section"},
			{ID: 603, Name: "锁与封锁协议", Type: "section"},
			{ID: 604, Name: "多版本并发控制", Type: "section"},
		},
		7: {
			{ID: 701, Name: "故障类型", Type: "section"},
			{ID: 702, Name: "恢复技术", Type: "section"},
			{ID: 703, Name: "备份与恢复策略", Type: "section"},
		},
	},
	Points: map[int][]model.Node{
		101: {
			{ID: 10101, Name: "数据库定义", Type: "knowledge", SectionId: 101, Description: "数据库是按照数据结构来组织、存储和管理数据的建立在计算机存储设备上的仓库。是一个长期存储在计算机内的、有组织的、可共享的、统一管理的数据集合。"},
			{ID: 10102, Name: "DBMS功能", Type: "knowledge", SectionId: 101, Description: "数据库管理系统（DBMS）是一种操纵和管理数据库的大型软件，用于建立、使用和维护数据库，对数据库进行统一管理和控制，以保证数据库的安全性和完整性。"},
			{ID: 10103, Name: "数据库系统体系结构", Type: "knowledge", SectionId: 101, Description: "数据库系统体系结构是从用户角度出发，定义数据库系统的组成结构，它包括数据库、数据库管理系统、应用程序和用户等部分。"},
		},
		102: {
			{ID: 10201, Name: "层次模型", Type: "knowledge", SectionId: 102, Description: "层次模型是数据库早期的数据模型之一，用树形结构表示数据，每个节点可以有多个子节点，但只能有一个父节点，这种限制使得它表示多对多关系较为困难。"},
			{ID: 10202, Name: "网状模型", Type: "knowledge", SectionId: 102, Description: "网状模型是对层次模型的扩展，允许节点有多个父节点，能更好地表示复杂关系，但结构复杂，操作灵活性不足。"},
			{ID: 10203, Name: "关系模型", Type: "knowledge", SectionId: 102, Description: "关系模型使用二维表组织数据，表中每行代表一个实体，每列代表实体的一个属性，通过外键表示实体间关系，结构简单，易于理解与操作。"},
		},
		201: {
			{ID: 20101, Name: "选择与投影", Type: "knowledge", SectionId: 201, Description: "选择是按条件从关系中选出满足条件的元组，投影是从关系中选出指定的属性列，这两种操作是关系代数中最基本的单目运算。"},
			{ID: 20102, Name: "连接操作", Type: "knowledge", SectionId: 201, Description: "连接操作用于将两个关系按照指定的连接条件组合起来，形成新关系，是关系代数中重要的多目运算。"},
			{ID: 20103, Name: "集合操作", Type: "knowledge", SectionId: 201, Description: "关系代数中的集合操作包括并、交、差，用于处理关系之间的集合操作，这些操作要求两个关系有相同的属性结构。"},
		},
		301: {
			{ID: 30101, Name: "SELECT语句基础", Type: "knowledge", SectionId: 301, Description: "SELECT语句是SQL中最常用的语句，用于从一个或多个表、视图中检索数据。基本语法为: SELECT column1, column2... FROM table_name;"},
			{ID: 30102, Name: "WHERE条件查询", Type: "knowledge", SectionId: 301, Description: "WHERE子句用于过滤满足特定条件的记录，可以使用比较运算符和逻辑运算符构建复杂条件。"},
			{ID: 30103, Name: "ORDER BY排序", Type: "knowledge", SectionId: 301, Description: "ORDER BY子句用于对结果集按照一个或多个列进行排序，可指定升序(ASC)或降序(DESC)。"},
		},
		401: {
			{ID: 40101, Name: "实体与联系", Type: "knowledge", SectionId: 401, Description: "实体是客观存在并可相互区别的事物；联系是指实体之间的关系，表现为不同实体集之间的关联。"},
			{ID: 40102, Name: "ER图绘制", Type: "knowledge", SectionId: 401, Description: "ER图是一种表示实体、属性和实体间联系的图形工具，常用矩形表示实体集，椭圆表示属性，菱形表示联系。"},
			{ID: 40103, Name: "从ER图到关系模式", Type: "knowledge", SectionId: 401, Description: "ER图转换为关系模式的基本规则包括：每个实体集转换为一个关系模式，多对多联系转换为一个单独的关系模式等。"},
		},
		501: {
			{ID: 50101, Name: "主键约束", Type: "knowledge", SectionId: 501, Description: "主键约束确保表中每行数据都有唯一标识，主键列不能包含NULL值，通常用于建立索引加速查询。"},
			{ID: 50102, Name: "唯一约束", Type: "knowledge", SectionId: 501, Description: "唯一约束确保列或列组合的值在表中是唯一的，但允许NULL值，可用于确保数据的唯一性而不用作主键。"},
			{ID: 50103, Name: "默认值约束", Type: "knowledge", SectionId: 501, Description: "默认值约束为列指定默认值，当插入数据时没有提供该列的值，系统会自动使用默认值填充。"},
		},
		601: {
			{ID: 60101, Name: "事务的ACID特性", Type: "knowledge", SectionId: 601, Description: "事务的ACID特性包括原子性(Atomicity)、一致性(Consistency)、隔离性(Isolation)和持久性(Durability)。"},
			{ID: 60102, Name: "事务的状态", Type: "knowledge", SectionId: 601, Description: "事务的基本状态包括活动的(Active)、部分提交的(Partially Committed)、失败的(Failed)、中止的(Aborted)和提交的(Committed)。"},
			{ID: 60103, Name: "事务隔离级别", Type: "knowledge", SectionId: 601, Description: "SQL标准定义了四种事务隔离级别：读未提交(Read Uncommitted)、读已提交(Read Committed)、可重复读(Repeatable Read)和可串行化(Serializable)，隔离级别越高，数据一致性越好但并发性能越低。"},
		},
		701: {
			{ID: 70101, Name: "事务故障", Type: "knowledge", SectionId: 701, Description: "事务故障是指事务执行过程中发生错误导致事务无法正常完成，可能由逻辑错误、系统错误或并发冲突引起。"},
			{ID: 70102, Name: "系统故障", Type: "knowledge", SectionId: 701, Description: "系统故障是指导致系统停止运行的事件，如断电、硬件故障等，需要重启系统。"},
			{ID: 70103, Name: "介质故障", Type: "knowledge", SectionId: 701, Description: "介质故障是指存储介质物理损坏导致的数据丢失，通常通过数据备份和冗余存储解决。"},
		},
	},
	Videos: []model.Video{
		{ID: 1, Title: "数据库入门", URL: "https://example.com/video1"},
		{ID: 2, Title: "SQL基础教程", URL: "https://example.com/video2"},
		{ID: 3, Title: "ER模型详解", URL: "https://example.com/video3"},
		{ID: 4, Title: "数据库事务机制", URL: "https://example.com/video4"},
		{ID: 5, Title: "索引优化技巧", URL: "https://example.com/video5"},
	},
	Exercises: []model.Exercise{
		{
			ID:         1,
			Title:      "数据库基础练习",
			Difficulty: "easy",
			URL:        "https://example.com/exercise1",
		},
		{
			ID:         2,
			Title:      "SQL查询综合练习",
			Difficulty: "medium",
			URL:        "https://example.com/exercise2",
		},
		{
			ID:         3,
			Title:      "数据库设计实战",
			Difficulty: "hard",
			URL:        "https://example.com/exercise3",
		},
	},
}

func GetChapterGraph(c *gin.Context) {
	links := []model.Link{
		{Source: 1, Target: 2, Type: "contains"},
		{Source: 2, Target: 3, Type: "prerequisite"},
		{Source: 3, Target: 4, Type: "prerequisite"},
		{Source: 4, Target: 5, Type: "prerequisite"},
		{Source: 5, Target: 6, Type: "prerequisite"},
		{Source: 6, Target: 7, Type: "prerequisite"},
		{Source: 1, Target: 3, Type: "related"},
		{Source: 2, Target: 4, Type: "related"},
		{Source: 4, Target: 6, Type: "related"},
		{Source: 5, Target: 7, Type: "related"},
	}

	c.JSON(200, response.Success(model.Graph{
		Nodes: mockData.Chapters,
		Links: links,
	}))
}

func GetSectionGraph(c *gin.Context) {
	chapterId, err := strconv.Atoi(c.Param("chapterId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的章节ID"))
		return
	}

	sections := mockData.Sections[chapterId]
	if sections == nil {
		c.JSON(404, response.Error(404, "章节不存在"))
		return
	}

	// 生成小节之间的连接
	links := make([]model.Link, 0)
	for i := 0; i < len(sections)-1; i++ {
		links = append(links, model.Link{
			Source: sections[i].ID,
			Target: sections[i+1].ID,
			Type:   "prerequisite",
		})
	}

	// 添加一些额外的关联关系（如果有多个小节）
	if len(sections) > 2 {
		links = append(links, model.Link{
			Source: sections[0].ID,
			Target: sections[len(sections)-1].ID,
			Type:   "related",
		})
	}

	// 查找对应的章节信息
	var chapter model.Node
	for _, ch := range mockData.Chapters {
		if ch.ID == chapterId {
			chapter = ch
			break
		}
	}

	c.JSON(200, response.Success(model.Section{
		Graph: model.Graph{
			Nodes: sections,
			Links: links,
		},
		Chapter: chapter,
	}))
}

func GetPointGraph(c *gin.Context) {
	sectionId, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的小节ID"))
		return
	}

	points := mockData.Points[sectionId]
	if points == nil {
		c.JSON(404, response.Error(404, "小节不存在"))
		return
	}

	// 生成知识点之间的连接
	links := make([]model.Link, 0)
	for i := 0; i < len(points)-1; i++ {
		links = append(links, model.Link{
			Source: points[i].ID,
			Target: points[i+1].ID,
			Type:   "prerequisite",
		})
	}

	// 添加一些额外的依赖关系
	if len(points) > 2 {
		links = append(links, model.Link{
			Source: points[0].ID,
			Target: points[len(points)-1].ID,
			Type:   "extends",
		})
	}

	// 查找对应的小节信息
	chapterId := sectionId / 100
	var section model.Node
	if sections := mockData.Sections[chapterId]; sections != nil {
		for _, s := range sections {
			if s.ID == sectionId {
				section = s
				break
			}
		}
	}

	c.JSON(200, response.Success(model.Point{
		Graph: model.Graph{
			Nodes: points,
			Links: links,
		},
		Section: section,
	}))
}

func GetKnowledgePoint(c *gin.Context) {
	pointId, err := strconv.Atoi(c.Param("pointId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的知识点ID"))
		return
	}

	sectionId := pointId / 100
	points := mockData.Points[sectionId]
	if points == nil {
		c.JSON(404, response.Error(404, "知识点不存在"))
		return
	}

	var point model.Node
	for _, p := range points {
		if p.ID == pointId {
			point = p
			break
		}
	}

	if point.ID == 0 {
		c.JSON(404, response.Error(404, "知识点不存在"))
		return
	}

	c.JSON(200, response.Success(point))
}

func GetKnowledgeVideos(c *gin.Context) {
	pointId, err := strconv.Atoi(c.Param("pointId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的知识点ID"))
		return
	}

	// 根据知识点ID筛选相关视频
	var videos []model.Video

	// 为特定知识点提供对应视频
	switch pointId / 10000 {
	case 1: // 数据库基础相关知识点
		videos = append(videos, mockData.Videos[0])
	case 2: // SQL相关知识点
		videos = append(videos, mockData.Videos[1])
	case 3: // 数据库设计相关知识点
		videos = append(videos, mockData.Videos[2])
	case 4: // 事务相关知识点
		videos = append(videos, mockData.Videos[3])
	case 5: // 索引优化相关知识点
		videos = append(videos, mockData.Videos[4])
	default:
		// 默认返回所有视频
		videos = mockData.Videos
	}

	c.JSON(200, response.Success(videos))
}

func GetKnowledgeExercises(c *gin.Context) {
	pointId, err := strconv.Atoi(c.Param("pointId"))
	if err != nil {
		c.JSON(400, response.Error(400, "无效的知识点ID"))
		return
	}

	// 根据知识点ID筛选相关练习题
	var exercises []model.Exercise

	// 为特定知识点提供对应练习题
	switch pointId / 10000 {
	case 1: // 数据库基础相关知识点
		exercises = append(exercises, mockData.Exercises[0])
	case 2: // SQL相关知识点
		exercises = append(exercises, mockData.Exercises[1])
	case 3: // 数据库设计相关知识点
		exercises = append(exercises, mockData.Exercises[2])
	default:
		// 默认返回所有练习题
		exercises = mockData.Exercises
	}

	c.JSON(200, response.Success(exercises))
}

// 搜索知识点
func SearchKnowledge(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(400, response.Error(400, "搜索关键词不能为空"))
		return
	}

	// 搜索结果
	var results []map[string]interface{}

	// 在章节中搜索
	for _, chapter := range mockData.Chapters {
		if strings.Contains(strings.ToLower(chapter.Name), strings.ToLower(keyword)) {
			results = append(results, map[string]interface{}{
				"id":    chapter.ID,
				"name":  chapter.Name,
				"type":  chapter.Type,
				"level": "chapter",
			})
		}
	}

	// 在小节中搜索
	for _, sections := range mockData.Sections {
		for _, section := range sections {
			if strings.Contains(strings.ToLower(section.Name), strings.ToLower(keyword)) {
				results = append(results, map[string]interface{}{
					"id":    section.ID,
					"name":  section.Name,
					"type":  section.Type,
					"level": "section",
				})
			}
		}
	}

	// 在知识点中搜索
	for _, points := range mockData.Points {
		for _, point := range points {
			if strings.Contains(strings.ToLower(point.Name), strings.ToLower(keyword)) {
				results = append(results, map[string]interface{}{
					"id":    point.ID,
					"name":  point.Name,
					"type":  point.Type,
					"level": "point",
				})
			}
		}
	}

	c.JSON(200, response.Success(results))
}

// ... 其他API处理函数实现
