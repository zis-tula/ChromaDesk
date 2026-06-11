// Copyright 2026 ChromaDesk Authors
// Skill bridge handler implementation

#include "SkillHandler.h"

#include <QJsonDocument>
#include <QMutexLocker>

#include "../controller/MainController.h"
#include "../manager/ClientManager.h"
#include "infra/log/log.h"

namespace chromadesk {

SkillHandler::SkillHandler(MainController* controller, QObject* parent)
    : QObject(parent)
    , m_controller(controller)
{
    connect(m_controller->clientManager(),
            &ClientManager::skillBridgeResponseReceived,
            this,
            [this](const QString& deviceId, const QJsonObject& response) {
                onSkillResponse(deviceId, response);
            });
}

// ---- Public API (called from worker threads) ----

QJsonObject SkillHandler::handleSkillExec(const QJsonObject& params)
{
    QString deviceId = params["deviceId"].toString();
    QString tool     = params["tool"].toString();
    QJsonValue args  = params.value("args");

    if (deviceId.isEmpty() || tool.isEmpty()) {
        return {{"error", "skillExec: device_id and tool are required"}};
    }

    if (!isConnectionValid(deviceId)) {
        return {{"error", QString("skillExec: device '%1' not found").arg(deviceId)}};
    }

    QJsonObject payload;
    payload["id"]   = nextRequestId();
    payload["type"] = "toolCall";
    payload["tool"] = tool;
    payload["args"] = args.isUndefined() ? QJsonValue(QJsonObject{}) : args;

    return sendAndWait(deviceId, payload);
}

QJsonObject SkillHandler::handleSkillListTools(const QJsonObject& params)
{
    QString deviceId = params["deviceId"].toString();

    if (deviceId.isEmpty()) {
        return {{"error", "skillListTools: device_id is required"}};
    }

    if (!isConnectionValid(deviceId)) {
        return {{"error", QString("skillListTools: device '%1' not found").arg(deviceId)}};
    }

    QJsonObject payload;
    payload["id"]   = nextRequestId();
    payload["type"] = "listTools";

    return sendAndWait(deviceId, payload);
}

// ---- Callback from ClientManager (main thread) ----

void SkillHandler::onSkillResponse(const QString& /*deviceId*/,
                                   const QJsonObject& response)
{
    QString id = response["id"].toString();
    if (id.isEmpty()) return;

    QMutexLocker pendingLock(&m_pendingMutex);
    auto it = m_pending.find(id);
    if (it == m_pending.end()) return;

    PendingRequest& req = it.value();
    QMutexLocker resultLock(req.mutex);
    *req.result = response;
    *req.done   = true;
    req.cond->wakeOne();
}

// ---- Private helpers ----

bool SkillHandler::isConnectionValid(const QString& deviceId) const
{
    auto* cm = m_controller->clientManager();
    return cm && cm->connectedDeviceIds().contains(deviceId);
}

QJsonObject SkillHandler::sendAndWait(const QString& deviceId,
                                      const QJsonObject& payload,
                                      int timeoutMs)
{
    QString id = payload["id"].toString();
    QByteArray bytes = QJsonDocument(payload).toJson(QJsonDocument::Compact);
    QString jsonData = QString::fromUtf8(bytes);

    QMutex mutex;
    QWaitCondition cond;
    QJsonObject result;
    bool done = false;

    {
        QMutexLocker locker(&m_pendingMutex);
        m_pending.insert(id, PendingRequest{&mutex, &cond, &result, &done});
    }

    // NativeMessaging is not thread-safe — send from the main thread.
    auto* cm = m_controller->clientManager();
    QMetaObject::invokeMethod(cm, [cm, deviceId, jsonData]() {
        cm->sendSkillCommand(deviceId, jsonData);
    }, Qt::QueuedConnection);

    // Block the *worker thread* (not the main thread) until response or timeout.
    {
        QMutexLocker locker(&mutex);
        if (!done) {
            cond.wait(&mutex, timeoutMs);
        }
    }

    {
        QMutexLocker locker(&m_pendingMutex);
        m_pending.remove(id);
    }

    if (!done) {
        return {{"error", QString("skillExec timed out after %1 ms").arg(timeoutMs)}};
    }
    return result;
}

QString SkillHandler::nextRequestId()
{
    return QString("skill-%1").arg(m_nextId.fetchAndAddRelaxed(1));
}

} // namespace chromadesk
