package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type BinlogSyncer struct {
	syncer *replication.BinlogSyncer
}

func NewBinlogSyncer() *BinlogSyncer {
	// 创建同步器配置
	cfg := replication.BinlogSyncerConfig{
		ServerID: 100, // 必须唯一，不能与MySQL其他slave重复
		Flavor:   "mysql",
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "abc521521521",
	}

	syncer := replication.NewBinlogSyncer(cfg)
	return &BinlogSyncer{syncer: syncer}
}

func (bs *BinlogSyncer) Start() error {
	// 从最新的binlog位置开始
	streamer, err := bs.syncer.StartSync(mysql.Position{})
	if err != nil {
		return err
	}

	// 处理信号量，优雅关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	fmt.Println("开始监听binlog...")

	// 创建带超时的context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case <-c:
			fmt.Println("收到停止信号，退出程序")
			bs.syncer.Close()
			return nil
		default:
			// 添加 context 参数和超时控制
			ev, err := streamer.GetEvent(ctx)
			if err != nil {
				fmt.Printf("获取事件错误: %v\n", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			bs.handleEvent(ev)
		}
	}
}

func (bs *BinlogSyncer) handleEvent(ev *replication.BinlogEvent) {
	switch ev.Header.EventType {
	case replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
		rowsEvent := ev.Event.(*replication.RowsEvent)
		fmt.Printf("插入数据 - 表: %s.%s, 行数: %d\n",
			string(rowsEvent.Table.Schema),
			string(rowsEvent.Table.Table),
			len(rowsEvent.Rows))

		// 处理插入的数据
		for _, row := range rowsEvent.Rows {
			fmt.Printf("插入的行数据: %v\n", row)
		}

	case replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
		rowsEvent := ev.Event.(*replication.RowsEvent)
		fmt.Printf("更新数据 - 表: %s.%s, 行数: %d\n",
			string(rowsEvent.Table.Schema),
			string(rowsEvent.Table.Table),
			len(rowsEvent.Rows)/2) // update事件包含旧值和新值

		// 处理更新的数据（每两行为一组：旧值、新值）
		for i := 0; i < len(rowsEvent.Rows); i += 2 {
			if i+1 < len(rowsEvent.Rows) {
				fmt.Printf("旧值: %v -> 新值: %v\n",
					rowsEvent.Rows[i], rowsEvent.Rows[i+1])
			}
		}

	case replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
		rowsEvent := ev.Event.(*replication.RowsEvent)
		fmt.Printf("删除数据 - 表: %s.%s, 行数: %d\n",
			string(rowsEvent.Table.Schema),
			string(rowsEvent.Table.Table),
			len(rowsEvent.Rows))

		// 处理删除的数据
		for _, row := range rowsEvent.Rows {
			fmt.Printf("删除的行数据: %v\n", row)
		}

	case replication.FORMAT_DESCRIPTION_EVENT:
		fmt.Println("Binlog格式描述事件")

	case replication.TABLE_MAP_EVENT:
		tableMapEvent := ev.Event.(*replication.TableMapEvent)
		fmt.Printf("表映射事件 - 数据库: %s, 表: %s\n",
			string(tableMapEvent.Schema),
			string(tableMapEvent.Table))

	case replication.QUERY_EVENT:
		queryEvent := ev.Event.(*replication.QueryEvent)
		fmt.Printf("SQL查询: %s\n", string(queryEvent.Query))
	}
}

// func main() {
// 	syncer := NewBinlogSyncer()
// 	if err := syncer.Start(); err != nil {
// 		fmt.Printf("启动失败: %v\n", err)
// 	}
// }
