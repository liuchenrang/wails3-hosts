package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chen/wails3-demo/internal/application/service"
	"github.com/chen/wails3-demo/internal/infrastructure/persistence"
	"github.com/chen/wails3-demo/internal/infrastructure/system"
	"github.com/chen/wails3-demo/internal/interface/handler"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 初始化基础设施层
	infra, err := initializeInfrastructure()
	if err != nil {
		log.Fatalf("初始化基础设施失败: %v", err)
	}

	// 初始化应用服务
	appService := service.NewHostsApplicationService(
		infra.hostsRepo,
		infra.versionRepo,
		infra.hostsFileOp,
		infra.sudoManager,
	)

	// 初始化 Wails 处理器
	hostsHandler := handler.NewHostsHandler(appService)

	// 创建 Wails 应用
	app := application.New(application.Options{
		Name:        "Hosts Manager",
		Description: "一个类似 SwitchHosts 的跨平台 hosts 文件管理工具",
		Services: []application.Service{
			application.NewService(hostsHandler),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// 创建主窗口
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Hosts Manager",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		Width:            1200,
		Height:           800,
	})

	// 运行应用
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// infrastructure 基础设施层实例集合
type infrastructure struct {
	hostsRepo   *persistence.HostsRepositoryImpl
	versionRepo *persistence.VersionRepositoryImpl
	hostsFileOp *system.HostsFileOperator
	sudoManager *system.SudoManager
}

// initializeInfrastructure 初始化基础设施层
// 单一职责: 创建和组装所有基础设施组件
// DDD: 依赖注入的根，负责组装整个对象图
func initializeInfrastructure() (*infrastructure, error) {
	// 获取配置目录
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("获取配置目录失败: %w", err)
	}

	// 初始化存储
	storage, err := persistence.NewJSONStorage(configDir)
	if err != nil {
		return nil, fmt.Errorf("创建存储失败: %w", err)
	}

	// 初始化仓储
	hostsRepo := persistence.NewHostsRepository(storage)
	versionRepo := persistence.NewVersionRepository(storage)

	// 初始化系统操作
	hostsFileOp, err := system.NewHostsFileOperator()
	if err != nil {
		return nil, fmt.Errorf("创建 hosts 文件操作器失败: %w", err)
	}

	sudoManager := system.NewSudoManager()

	return &infrastructure{
		hostsRepo:   hostsRepo.(*persistence.HostsRepositoryImpl),
		versionRepo: versionRepo.(*persistence.VersionRepositoryImpl),
		hostsFileOp: hostsFileOp,
		sudoManager: sudoManager,
	}, nil
}

// getConfigDir 获取应用配置目录
// KISS: 根据操作系统返回对应的配置目录
func getConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "hosts-manager"), nil
}
