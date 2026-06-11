// Copyright 2026 ChromaDesk Authors
// Skill bridge handler — routes skillExec / skillListTools from MCP to the
// remote skill host via the qd-skill-bridge WebRTC data channel.

#ifndef CHROMADESK_API_SKILLHANDLER_H
#define CHROMADESK_API_SKILLHANDLER_H

#include <QAtomicInt>
#include <QJsonObject>
#include <QMap>
#include <QMutex>
#include <QObject>
#include <QString>
#include <QWaitCondition>

namespace chromadesk {

class MainController;

/**
 * @brief Handles skillExec and skillListTools MCP tool calls.
 *
 * Thread safety: handleSkillExec / handleSkillListTools are designed to be
 * called from *worker threads* (e.g. QThreadPool).  They block the calling
 * worker thread via QWaitCondition (NOT QEventLoop) until the remote skill host
 * replies or the timeout elapses.
 *
 * onSkillResponse runs on the *main thread* (direct signal connection from
 * ClientManager) and wakes the waiting worker thread.
 */
class SkillHandler : public QObject {
    Q_OBJECT

public:
    explicit SkillHandler(MainController* controller, QObject* parent = nullptr);

    // Thread-safe: can be called from any thread (dispatched via QThreadPool).
    QJsonObject handleSkillExec(const QJsonObject& params);
    QJsonObject handleSkillListTools(const QJsonObject& params);

    // Called on the main thread when a skillBridgeResponse arrives.
    void onSkillResponse(const QString& deviceId,
                         const QJsonObject& response);

private:
    bool isConnectionValid(const QString& deviceId) const;

    QJsonObject sendAndWait(const QString& deviceId,
                            const QJsonObject& payload,
                            int timeoutMs = 10000);

    QString nextRequestId();

    MainController* m_controller;
    QAtomicInt      m_nextId{1};

    struct PendingRequest {
        QMutex*         mutex  = nullptr;
        QWaitCondition* cond   = nullptr;
        QJsonObject*    result = nullptr;
        bool*           done   = nullptr;
    };
    QMutex m_pendingMutex;
    QMap<QString, PendingRequest> m_pending;
};

} // namespace chromadesk

#endif // CHROMADESK_API_SKILLHANDLER_H
