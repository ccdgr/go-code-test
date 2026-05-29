package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 一、定义心理测评平台核心结构体（新增学生表，完善关联关系）
// School 学校表：存储测评平台所属学校基础信息
type School struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolNo  string         `gorm:"type:varchar(30);not null;unique" json:"school_no"` // 学校编号，唯一
	Name      string         `gorm:"type:varchar(100);not null;unique" json:"name"`     // 学校名称，唯一
	Address   string         `gorm:"type:varchar(255);not null" json:"address"`         // 学校地址
	Phone     string         `gorm:"type:varchar(20);not null" json:"phone"`            // 联系电话
	Level     string         `gorm:"type:varchar(20);not null" json:"level"`            // 学校等级（小学/初中/高中/大学）
	CreatedAt time.Time      `json:"created_at"`                                        // GORM自动维护
	UpdatedAt time.Time      `json:"updated_at"`                                        // GORM自动维护
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`                           // 软删除字段
}

// Student 学生表（新增）：存储学生基础信息，关联学校
type Student struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentNo   string         `gorm:"type:varchar(30);not null;unique" json:"student_no"` // 学生学号，唯一
	StudentName string         `gorm:"type:varchar(50);not null" json:"student_name"`      // 学生姓名
	Gender      int            `gorm:"type:tinyint;not null;default:1" json:"gender"`      // 性别：1-男，2-女
	Grade       string         `gorm:"type:varchar(20);not null" json:"grade"`             // 年级（如：三年级2班）
	Phone       string         `gorm:"type:varchar(20);nullable" json:"phone"`             // 家长联系电话
	SchoolID    uint           `gorm:"not null;index" json:"school_id"`                    // 关联学校ID
	CreatedAt   time.Time      `json:"created_at"`                                         // GORM自动维护
	UpdatedAt   time.Time      `json:"updated_at"`                                         // GORM自动维护
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`                            // 软删除字段
}

// EvaluationPaper 测评试卷表：存储心理测评试卷基础信息
type EvaluationPaper struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	PaperNo     string         `gorm:"type:varchar(30);not null;unique" json:"paper_no"` // 试卷编号，唯一
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`          // 试卷标题（如：中学生心理健康测评试卷）
	Description string         `gorm:"type:text;nullable" json:"description"`            // 试卷描述
	QuestionNum int            `gorm:"type:int;not null;default:0" json:"question_num"`  // 题目数量
	Duration    int            `gorm:"type:int;not null;default:60" json:"duration"`     // 考试时长（分钟）
	Status      int            `gorm:"type:tinyint;not null;default:1" json:"status"`    // 状态：1-可用，2-禁用
	CreatedAt   time.Time      `json:"created_at"`                                       // GORM自动维护
	UpdatedAt   time.Time      `json:"updated_at"`                                       // GORM自动维护
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`                          // 软删除字段
}

// Exam 考试表：存储心理测评考试信息（如：2025年秋季学期心理健康测评）
type Exam struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ExamNo    string         `gorm:"type:varchar(30);not null;unique" json:"exam_no"` // 考试编号，唯一
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`         // 考试标题
	SchoolID  uint           `gorm:"not null;index" json:"school_id"`                 // 关联学校ID（哪个学校的考试）
	StartTime time.Time      `json:"start_time"`                                      // 考试开始时间
	EndTime   time.Time      `json:"end_time"`                                        // 考试结束时间
	Status    int            `gorm:"type:tinyint;not null;default:1" json:"status"`   // 状态：1-未开始，2-进行中，3-已结束
	CreatedAt time.Time      `json:"created_at"`                                      // GORM自动维护
	UpdatedAt time.Time      `json:"updated_at"`                                      // GORM自动维护
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`                         // 软删除字段
}

// ExamPaperRelation 考试关联试卷表（1对多：一个考试对应多份试卷）
type ExamPaperRelation struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ExamID    uint           `gorm:"not null;index" json:"exam_id"`           // 关联考试ID
	PaperID   uint           `gorm:"not null;index" json:"paper_id"`          // 关联试卷ID
	Sort      int            `gorm:"type:int;not null;default:1" json:"sort"` // 试卷排序（一个考试多份试卷的展示顺序）
	CreatedAt time.Time      `json:"created_at"`                              // GORM自动维护
	UpdatedAt time.Time      `json:"updated_at"`                              // GORM自动维护
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`                 // 软删除字段
	// 联合唯一索引：避免同一考试重复关联同一试卷
	_ struct{} `gorm:"uniqueIndex:uk_exam_paper(exam_id,paper_id)"`
}

// PaperQuestionRelation 测评试卷题目关联表（1对多：一份试卷对应多道题目）
type PaperQuestionRelation struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	PaperID    uint           `gorm:"not null;index" json:"paper_id"`               // 关联试卷ID
	QuestionID string         `gorm:"type:varchar(30);not null" json:"question_id"` // 题目ID（可关联题目表，此处简化直接存储ID）
	Question   string         `gorm:"type:text;not null" json:"question"`           // 题目内容（简化设计，无需单独建题目表）
	Score      int            `gorm:"type:int;not null;default:2" json:"score"`     // 题目分值
	Sort       int            `gorm:"type:int;not null;default:1" json:"sort"`      // 题目排序
	CreatedAt  time.Time      `json:"created_at"`                                   // GORM自动维护
	UpdatedAt  time.Time      `json:"updated_at"`                                   // GORM自动维护
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`                      // 软删除字段
	// 联合唯一索引：避免同一试卷重复关联同一题目
	_ struct{} `gorm:"uniqueIndex:uk_paper_question(paper_id,question_id)"`
}

// StudentScore 学生成绩表（修改：关联学生ID，移除冗余的学号/姓名字段）
type StudentScore struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID  uint           `gorm:"not null;index" json:"student_id"`               // 关联学生ID（核心修改：关联独立学生表）
	SchoolID   uint           `gorm:"not null;index" json:"school_id"`                // 关联学校ID
	ExamID     uint           `gorm:"not null;index" json:"exam_id"`                  // 关联考试ID
	PaperID    uint           `gorm:"not null;index" json:"paper_id"`                 // 关联试卷ID
	TotalScore int            `gorm:"type:int;not null;default:0" json:"total_score"` // 试卷总分
	SubmitTime time.Time      `json:"submit_time"`                                    // 提交时间
	Remark     string         `gorm:"type:varchar(255);nullable" json:"remark"`       // 备注（如：缺考/作弊）
	CreatedAt  time.Time      `json:"created_at"`                                     // GORM自动维护
	UpdatedAt  time.Time      `json:"updated_at"`                                     // GORM自动维护
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`                        // 软删除字段
	// 联合索引：优化学生+考试+试卷的查询效率，确保数据唯一性
	_ struct{} `gorm:"index:idx_student_exam_paper(student_id,exam_id,paper_id)"`
}

// 全局数据库连接实例
var db *gorm.DB

// 初始化数据库连接
func initDB() error {
	// MySQL DSN格式：user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	// 请替换为你的数据库账号、密码、数据库名（需先手动创建school_psychology数据库）
	dsn := "root:abc521521521@tcp(127.0.0.1:3306)/psy?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	// 连接数据库，开启SQL日志（开发环境便于调试）
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("数据库连接失败：%v", err)
	}

	// 自动迁移表结构（包含新增的Student表）
	err = db.AutoMigrate(
		&School{},
		&Student{}, // 新增学生表迁移
		&EvaluationPaper{},
		&Exam{},
		&ExamPaperRelation{},
		&PaperQuestionRelation{},
		&StudentScore{},
	)
	if err != nil {
		return fmt.Errorf("表结构迁移失败：%v", err)
	}

	fmt.Println("数据库初始化成功，表结构迁移完成")
	return nil
}

// -------------- 模拟数据生成函数（新增学生表数据生成，调整成绩表关联）--------------
// 生成模拟学校数据
func generateMockSchools(count int) []School {
	schools := make([]School, 0, count)
	schoolNames := []string{"阳光小学", "育才中学", "实验外国语学校", "蓝天中学", "明德高中", "理工大学附属中学", "师范大学附属小学", "光华中学", "星辰小学", "希望高中"}
	schoolLevels := []string{"小学", "初中", "高中", "大学"}
	addresses := []string{"北京市海淀区中关村大街", "上海市浦东新区陆家嘴环路", "广州市天河区天河路", "深圳市南山区科技园路", "杭州市西湖区西湖大道", "成都市锦江区春熙路", "重庆市渝中区解放碑", "南京市秦淮区夫子庙路", "武汉市武昌区中南路", "西安市雁塔区长安路"}
	phones := []string{"010-12345678", "021-87654321", "020-11223344", "0755-55667788", "0571-99887766", "028-66778899", "023-88990011", "025-77889900", "027-66554433", "029-55443322"}

	for i := 0; i < count; i++ {
		school := School{
			SchoolNo: fmt.Sprintf("SCH%06d", i+1),
			Name:     fmt.Sprintf("%s-%d", schoolNames[i%len(schoolNames)], i/len(schoolNames)+1),
			Address:  fmt.Sprintf("%s%d号", addresses[rand.IntN(len(addresses))], rand.IntN(100)),
			Phone:    phones[rand.IntN(len(phones))],
			Level:    schoolLevels[rand.IntN(len(schoolLevels))],
		}
		schools = append(schools, school)
	}
	return schools
}

// 生成模拟学生数据（新增：每所学校指定数量学生）
func generateMockStudents(schoolIDs []uint, studentsPerSchool int) []Student {
	totalCount := len(schoolIDs) * studentsPerSchool
	students := make([]Student, 0, totalCount)
	studentNames := []string{"张三", "李四", "王五", "赵六", "钱七", "孙八", "周九", "吴十", "郑十一", "王十二", "李十三", "张十四", "刘十五", "陈十六", "杨十七"}
	grades := []string{"一年级1班", "一年级2班", "二年级1班", "二年级2班", "三年级1班", "三年级2班", "四年级1班", "四年级2班", "五年级1班", "五年级2班", "六年级1班", "六年级2班", "初一1班", "初一2班", "初二1班", "初二2班", "初三1班", "初三2班", "高一1班", "高一2班", "高二1班", "高二2班", "高三1班", "高三2班"}
	phonePrefixes := []string{"138", "139", "137", "136", "135", "186", "189", "177"}

	// 为每所学校生成指定数量学生
	for _, schoolID := range schoolIDs {
		for i := 0; i < studentsPerSchool; i++ {
			studentNo := fmt.Sprintf("STU_%06d_%04d", schoolID, i+1) // 学号格式：学校ID_学生序号
			studentName := fmt.Sprintf("%s%d", studentNames[rand.IntN(len(studentNames))], rand.IntN(100))
			gender := rand.IntN(2) + 1 // 1-男，2-女
			grade := grades[rand.IntN(len(grades))]
			phone := fmt.Sprintf("%s%08d", phonePrefixes[rand.IntN(len(phonePrefixes))], rand.IntN(100000000))

			student := Student{
				StudentNo:   studentNo,
				StudentName: studentName,
				Gender:      gender,
				Grade:       grade,
				Phone:       phone,
				SchoolID:    schoolID,
			}
			students = append(students, student)
		}
	}
	return students
}

// 生成模拟测评试卷数据
func generateMockPapers(count int) []EvaluationPaper {
	papers := make([]EvaluationPaper, 0, count)
	paperTitles := []string{
		"小学生心理健康测评试卷",
		"中学生焦虑情绪测评试卷",
		"高中生压力应对能力测评试卷",
		"大学生人际交往能力测评试卷",
		"青少年抑郁倾向筛查试卷",
		"师生关系满意度测评试卷",
		"小学生情绪管理能力测评试卷",
		"中学生自我认知测评试卷",
		"高中生升学压力测评试卷",
		"青少年人际交往障碍筛查试卷",
	}

	for i := 0; i < count; i++ {
		questionNum := rand.IntN(40) + 10 // 题目数量：10-50题
		duration := rand.IntN(60) + 30    // 考试时长：30-90分钟
		status := rand.IntN(2) + 1        // 状态：1-可用，2-禁用

		paper := EvaluationPaper{
			PaperNo:     fmt.Sprintf("PAP%06d", i+1),
			Title:       fmt.Sprintf("%s-%d", paperTitles[i%len(paperTitles)], i/len(paperTitles)+1),
			Description: fmt.Sprintf("本试卷共%d题，考试时长%d分钟，用于测评学生心理健康状态", questionNum, duration),
			QuestionNum: questionNum,
			Duration:    duration,
			Status:      status,
		}
		papers = append(papers, paper)
	}
	return papers
}

// 生成模拟考试数据
func generateMockExams(count int, schoolIDs []uint) []Exam {
	exams := make([]Exam, 0, count)
	examTitles := []string{
		"春季学期心理健康测评",
		"秋季学期心理健康测评",
		"开学初心理健康筛查",
		"期末心理健康总结测评",
		"考前心理状态测评",
		"寒假前心理健康测评",
		"暑假前心理健康测评",
	}

	for i := 0; i < count; i++ {
		// 随机选择学校ID
		schoolID := schoolIDs[rand.IntN(len(schoolIDs))]
		// 生成开始时间和结束时间（结束时间比开始时间晚2小时）
		startTime := time.Now().Add(-time.Duration(rand.IntN(30)) * 24 * time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		status := rand.IntN(3) + 1 // 状态：1-未开始，2-进行中，3-已结束

		exam := Exam{
			ExamNo:    fmt.Sprintf("EXM%06d", i+1),
			Title:     fmt.Sprintf("%s-%d", examTitles[i%len(examTitles)], i/len(examTitles)+1),
			SchoolID:  schoolID,
			StartTime: startTime,
			EndTime:   endTime,
			Status:    status,
		}
		exams = append(exams, exam)
	}
	return exams
}

// 生成考试-试卷关联数据
func generateMockExamPaperRelations(count int, examIDs []uint, paperIDs []uint) []ExamPaperRelation {
	relations := make([]ExamPaperRelation, 0, count)

	for i := 0; i < count; i++ {
		examID := examIDs[rand.IntN(len(examIDs))]
		paperID := paperIDs[rand.IntN(len(paperIDs))]
		sort := rand.IntN(10) + 1 // 排序值：1-10

		relation := ExamPaperRelation{
			ExamID:  examID,
			PaperID: paperID,
			Sort:    sort,
		}
		relations = append(relations, relation)
	}
	return relations
}

// 生成试卷-题目关联数据
func generateMockPaperQuestionRelations(count int, paperIDs []uint) []PaperQuestionRelation {
	relations := make([]PaperQuestionRelation, 0, count)
	questionContents := []string{
		"你最近一周的睡眠质量如何？",
		"你是否经常感到学习压力过大？",
		"你与同学的人际关系是否融洽？",
		"你是否对未来有明确的规划？",
		"你是否经常感到情绪低落？",
		"你是否能快速适应新的学习环境？",
		"你与父母的沟通是否顺畅？",
		"你是否有自己的兴趣爱好？",
		"你是否经常感到烦躁易怒？",
		"你对自己的学习成绩是否满意？",
	}

	for i := 0; i < count; i++ {
		paperID := paperIDs[rand.IntN(len(paperIDs))]
		questionID := fmt.Sprintf("QUE%06d", i+1)
		question := questionContents[rand.IntN(len(questionContents))]
		score := rand.IntN(5) + 1 // 分值：1-5分
		sort := rand.IntN(50) + 1 // 排序值：1-50

		relation := PaperQuestionRelation{
			PaperID:    paperID,
			QuestionID: questionID,
			Question:   question,
			Score:      score,
			Sort:       sort,
		}
		relations = append(relations, relation)
	}
	return relations
}

// 生成学生成绩数据（修改：关联学生ID，支持50万+数据）
func generateMockStudentScores(totalCount int, studentIDs []uint, schoolIDs []uint, examIDs []uint, paperIDs []uint) []StudentScore {
	scores := make([]StudentScore, 0, totalCount)

	for i := 0; i < totalCount; i++ {
		studentID := studentIDs[rand.IntN(len(studentIDs))] // 随机关联学生ID
		schoolID := schoolIDs[rand.IntN(len(schoolIDs))]    // 随机关联学校ID
		examID := examIDs[rand.IntN(len(examIDs))]          // 随机关联考试ID
		paperID := paperIDs[rand.IntN(len(paperIDs))]       // 随机关联试卷ID
		totalScore := rand.IntN(100)                        // 总分：0-100分
		submitTime := time.Now().Add(-time.Duration(rand.IntN(20)) * 24 * time.Hour)
		// 随机生成备注（大部分为空）
		remark := ""
		if rand.IntN(100) < 5 { // 5%的概率有备注
			remarks := []string{"缺考", "作弊", "超时提交", "漏答题"}
			remark = remarks[rand.IntN(len(remarks))]
		}

		score := StudentScore{
			StudentID:  studentID,
			SchoolID:   schoolID,
			ExamID:     examID,
			PaperID:    paperID,
			TotalScore: totalScore,
			SubmitTime: submitTime,
			Remark:     remark,
		}
		scores = append(scores, score)
	}
	return scores
}

// 分批插入数据（核心优化：避免单次插入过大导致内存溢出/超时）
// batchSize：每批次插入数量（建议1000-5000，根据数据库性能调整）
func batchInsert[T any](data []T, batchSize int) error {
	if len(data) == 0 {
		return nil
	}

	start := 0
	total := len(data)
	for start < total {
		end := min(start+batchSize, total)
		// 直接传入子切片，无需取地址
		if err := db.Create(data[start:end]).Error; err != nil {
			return fmt.Errorf("批次插入失败（%d-%d）：%v", start, end, err)
		}
		fmt.Printf("已插入%d条数据，累计插入%d条\n", end-start, end)
		start = end
		// 短暂休眠，避免数据库压力过大
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func batch() {
	// 1. 初始化数据库
	err := initDB()
	if err != nil {
		fmt.Printf("初始化失败：%v\n", err)
		return
	}

	// 2. 配置数据量
	schoolCount := 10                  // 10所学校
	studentsPerSchool := 5000          // 每校5000学生，总计5万学生
	paperCount := 20                   // 20份试卷
	examCount := 15                    // 15场考试
	examPaperRelCount := 30            // 30条考试-试卷关联
	paperQuestionRelCount := 500       // 500条试卷-题目关联
	studentScoreTotalCount := 10000000 // 50万条学生成绩
	batchSize := 2000                  // 每批次插入2000条

	// 2.1 生成并插入学校数据
	schools := generateMockSchools(schoolCount)
	err = batchInsert(schools, batchSize)
	if err != nil {
		fmt.Printf("学校数据插入失败：%v\n", err)
		return
	}
	// 查询获取自增后的学校ID
	var schoolIDs []uint
	schools = []School{}
	db.Find(&schools)
	for _, s := range schools {
		schoolIDs = append(schoolIDs, s.ID)
	}
	fmt.Printf("学校数据插入完成，共%d所学校，获取到学校ID：%v\n", len(schoolIDs), schoolIDs)

	// 2.2 生成并插入学生数据（每校5000人）
	students := generateMockStudents(schoolIDs, studentsPerSchool)
	err = batchInsert(students, batchSize)
	if err != nil {
		fmt.Printf("学生数据插入失败：%v\n", err)
		return
	}
	// 查询获取自增后的学生ID
	var studentIDs []uint
	students = []Student{}
	db.Find(&students)
	for _, stu := range students {
		studentIDs = append(studentIDs, stu.ID)
	}
	fmt.Printf("学生数据插入完成，共%d名学生，获取到学生ID总数：%d\n", len(studentIDs), len(studentIDs))

	// 2.3 生成并插入试卷数据
	papers := generateMockPapers(paperCount)
	err = batchInsert(papers, batchSize)
	if err != nil {
		fmt.Printf("试卷数据插入失败：%v\n", err)
		return
	}
	// 查询获取自增后的试卷ID
	var paperIDs []uint
	papers = []EvaluationPaper{}
	db.Find(&papers)
	for _, p := range papers {
		paperIDs = append(paperIDs, p.ID)
	}
	fmt.Printf("试卷数据插入完成，共%d份试卷\n", len(paperIDs))

	// 2.4 生成并插入考试数据
	exams := generateMockExams(examCount, schoolIDs)
	err = batchInsert(exams, batchSize)
	if err != nil {
		fmt.Printf("考试数据插入失败：%v\n", err)
		return
	}
	// 查询获取自增后的考试ID
	var examIDs []uint
	exams = []Exam{}
	db.Find(&exams)
	for _, e := range exams {
		examIDs = append(examIDs, e.ID)
	}
	fmt.Printf("考试数据插入完成，共%d场考试\n", len(examIDs))

	// 2.5 生成并插入考试-试卷关联数据
	examPaperRels := generateMockExamPaperRelations(examPaperRelCount, examIDs, paperIDs)
	err = batchInsert(examPaperRels, batchSize)
	if err != nil {
		fmt.Printf("考试-试卷关联数据插入失败：%v\n", err)
		return
	}
	fmt.Printf("考试-试卷关联数据插入完成，共%d条关联\n", len(examPaperRels))

	// 2.6 生成并插入试卷-题目关联数据
	paperQuestionRels := generateMockPaperQuestionRelations(paperQuestionRelCount, paperIDs)
	err = batchInsert(paperQuestionRels, batchSize)
	if err != nil {
		fmt.Printf("试卷-题目关联数据插入失败：%v\n", err)
		return
	}
	fmt.Printf("试卷-题目关联数据插入完成，共%d条关联\n", len(paperQuestionRels))

	// 2.7 生成并插入50万条学生成绩数据（关联学生ID）
	fmt.Println("开始插入50万条学生成绩数据，请耐心等待...")
	studentScores := generateMockStudentScores(studentScoreTotalCount, studentIDs, schoolIDs, examIDs, paperIDs)
	err = batchInsert(studentScores, batchSize)
	if err != nil {
		fmt.Printf("学生成绩数据插入失败：%v\n", err)
		return
	}

	// 输出统计信息
	fmt.Printf("\n所有数据插入完成！统计如下：\n")
	fmt.Printf("1. 学校：%d所\n", schoolCount)
	fmt.Printf("2. 学生：%d名（每校%d名）\n", len(studentIDs), studentsPerSchool)
	fmt.Printf("3. 试卷：%d份\n", paperCount)
	fmt.Printf("4. 考试：%d场\n", examCount)
	fmt.Printf("5. 考试-试卷关联：%d条\n", examPaperRelCount)
	fmt.Printf("6. 试卷-题目关联：%d条\n", paperQuestionRelCount)
	fmt.Printf("7. 学生成绩：%d条\n", studentScoreTotalCount)
}
