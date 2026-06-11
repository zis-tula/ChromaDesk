#include "SystemTrayManager.h"

#include <QApplication>
#include <QTimer>

#include "core/localconfigcenter.h"
#include "infra/log/log.h"

namespace chromadesk {

SystemTrayManager::SystemTrayManager(QObject* parent)
    : QObject(parent)
{
    initTrayIcon();
}

SystemTrayManager::~SystemTrayManager()
{
    if (m_trayIcon) {
        m_trayIcon->hide();
    }
}

SystemTrayManager& SystemTrayManager::instance()
{
    static SystemTrayManager inst;
    return inst;
}

void SystemTrayManager::initTrayIcon()
{
    m_trayIcon = new QSystemTrayIcon(this);
    m_trayIcon->setIcon(QIcon(":/image/tray/logo.png"));
    m_trayIcon->setToolTip("ChromaDesk");

    m_menu = new QMenu();
    m_showAction = new QAction(tr("Show"), this);
    m_quitAction = new QAction(tr("Quit"), this);
    m_menu->addAction(m_showAction);
    m_menu->addAction(m_quitAction);
    m_trayIcon->setContextMenu(m_menu);
    m_trayIcon->show();

    connect(m_showAction, &QAction::triggered, this, [this]() {
        Q_EMIT showWindowRequested();
    });

    connect(m_quitAction, &QAction::triggered, this, [this]() {
        quit();
    });

    connect(m_trayIcon, &QSystemTrayIcon::activated, this,
        [this](QSystemTrayIcon::ActivationReason reason) {
            if (reason == QSystemTrayIcon::Trigger ||
                reason == QSystemTrayIcon::DoubleClick) {
                Q_EMIT showWindowRequested();
            }
        });
}

void SystemTrayManager::minimizeToTray()
{
    if (!core::LocalConfigCenter::instance().trayMessageShown()) {
        core::LocalConfigCenter::instance().setTrayMessageShown(true);
        m_trayIcon->showMessage(
            tr("ChromaDesk"),
            tr("ChromaDesk is still running in the background."),
            QSystemTrayIcon::Information,
            3000);
    }
}

void SystemTrayManager::quit()
{
    LOG_INFO("SystemTrayManager: quit requested");
    m_trayIcon->hide();
    QTimer::singleShot(200, qApp, &QApplication::quit);
}

}
