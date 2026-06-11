// Copyright 2026 ChromaDesk Authors
// ChromaDesk Qt Application Entry Point

#include <QApplication>
#include <QFontDatabase>
#include <QIcon>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QSharedMemory>
#include <QJSEngine>
#include <QtQuickControls2/QQuickStyle>

#include "core/localconfigcenter.h"
#include "core/userdatacenter.h"
#include "infra/env/applicationcontext.h"
#include "infra/log/log.h"

#include "api/OcrEngine.h"
#include "controller/MainController.h"
#include "manager/ServerManager.h"
#include "manager/HostManager.h"
#include "manager/ClientManager.h"
#include "manager/SharedMemoryManager.h"
#include "component/VideoFrameProvider.h"
#include "component/KeycodeMapper.h"
#include "component/CursorImageProvider.h"
#include "component/SystemTrayManager.h"
#include "viewmodel/configviewmodel.h"
#include "viewmodel/connectionlistmodel.h"
#include "language/languagemanage.h"
#include "common/ProcessStatus.h"

int main(int argc, char *argv[])
{
    QSharedMemory sharedMemory("quick_coder_qdsm");
    // attach成功说明已经create过了，直接退出
    if (sharedMemory.attach()) {
        return 0;
    }
    sharedMemory.create(1);

    QApplication::setHighDpiScaleFactorRoundingPolicy(
        Qt::HighDpiScaleFactorRoundingPolicy::PassThrough);

    QApplication app(argc, argv);
    
    // 设置应用图标
    app.setWindowIcon(QIcon(":/ChromaDesk.ico"));

    // infra
    infra::ApplicationContext::instance().init();
    infra::Log::instance().init(infra::ApplicationContext::instance().logPath());
    LOG_INFO("start app {}, log level:{} ********", infra::ApplicationContext::instance().applicationName().toStdString(), SPDLOG_ACTIVE_LEVEL);
    infra::ApplicationContext::instance().setApplicationVersion(APP_VERSION_STR);
    LOG_INFO("current version:{}", infra::ApplicationContext::instance().applicationVersion().toStdString());

    // QML FileDialog需要这个
    // 放在ApplicationContext::init后面设置，否则OrganizationName会包含在QStandardPaths::AppLocalDataLocation目录中
    app.setOrganizationName("QuickCoder");
    app.setApplicationName("ChromaDesk");

    // 设置 Qt Quick Controls 使用 Basic 样式（允许完全自定义）
    QQuickStyle::setStyle("Basic");
    
    LOG_INFO("ChromaDesk starting...");
    LOG_INFO("Qt version: {}", qVersion());

    core::LocalConfigCenter::instance().init();
    core::UserDataCenter::instance().init();
    LanguageManage::instance().init();

    // Register C++ types for QML
    qmlRegisterType<ConfigViewModel>("ChromaDesk", 1, 0, "ConfigViewModel");
    qmlRegisterType<chromadesk::MainController>("ChromaDesk", 1, 0, "MainController");
    qmlRegisterUncreatableType<chromadesk::ServerManager>("ChromaDesk", 1, 0, "ServerManager",
        "ServerManager is accessed through MainController");
    qmlRegisterUncreatableType<chromadesk::HostManager>("ChromaDesk", 1, 0, "HostManager",
        "HostManager is accessed through MainController");
    qmlRegisterUncreatableType<chromadesk::ClientManager>("ChromaDesk", 1, 0, "ClientManager",
        "ClientManager is accessed through MainController");
    qmlRegisterUncreatableType<chromadesk::SharedMemoryManager>("ChromaDesk", 1, 0, "SharedMemoryManager",
        "SharedMemoryManager is accessed through ClientManager");
    qmlRegisterType<chromadesk::VideoFrameProvider>("ChromaDesk", 1, 0, "VideoFrameProvider");
    qmlRegisterType<chromadesk::ConnectionListModel>("ChromaDesk", 1, 0, "ConnectionListModel");
    
    // Register enums for QML
    qmlRegisterUncreatableType<chromadesk::ProcessStatus>("ChromaDesk", 1, 0, "ProcessStatus",
        "ProcessStatus is an enum container");
    qmlRegisterUncreatableType<chromadesk::ServerStatus>("ChromaDesk", 1, 0, "ServerStatus",
        "ServerStatus is an enum container");
    qmlRegisterUncreatableType<chromadesk::RtcStatus>("ChromaDesk", 1, 0, "RtcStatus",
        "RtcStatus is an enum container");
    qmlRegisterUncreatableType<chromadesk::HostLaunchMode>("ChromaDesk", 1, 0, "HostLaunchMode",
        "HostLaunchMode is an enum container");
    
    // Register KeyboardStateTracker as singleton
    qmlRegisterSingletonType<chromadesk::KeyboardStateTracker>("ChromaDesk", 1, 0, "KeyboardStateTracker",
        [](QQmlEngine*, QJSEngine*) -> QObject* {
            return &chromadesk::KeyboardStateTracker::instance();
        });
    
    // Register LanguageManage as singleton
    qmlRegisterSingletonType<LanguageManage>("ChromaDesk", 1, 0, "LanguageManage",
        [](QQmlEngine*, QJSEngine*) -> QObject* {
            return &LanguageManage::instance();
        });

    // Register SystemTrayManager as singleton
    qmlRegisterSingletonType<chromadesk::SystemTrayManager>("ChromaDesk", 1, 0, "SystemTrayManager",
        [](QQmlEngine*, QJSEngine*) -> QObject* {
            auto* mgr = &chromadesk::SystemTrayManager::instance();
            QJSEngine::setObjectOwnership(mgr, QJSEngine::CppOwnership);
            return mgr;
        });

    QQmlApplicationEngine engine;
    
    // Register cursor image provider
    engine.addImageProvider("cursor", new chromadesk::CursorImageProvider());
    
    // Expose version to QML
    engine.rootContext()->setContextProperty("APP_VERSION", APP_VERSION_STR);
    
    QFontDatabase::addApplicationFont(":/font/SegoeFluentIcons.ttf");

    // Handle QML creation failures
    QObject::connect(&engine, &QQmlApplicationEngine::objectCreationFailed,
        &app, []() {
            LOG_CRITICAL("QML object creation failed!");
            QCoreApplication::exit(-1);
        },
        Qt::QueuedConnection);

    // Load main QML
    engine.loadFromModule("ChromaDesk", "MainWindow");
    LOG_INFO("ChromaDesk started successfully");
    
    int runRet = app.exec();

    chromadesk::OcrEngine::instance().uninitialize();

    LOG_INFO("ChromaDesk exiting with code {}", runRet);

    return runRet;
}
