package database

import (
	"strconv"
	"time"

	"math/rand/v2" // 导入rand/v2包

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 定义表对应的结构体（与SQL表字段对应）
type Vulnerability struct {
	ID            int    `gorm:"primaryKey;autoIncrement"`
	PkgName       string `gorm:"type:varchar(512)"`
	Name          string `gorm:"type:varchar(255)"`
	GroupName     string `gorm:"type:varchar(255)"`
	VersionId     int    `gorm:"type:int"`
	Level         int8   `gorm:"type:tinyint"` // tinyint对应Go的int8
	Reachable     int    `gorm:"type:int"`
	ApplicationId string `gorm:"type:varchar(255)"`
}

// 表名（Gorm默认会复数，这里指定与SQL表名一致）
func (Vulnerability) TableName() string {
	return "vulnerability"
}

func BatchInsert() {
	// 1. 数据库连接配置
	dsn := "root:abc521521521@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 静默日志（批量插入时减少IO开销）
	})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	err = db.AutoMigrate(
		&Vulnerability{},
	)
	if err != nil {
		panic(err)
	}

	// 2. 预生成500种group_name（模拟真实分组，比如组件分组、业务分组）
	var groupNames []string
	for i := 0; i < 500; i++ {
		groupNames = append(groupNames, "group_"+strconv.Itoa(i)+"_"+randString(6)) // 格式：group_0_xxxxxx
	}

	// 3. 降低批次大小（关键修改：从10000调整为500，避免占位符超限；也可根据实际情况调整为1000/2000）
	batchSize := 500 // 核心修改点：减小每批次数据量，规避占位符数量限制
	total := 10000000
	startTime := time.Now()

	// 计算总批次（若total无法被batchSize整除，最后单独处理剩余数据）
	totalBatches := total / batchSize
	remaining := total % batchSize

	for batch := 0; batch < totalBatches; batch++ {
		var dataList []Vulnerability
		for i := 0; i < batchSize; i++ {
			// 随机生成各字段的真实模拟数据（使用rand/v2的IntN方法）
			groupIdx := rand.IntN(500)    // 从500种group_name中随机选
			pkgVer := rand.IntN(1000) + 1 // 版本号（1-1000）
			vulnLevel := rand.IntN(4) + 1 // 风险等级（1-4：低/中/高/严重）
			reachable := rand.IntN(2)     // 是否可达（0/1）

			dataList = append(dataList, Vulnerability{
				PkgName:       "pkg_" + randString(8) + "_v" + strconv.Itoa(pkgVer),
				Name:          "CVE-" + strconv.Itoa(rand.IntN(2026)+2000) + "-" + randString(6),
				GroupName:     groupNames[groupIdx],
				VersionId:     pkgVer,
				Level:         int8(vulnLevel),
				Reachable:     reachable,
				ApplicationId: "app_" + strconv.Itoa(rand.IntN(10000)+1),
			})
		}

		// 批量插入（Gorm的CreateInBatches）
		if err := db.CreateInBatches(dataList, batchSize).Error; err != nil {
			panic("批次" + strconv.Itoa(batch) + "插入失败: " + err.Error())
		}

		// 打印进度
		processed := (batch + 1) * batchSize
		speed := float64(processed) / time.Since(startTime).Seconds()
		println("已插入", processed, "条，速度约", int(speed), "条/秒")
	}

	// 处理剩余数据（避免总数据量遗漏）
	if remaining > 0 {
		var dataList []Vulnerability
		for i := 0; i < remaining; i++ {
			groupIdx := rand.IntN(500)
			pkgVer := rand.IntN(1000) + 1
			vulnLevel := rand.IntN(4) + 1
			reachable := rand.IntN(2)

			dataList = append(dataList, Vulnerability{
				PkgName:       "pkg_" + randString(8) + "_v" + strconv.Itoa(pkgVer),
				Name:          "CVE-" + strconv.Itoa(rand.IntN(2026)+2000) + "-" + randString(6),
				GroupName:     groupNames[groupIdx],
				VersionId:     pkgVer,
				Level:         int8(vulnLevel),
				Reachable:     reachable,
				ApplicationId: "app_" + strconv.Itoa(rand.IntN(10000)+1),
			})
		}
		if err := db.CreateInBatches(dataList, remaining).Error; err != nil {
			panic("剩余数据插入失败: " + err.Error())
		}
		processed := totalBatches*batchSize + remaining
		println("已插入", processed, "条（含剩余数据）")
	}

	// 完成提示
	cost := time.Since(startTime).Minutes()
	println("500w条数据插入完成，耗时约", cost, "分钟")
}

// 生成随机字符串（基于rand/v2实现）
func randString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}
	return string(b)
}
